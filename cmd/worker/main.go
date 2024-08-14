/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
	pb "workshop/protos"

	mapper "workshop/cmd/worker/mapper"
	reducer "workshop/cmd/worker/reducer"

	gateway "workshop/common/gateways"

	glog "github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	// clientv3 "go.etcd.io/etcd/client/v3"
	// "go.etcd.io/etcd/client/v3/concurrency"
)

const (
	defaultMessage = "Hello"
)

var (
	name        = flag.String("name", "worker", "the name of the worekr")
	serviceName = flag.String("serviceName", "worker-service", "the name of the worker service")
	port        = flag.Int("port", 8080, "The server port")
	// port                               = flag.Int("port", 8081, "The server port")
	failurePoint                       = flag.Int("failurePoint", 0, "Set failure point for this worker")
	azureGateway *gateway.AzureGateway = nil
	masterAddr                         = flag.String("masterAddr", "master-service.mapreduce.svc.cluster.local", "the address to connect to master")
	masterPort                         = flag.Int("masterPort", 8080, "The master port")
	// masterAddr    = flag.String("masterAddr", "localhost:8080", "the address to connect to master")
	taskMutex     sync.Mutex
	mapTaskNum    int32 = -1
	reduceTaskNum int32 = -1

	taskCount                 = atomic.Int32{}
	taskStatus pb.TaskStatus  = pb.TaskStatus_COMPLETED
	taskResult []*pb.BlobInfo = nil
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedWorkerServer
}

// Mapper
func (s *server) PerformMap(ctx context.Context, in *pb.MapRequest) (*pb.MapReply, error) {
	err := setTask(in.GetMapTaskNum(), -1)
	if err != nil {
		glog.Error(err.Error())
		reply := &pb.MapReply{
			Success: false,
		}
		return reply, err
	}
	glog.Infof("Received map: %v", in.GetShardInfos())
	reply, err := mapper.Map(azureGateway, ctx, in)
	taskResult = reply.GetKeyValueInfos()
	taskStatus = pb.TaskStatus_COMPLETED
	glog.Infof("Reply: %v", reply)
	return reply, err
}

// Reducer
func (s *server) PerformReduce(ctx context.Context, in *pb.ReduceRequest) (*pb.ReduceReply, error) {
	err := setTask(-1, in.GetReduceTaskNum())
	if err != nil {
		glog.Error(err.Error())
		reply := &pb.ReduceReply{
			Success: false,
		}
		return reply, err
	}
	glog.Infof("Received reduce: %v", in.GetKeyValueInfos())
	reply, err := reducer.Reduce(azureGateway, ctx, in)
	taskResult = []*pb.BlobInfo{reply.GetResultInfo()}
	taskStatus = pb.TaskStatus_COMPLETED
	glog.Infof("Reply: %v", reply)
	return reply, err
}

func (s *server) Heartbeat(ctx context.Context, in *pb.HeartbeatRequest) (*pb.HeartbeatReply, error) {
	glog.Infof("Received heartbeat, sending reply: status: %v, mapTaskNum: %v, reduceTaskNum: %v, taskResult: %v", taskStatus, mapTaskNum, reduceTaskNum, taskResult)
	taskMutex.Lock()
	defer taskMutex.Unlock()
	return &pb.HeartbeatReply{
		Status:        taskStatus,
		MapTaskNum:    mapTaskNum,
		ReduceTaskNum: reduceTaskNum,
		BlobInfos:     taskResult,
	}, nil
}

func setTask(mtNum int32, rtNum int32) error {
	taskMutex.Lock()
	defer taskMutex.Unlock()

	if taskStatus == pb.TaskStatus_IN_PROGRESS {
		return errors.New("task already in progress")
	}
	// Reset
	mapTaskNum = -1
	reduceTaskNum = -1
	taskStatus = pb.TaskStatus_IN_PROGRESS
	taskResult = nil

	// Check for failure
	taskCount.Add(1)
	if taskCount.Load() == int32(*failurePoint) {
		glog.Fatalf("Failure on %s as scheduled after %v map request", *name, taskCount.Load())
		os.Exit(0)
	}

	if mtNum >= 0 {
		mapTaskNum = mtNum
	}
	if rtNum >= 0 {
		reduceTaskNum = rtNum
	}

	return nil
}

func registerToMaster() error {
	// connect to registered worker
	glog.Infof("Connecting to master at %s", *masterAddr)
	addrs, err := net.LookupHost(*masterAddr)
	glog.Infof("Found master at %s", addrs)
	i := 0
	if err != nil {
		return err
	}
	for {
		glog.Infof("Dialing %s", addrs[i]+":"+fmt.Sprint(*masterPort))
		conn, err := grpc.Dial(addrs[i]+":"+fmt.Sprint(*masterPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return err
		}
		defer conn.Close()
		client := pb.NewMasterClient(conn)

		client_ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Send shard info to mapper
		// TODO: get current dns server name
		// currentServerName := "localhost:8081"
		currentServerName := *serviceName + ".mapreduce.svc.cluster.local:8080"
		request := &pb.MapperInfo{Name: currentServerName}

		// worker shold get back as registration message again after perform
		glog.Info("Registering to master...")
		reply, err := client.MapperRegister(client_ctx, request)
		if err != nil {
			return err
		}
		if reply.GetSuccess() {
			glog.Info("Registered to master!")
			return nil
		}
		time.Sleep(1 * time.Second)
		// Loop again for next address
		i = (i + 1) % len(addrs)
	}
}

func startGRPCServer(portNum int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", portNum))
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	pb.RegisterWorkerServer(s, &server{})
	glog.Infof("Master server listening at %v", lis.Addr())
	err = s.Serve(lis)
	if err != nil {
		return err
	}
	return nil
}

func handleError(err error) {
	if err != nil {
		glog.Error(err.Error())
	}
}

func initAzureGateway() gateway.AzureGateway {
	if azureGateway != nil {
		return *azureGateway
	}
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	azureGateway = &gateway.AzureGateway{AccountName: accountName, AccountKey: accountKey}
	return *azureGateway
}

func main() {
	flag.Parse()
	flag.Lookup("alsologtostderr").Value.Set("true")

	defer glog.Flush()
	var message = fmt.Sprintf("%s from %s!", defaultMessage, *name)
	glog.Info(message)

	initAzureGateway()

	// Wait 5s for master initialization
	time.Sleep(5 * time.Second)

	// register mapper to master
	err := registerToMaster()
	handleError(err)

	// Set up a gRPC server to listen to mappers
	err = startGRPCServer(*port)
	handleError(err)
}
