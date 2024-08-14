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

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"
	gateway "workshop/common/gateways"
	pb "workshop/protos"

	glog "github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"go.etcd.io/etcd/client/v3/concurrency"
)

const (
	defaultMessage = "Hello"
)

var (
	name              = flag.String("name", "master", "the name of the master")
	port              = flag.Int("port", 8080, "The server port")
	httpAddress       = flag.String("httpAddress", ":8090", "Address listening for http requests")
	failureNode       = flag.Int("failureNode", -1, "Node to fail")
	failurePoint      = flag.Int("failurePoint", -1, "Set failure point for master")
	launchTime        = flag.Int("launchTime", 0, "Time deployment was launched")
	taskCount         = atomic.Int32{}
	isLeader     bool = false

	userRequestQueue      = make([]ContainerRequest, 0)
	userRequestQueueMutex sync.Mutex

	azureGateway *gateway.AzureGateway = nil
	workerPool   *WorkerPool           = nil
	taskPool     *TaskPool             = nil
)

type ContainerRequest struct {
	ContainerName string `json:"containerName"`
	M             int    `json:"m"`
	R             int    `json:"r"`
}

type BlobData struct {
	Name string
	Size int
}

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedMasterServer
}

// MapperRegister implements masterServer/MapperRegister
func (server) MapperRegister(ctx context.Context, in *pb.MapperInfo) (*pb.MapperRegisterReply, error) {
	if !isLeader {
		glog.Infof("Reject registration of %v", in.GetName())
		return &pb.MapperRegisterReply{
			Success: false,
		}, nil
	}
	glog.Infof("Registration of %v", in.GetName())
	workerPool.RegisterWorker(in.GetName())
	err := Replicate()

	if err != nil {
		return nil, err
	}
	return &pb.MapperRegisterReply{
		Success: true,
	}, nil
}

func (server) Replicate(ctx context.Context, in *pb.ReplicationRequest) (*pb.ReplicationReply, error) {
	glog.Infof("Received replication request %d", in.Seq)
	reply, err := ReceiveReplication(in)
	if err != nil {
		glog.Error(err.Error())
		return nil, err
	}
	return reply, nil
}

func startGRPCServer(portNum int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", portNum))
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	pb.RegisterMasterServer(s, &server{})
	glog.Infof("Master server listening at %v", lis.Addr().String())
	err = s.Serve(lis)
	if err != nil {
		return err
	}
	return nil
}

func sendReduceRequest(workerName string, currentBlobInfos []*pb.BlobInfo, currentReduceTaskNum int32) error {
	// connect to registered worker
	conn, err := grpc.Dial(workerName, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := pb.NewWorkerClient(conn)

	client_ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	request := &pb.ReduceRequest{
		ReduceTaskNum: currentReduceTaskNum,
		KeyValueInfos: currentBlobInfos,
	}

	// worker shold get back as registration message again after perform
	glog.Infof("Sending reduce request %d to %s", request.ReduceTaskNum, workerName)
	reply, err := client.PerformReduce(client_ctx, request)
	if err != nil {
		return err
	}
	if reply.GetSuccess() {
		workerPool.AddIdleWorker(workerName)
		if taskPool != nil {
			taskPool.HandleReduceSuccess(workerName, currentReduceTaskNum, []*pb.BlobInfo{reply.GetResultInfo()})
		}
	} else {
		glog.Infof("Reduce request %d to %s failed", currentReduceTaskNum, workerName)
		taskPool.HandleReduceFailure(workerName, currentReduceTaskNum)
	}

	if client_ctx.Err() == context.Canceled {
		// Worker heartbeat success but failed to perform map within the deadline
		glog.Infof("Worker %s failed to perform reduce %d within the deadline", workerName, currentReduceTaskNum)
		if taskPool != nil {
			taskPool.MarkWorkerFailure(workerName)
		}
	}
	return nil
}

func estimateWorkerTaskTime(blobInfos []*pb.BlobInfo) int {
	totalSize := 0
	for _, blobInfo := range blobInfos {
		totalSize += int((*blobInfo).RangeEnd) - int((*blobInfo).RangeStart)
	}
	// glog.Infof("Got size of %d for map work", totalSize)
	// add 10s per 1MB
	defaultTime := 5
	additionalTimePerMB := 10
	return defaultTime + additionalTimePerMB*(totalSize/1000000)
}

func sendMapRequest(mapperName string, currentBlobInfos []*pb.BlobInfo, currentMapTaskNum int32, reduceSize int) error {
	// connect to registered worker
	conn, err := grpc.Dial(mapperName, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := pb.NewWorkerClient(conn)

	estimatedDeadline := estimateWorkerTaskTime(currentBlobInfos)

	start := time.Now()
	client_ctx, cancel := context.WithTimeout(context.Background(), time.Duration(estimatedDeadline)*time.Second)
	defer cancel()

	// Send shard info to mapper
	request := &pb.MapRequest{
		MapTaskNum:  currentMapTaskNum,
		ShardInfos:  currentBlobInfos,
		ReduceTasks: int32(reduceSize),
	}
	// worker shold get back as registration message again after perform
	glog.Infof("Sending map request %d to %s", currentMapTaskNum, mapperName)
	reply, err := client.PerformMap(client_ctx, request)
	if err != nil {
		return err
	}

	if reply.GetSuccess() {
		elapsed := time.Since(start)
		glog.Infof("Time taken for map work %d: %s\n", currentMapTaskNum, elapsed)

		workerPool.AddIdleWorker(mapperName)
		if taskPool != nil {
			taskPool.HandleMapSuccess(mapperName, currentMapTaskNum, reply.GetKeyValueInfos())
		}
	} else {
		glog.Infof("Map request %d to %s failed", currentMapTaskNum, mapperName)
		taskPool.HandleMapFailure(mapperName, currentMapTaskNum)
	}

	if client_ctx.Err() == context.Canceled {
		// Worker heartbeat success but failed to perform map within the deadline
		glog.Infof("Worker %s failed to perform map %d within the deadline", mapperName, currentMapTaskNum)
		if taskPool != nil {
			taskPool.MarkWorkerFailure(mapperName)
		}
	}
	return nil
}

func shardBlobs(blobMetadata []BlobData, containerName string, mapSize int) {
	currentMapTaskNum := 0

	for currentMapTaskNum = 0; currentMapTaskNum < mapSize; currentMapTaskNum += 1 {
		currentBlobInfos := make([]*pb.BlobInfo, 0)
		for _, blobData := range blobMetadata {
			rangeStart := currentMapTaskNum * blobData.Size / mapSize
			rangeEnd := (currentMapTaskNum+1)*blobData.Size/mapSize - 1
			if currentMapTaskNum == mapSize {
				rangeEnd = blobData.Size
			}
			currentBlobInfo := pb.BlobInfo{
				RangeStart:    int64(rangeStart),
				RangeEnd:      int64(rangeEnd),
				ContainerName: containerName,
				BlobName:      blobData.Name,
			}
			currentBlobInfos = append(currentBlobInfos, &currentBlobInfo)
		}

		taskPool.InitMapTask(int32(currentMapTaskNum), currentBlobInfos)
	}
}

func handleError(err error) {
	if err != nil {
		glog.Error(err.Error())
	}
}

func getBlobMetadata(containerName string) []BlobData {
	client, err := azureGateway.GetBlobClient(nil)
	handleError(err)
	pager := client.NewListBlobsFlatPager(containerName, &azblob.ListBlobsFlatOptions{
		Include: azblob.ListBlobsInclude{Metadata: true, Snapshots: true, Versions: true},
	})
	blobMetadata := []BlobData{}
	glog.Infof("Collecting blob data from %s", containerName)
	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		handleError(err)
		for _, blob := range resp.Segment.BlobItems {
			newData := BlobData{*blob.Name, int(*blob.Properties.ContentLength)}
			blobMetadata = append(blobMetadata, newData)
			glog.Infof("Read blob name: %s, size: %v", newData.Name, newData.Size)
		}
	}
	return blobMetadata
}

func initAzureGateway() gateway.AzureGateway {
	if azureGateway != nil {
		return *azureGateway
	}
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	azureGateway = &gateway.AzureGateway{AccountName: accountName, AccountKey: accountKey}
	return *azureGateway
}

func handleDeadWorkers(deadWorkerAddresses []string) {
	for _, deadWorkerAddress := range deadWorkerAddresses {
		if taskPool != nil {
			taskPool.MarkWorkerFailure(deadWorkerAddress)
		}
		glog.Infof("Worker %s dead, re-assigning task", deadWorkerAddress)
	}
}

// HEARTBEAT //
func sendHeartbeatToWorkers() {
	heartbeatCounter := 0
	heartbeatTimeout := 10 * time.Second
	for {
		// Send heartbeat to all workers
		workerAddresses := workerPool.GetWorkerAddresses()
		for _, workerAddress := range workerAddresses {
			go func(workerAddress string) {
				err := sendHeartbeat(workerAddress, &pb.HeartbeatRequest{}, heartbeatTimeout)
				if err != nil {
					handleError(err)
				}
			}(workerAddress)
		}
		heartbeatCounter += 1
		// Wait for heartbeat timeout
		time.Sleep(heartbeatTimeout)
		// Mark workers dead if no heartbeat
		workerPool.MarkWorkersDeadIfNoHeartbeat()
		deadWorkerAddresses := workerPool.GetDeadWorkerAddresses()
		handleDeadWorkers(deadWorkerAddresses)
		if heartbeatCounter%6 == 0 {
			prettyPrint := workerPool.PrettyPrintWorkerStatuses()
			glog.Infof("Worker pool: \n%s", prettyPrint)
		}
	}
}

func sendHeartbeat(workerAddress string, heartbeatRequest *pb.HeartbeatRequest, heartbeatTimeout time.Duration) error {
	conn, err := grpc.Dial(workerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := pb.NewWorkerClient(conn)

	client_ctx, cancel := context.WithTimeout(context.Background(), heartbeatTimeout)
	defer cancel()

	heartbeatReply, err := client.Heartbeat(client_ctx, heartbeatRequest)
	if err != nil {
		return err
	}
	workerPool.MarkWorkerHeartbeat(workerAddress)
	if taskPool != nil {
		// Check if task is correct for worker; else, mark current active task as failed
		activeTask, err := taskPool.GetActiveTaskForWorker(workerAddress)
		if err == nil { // Active task
			// Mark worker as failed if task is incorrect
			if taskPool.GetStage() == StageMap && heartbeatReply.GetMapTaskNum() != activeTask.Id {
				glog.Infof("Worker %s has incorrect map task %d, marking as failed", workerAddress, activeTask.Id)
				if heartbeatReply.GetStatus() == pb.TaskStatus_COMPLETED {
					workerPool.AddIdleWorker(workerAddress)
				}
				taskPool.HandleMapFailure(workerAddress, activeTask.Id)
			} else if taskPool.GetStage() == StageReduce && heartbeatReply.GetReduceTaskNum() != activeTask.Id {
				glog.Infof("Worker %s has incorrect reduce task %d, marking as failed", workerAddress, activeTask.Id)
				if heartbeatReply.GetStatus() == pb.TaskStatus_COMPLETED {
					workerPool.AddIdleWorker(workerAddress)
				}
				taskPool.HandleReduceFailure(workerAddress, activeTask.Id)
			}
			// Mark task as success if worker has completed task
			if taskPool.GetStage() == StageMap && heartbeatReply.GetStatus() == pb.TaskStatus_COMPLETED && heartbeatReply.GetMapTaskNum() == activeTask.Id {
				workerPool.AddIdleWorker(workerAddress)
				taskPool.HandleMapSuccess(workerAddress, heartbeatReply.GetMapTaskNum(), heartbeatReply.GetBlobInfos())
			} else if taskPool.GetStage() == StageReduce && heartbeatReply.GetStatus() == pb.TaskStatus_COMPLETED && heartbeatReply.GetReduceTaskNum() == activeTask.Id {
				workerPool.AddIdleWorker(workerAddress)
				taskPool.HandleReduceSuccess(workerAddress, heartbeatReply.GetReduceTaskNum(), heartbeatReply.GetBlobInfos())
			}
		}
	}

	return nil
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		glog.Infof("Received request with method %s, but only POST is allowed", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !isLeader {
		err := Proxy(w, r)
		handleError(err)
		return
	}

	// Read the request body
	var containerRequest ContainerRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&containerRequest)
	if err != nil {
		http.Error(w, "Error decoding JSON body", http.StatusBadRequest)
		return
	}

	// Process the request body as needed
	glog.Infof("Received POST request with container name: %v\n", containerRequest)

	userRequestQueueMutex.Lock()
	userRequestQueue = append(userRequestQueue, containerRequest)
	userRequestQueueMutex.Unlock()

	// Send a response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("POST request received successfully"))
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Health check OK"))
}

func incrementTaskCount() {
	taskCount.Add(1)
	nodeNumber, _ = getNodeNumber()
	if nodeNumber == *failureNode {
		if taskCount.Load() == int32(*failurePoint) {
			// Resign
			cli, err := getEtcdClient()
			if err != nil {
				glog.Fatal(err)
			}
			// create a sessions to elect a Leader
			c, err := concurrency.NewSession(cli)
			if err != nil {
				glog.Fatal(err)
			}
			defer c.Close()
			election := concurrency.NewElection(c, "/mapreduce")
			ctx := context.Background()
			if err := election.Resign(ctx); err != nil {
				glog.Fatal(err)
			}
			glog.Infof("I am no longer the leader")
			// Do something as a leader
			isLeader = false
			glog.Fatalf("Failure on %s as scheduled after %v task assignment", *name, taskCount.Load())
		} else {
			glog.Infof("[MASTER NODE FAILURE] %d/%d tasks assigned", taskCount.Load(), *failurePoint)
		}
	}
}

func main() {
	flag.Parse()
	flag.Lookup("alsologtostderr").Value.Set("true")

	defer glog.Flush()
	var message = fmt.Sprintf("%s from %s!", defaultMessage, *name)
	glog.Info(message)

	initAzureGateway()
	workerPool = initWorkerPool()

	go func() {
		http.HandleFunc("/word-count", handlePost)
		http.HandleFunc("/healthz", handleHealthz)
		glog.Infof("%s listening on %v to serve user POST request", *name, *httpAddress)
		http.ListenAndServe(*httpAddress, nil)
	}()
	// Set up a gRPC server to listen to mappers with goroutine
	go func() {
		err := startGRPCServer(*port)
		handleError(err)
	}()

	time.Sleep(5 * time.Second)
	podIP := os.Getenv("POD_IP")
	// Force lower IP to be leader in order for failure to work
	firstLeader, err := getLeader()
	if err != nil {
		glog.Fatal(err)
	}
	if firstLeader != podIP {
		time.Sleep(5 * time.Second)
	}

	// Poll for leader
	glog.Infof("Polling for leader")
	cli, err := getEtcdClient()
	if err != nil {
		glog.Fatal(err)
	}
	// create a sessions to elect a Leader
	c, err := concurrency.NewSession(cli)
	if err != nil {
		glog.Fatal(err)
	}
	defer c.Close()
	election := concurrency.NewElection(c, "/mapreduce")
	ctx := context.Background()

	// Elect a leader (or wait that the leader resign)
	if err := election.Campaign(ctx, podIP); err != nil {
		glog.Fatal(err)
	}

	glog.Infof("I am the leader")
	// Do something as a leader
	isLeader = true
	// log all of the data structures
	glog.Infof("workerPool: %v", workerPool)
	glog.Infof("taskPool: %v", taskPool)
	glog.Infof("userRequestQueue: %v", userRequestQueue)

	if err != nil {
		glog.Fatal(err)
	}

	// Send heartbeat to workers
	go func() {
		sendHeartbeatToWorkers()
	}()

	userRequestLogged := false
	for {
		// If backup, resume current task
		if taskPool == nil {
			if len(userRequestQueue) == 0 {
				if !userRequestLogged {
					glog.Info("Waiting for user requests ...")
					userRequestLogged = true
				}
				time.Sleep(1 * time.Second)
				continue
			}
			userRequestLogged = false
			userRequestQueueMutex.Lock()
			containerRequest := userRequestQueue[0]
			userRequestQueue = userRequestQueue[1:]
			userRequestQueueMutex.Unlock()
			containerName := containerRequest.ContainerName
			mapSize := containerRequest.M
			reduceSize := containerRequest.R

			client, err := azureGateway.GetBlobClient(nil)
			handleError(err)
			ctx := context.Background()
			// Create the mapper container
			mapperContainerName := containerName + "-mapped"
			_, err = client.CreateContainer(ctx, mapperContainerName, nil)
			// handleError(err)  // Ignore error if container already exists

			// Create the reducer container
			reducerContainerName := containerName + "-mapped-reduced"
			_, err = client.CreateContainer(ctx, reducerContainerName, nil)
			// handleError(err)  // Ignore error if container already exists

			taskPool = initTaskPool(mapSize, reduceSize)
			blobMetadata := getBlobMetadata(containerName)
			glog.Infof("Got blob metadata: %v", blobMetadata)
			// Shard blobs
			shardBlobs(blobMetadata, containerName, mapSize)
			glog.Infof("Sharded blobs into %d map tasks", mapSize)
		}
		glog.Infof("[MILESTONE] Finished initialization")

		// Block for map tasks
		mapSize := len(taskPool.MapTasks)
		reduceSize := len(taskPool.ReduceTasks)
		mapPrettyPrint, _ := taskPool.PrettyPrintMapTasks()
		reducePrettyPrint, _ := taskPool.PrettyPrintReduceTasks()
		activePrettyPrint, _ := taskPool.PrettyPrintActiveTasks()
		glog.Infof("Map size: %d, reduce size: %d", mapSize, reduceSize)
		glog.Infof("Map task map: \n%s", mapPrettyPrint)
		glog.Infof("Reduce task map: \n%s", reducePrettyPrint)
		glog.Infof("Active task map: \n%s", activePrettyPrint)
		glog.Infof("Worker pool: %v", workerPool)
		mapLogCounter := 0
		for taskPool.GetStage() == StageMap {
			time.Sleep(1 * time.Second)
			mapLogCounter += 1
			if mapLogCounter%10 == 0 {
				mapTasksRemaining := taskPool.GetMapTasksRemaining()
				glog.Infof("Waiting for %d/%d map tasks to finish", mapTasksRemaining, mapSize)
				// mapPrettyPrint, _ := taskPool.PrettyPrintMapTasks()
				// activePrettyPrint, _ := taskPool.PrettyPrintActiveTasks()
				// glog.Infof("Map task map: \n%s", mapPrettyPrint)
				// glog.Infof("Active task map: \n%s", activePrettyPrint)
			}
			worker := workerPool.PollForIdleWorker()
			task, err := taskPool.AssignIdleTask(worker)
			if err != nil {
				workerPool.AddIdleWorker(worker)
				continue
			}
			// glog.Infof("Assigned task %d to worker %s", task.Id, worker)
			err = Replicate()
			handleError(err)
			go func(worker string, task *Task) {
				sendMapRequest(worker, task.Input, task.Id, reduceSize)
			}(worker, task)
			incrementTaskCount()
		}
		glog.Infof("[MILESTONE] Finished map tasks")
		// Block for reduce tasks
		reduceLogCounter := 0
		for taskPool.GetStage() == StageReduce {
			time.Sleep(1 * time.Second)
			reduceLogCounter += 1
			if reduceLogCounter%10 == 0 {
				reduceTasksRemaining := taskPool.GetReduceTasksRemaining()
				glog.Infof("Waiting for %d/%d reduce tasks to finish", reduceTasksRemaining, reduceSize)
			}
			worker := workerPool.PollForIdleWorker()
			task, err := taskPool.AssignIdleTask(worker)
			if err != nil {
				workerPool.AddIdleWorker(worker)
				continue
			}
			err = Replicate()
			handleError(err)
			go func(worker string, task *Task) {
				sendReduceRequest(worker, task.Input, task.Id)
			}(worker, task)
			incrementTaskCount()
		}
		glog.Infof("[MILESTONE] Finished reduce tasks")
		taskPool = nil
		err = Replicate()
		handleError(err)
	}

}
