package main

import (
	"context"
	"flag"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sort"
	"sync/atomic"
	"time"
	pb "workshop/protos"

	glog "github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

var (
	seq                         = atomic.Int64{}
	nodeNumber                  = -1
	etcdClient *clientv3.Client = nil
	namespace                   = flag.String("namespace", "mapreduce", "the namespace of the cluster")
	masterAddr                  = flag.String("masterAddr", "master-service."+*namespace+".svc.cluster.local", "the master service's DNS")
)

func getEtcdClient() (*clientv3.Client, error) {
	if etcdClient != nil {
		return etcdClient, nil
	}
	cli, err := clientv3.New(clientv3.Config{Endpoints: []string{"etcd." + *namespace + ".svc.cluster.local:2379"}})
	if err != nil {
		return nil, err
	}
	etcdClient = cli
	return cli, nil
}

func getNodeNumber() (int, error) {
	// get node number ascending
	// if nodeNumber != -1 {
	// 	return nodeNumber, nil
	// }
	addrs, err := net.LookupHost(*masterAddr)
	if err != nil {
		return -1, err
	}
	// sort addrs
	sort.Strings(addrs)
	for i, addr := range addrs {
		if addr == os.Getenv("POD_IP") {
			nodeNumber = i
			break
		}
	}
	return nodeNumber, nil
}

func getLeader() (string, error) {
	addrs, err := net.LookupHost(*masterAddr)
	// find smallest address alphabetically
	sort.Strings(addrs)
	smallest := addrs[0]

	return smallest, err
}

func PollForLeader() (*concurrency.Election, error) {
	glog.Infof("Polling for leader")
	cli, err := getEtcdClient()
	if err != nil {
		return nil, err
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
	time.Sleep(10 * time.Second)
	podIP := os.Getenv("POD_IP")
	if err := election.Campaign(ctx, podIP); err != nil {
		return nil, err
	}

	glog.Infof("I am the leader")
	// Do something as a leader
	isLeader = true
	return election, nil
}

func Resign() error {
	glog.Infof("Resigning")
	cli, err := getEtcdClient()
	if err != nil {
		return err
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
		return err
	}
	glog.Infof("I am no longer the leader")
	// Do something as a leader
	isLeader = false
	return nil
}

func Proxy(w http.ResponseWriter, r *http.Request) error {
	cli, err := clientv3.New(clientv3.Config{Endpoints: []string{"etcd:2379"}})
	if err != nil {
		return err
	}
	c, err := concurrency.NewSession(cli)
	if err != nil {
		glog.Fatal(err)
	}
	defer c.Close()
	election := concurrency.NewElection(c, "/mapreduce")
	ctx := context.Background()
	resp, err := election.Leader(ctx)
	if err != nil {
		return err
	}

	addr := resp.Kvs[0].Value
	glog.Infof("Proxying to %s", string(addr))
	target, err := url.Parse("http://" + string(addr) + ":8090/word-count")
	if err != nil {
		return err
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Host = target.Host
	}
	glog.Infof("Proxying to %s", target)

	proxy.ServeHTTP(w, r)

	return nil
}

func Replicate() error {
	/*
		Workflow:
		1. Get other master services from kubernetes
		2. Transform userRequestQueue, workerPool, and taskPool into protobuf versions
		3. Send (blocking) to other master services
	*/
	// 1. Get other master services from kubernetes
	addrs, err := net.LookupHost(*masterAddr)
	if err != nil {
		return err
	}
	// glog.Infof("Found master at %s", addrs)

	// 2. Transform userRequestQueue, workerPool, and taskPool into protobuf versions
	seq.Add(1)
	urq := []*pb.ContainerRequest{}
	for _, req := range userRequestQueue {
		urq = append(urq, &pb.ContainerRequest{
			ContainerName: req.ContainerName,
			M:             int32(req.M),
			R:             int32(req.R),
		})
	}
	wp := &pb.WorkerPool{}
	wp.Queue = []string{}
	wp.State = make(map[string]*pb.WorkerInfo)
	wp.DeadWorkerQueue = []string{}
	wp.Queue = append(wp.Queue, workerPool.Queue...)
	for worker, info := range workerPool.State {
		wp.State[worker] = &pb.WorkerInfo{
			Address:   info.Address,
			Heartbeat: info.Heartbeat,
			Status:    int32(info.Status),
		}
	}
	wp.DeadWorkerQueue = append(wp.DeadWorkerQueue, workerPool.DeadWorkerQueue...)
	tp := &pb.TaskPool{}
	if taskPool != nil {
		tp.MapTasks = make(map[int32]*pb.Task)
		tp.ReduceTasks = make(map[int32]*pb.Task)
		tp.ActiveTasks = make(map[string]*pb.Task)
		for taskNum, mapTask := range taskPool.MapTasks {
			tp.MapTasks[taskNum] = &pb.Task{
				Id:     mapTask.Id,
				Input:  mapTask.Input,
				Output: mapTask.Output,
				Status: int32(mapTask.Status),
			}
		}
		for taskNum, reduceTask := range taskPool.ReduceTasks {
			tp.ReduceTasks[taskNum] = &pb.Task{
				Id:     reduceTask.Id,
				Input:  reduceTask.Input,
				Output: reduceTask.Output,
				Status: int32(reduceTask.Status),
			}
		}
		for worker, task := range taskPool.ActiveTasks {
			tp.ActiveTasks[worker] = &pb.Task{
				Id:     task.Id,
				Input:  task.Input,
				Output: task.Output,
				Status: int32(task.Status),
			}
		}
	}

	// 3. Send (blocking) to other master services
	// Remove own IP
	podIP := os.Getenv("POD_IP")
	for i, addr := range addrs {
		if addr == podIP {
			addrs = append(addrs[:i], addrs[i+1:]...)
			break
		}
	}
	for i := 0; i < len(addrs); i++ {
		conn, err := grpc.Dial(addrs[i]+":8080", grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return err
		}
		defer conn.Close()
		client := pb.NewMasterClient(conn)

		client_ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		replicationRequest := &pb.ReplicationRequest{
			Seq:              seq.Load(),
			UserRequestQueue: urq,
			WorkerPool:       wp,
			TaskPool:         tp,
		}
		// glog.Infof("Sending replication request to %s", addrs[i])
		_, err = client.Replicate(client_ctx, replicationRequest)
		if err != nil {
			glog.Infof("Error replicating seq %d on %s", seq.Load(), addrs[i])
			return err
		}
		glog.Infof("Successfully replicated seq %d on %s", seq.Load(), addrs[i])
	}
	return nil
}

func ReceiveReplication(in *pb.ReplicationRequest) (*pb.ReplicationReply, error) {
	if isLeader {
		return &pb.ReplicationReply{
			Success: true,
		}, nil
	}
	if in.Seq <= seq.Load() {
		return &pb.ReplicationReply{
			Success: true,
		}, nil
	}

	// Transform protobuf versions into userRequestQueue, workerPool, and taskPool
	userRequestQueue = []ContainerRequest{}
	for _, req := range in.UserRequestQueue {
		userRequestQueue = append(userRequestQueue, ContainerRequest{
			ContainerName: req.ContainerName,
			M:             int(req.M),
			R:             int(req.R),
		})
	}
	workerPool = &WorkerPool{}
	workerPool.Queue = []string{}
	workerPool.State = make(map[string]WorkerInfo)
	workerPool.DeadWorkerQueue = []string{}
	workerPool.Queue = append(workerPool.Queue, in.WorkerPool.Queue...)
	for worker, info := range in.WorkerPool.State {
		workerPool.State[worker] = WorkerInfo{
			Address:   info.Address,
			Heartbeat: info.Heartbeat,
			Status:    int(info.Status),
		}
	}
	workerPool.DeadWorkerQueue = append(workerPool.DeadWorkerQueue, in.WorkerPool.DeadWorkerQueue...)
	taskPool = &TaskPool{}
	if in.TaskPool != nil {
		taskPool.MapTasks = make(map[int32]*Task)
		taskPool.ReduceTasks = make(map[int32]*Task)
		taskPool.ActiveTasks = make(map[string]*Task)
		for taskNum, mapTask := range in.TaskPool.MapTasks {
			taskPool.MapTasks[taskNum] = &Task{
				Id:     mapTask.Id,
				Input:  mapTask.Input,
				Output: mapTask.Output,
				Status: int(mapTask.Status),
			}
		}
		for taskNum, reduceTask := range in.TaskPool.ReduceTasks {
			taskPool.ReduceTasks[taskNum] = &Task{
				Id:     reduceTask.Id,
				Input:  reduceTask.Input,
				Output: reduceTask.Output,
				Status: int(reduceTask.Status),
			}
		}
		// Check for stage in order to point ActiveTasks to the correct map/reduce tasks
		stage := taskPool.GetStage()
		if stage == StageMap {
			for worker, task := range in.TaskPool.ActiveTasks {
				taskPool.ActiveTasks[worker] = taskPool.MapTasks[task.Id]
			}
		} else {
			for worker, task := range in.TaskPool.ActiveTasks {
				taskPool.ActiveTasks[worker] = taskPool.ReduceTasks[task.Id]
			}
		}
	}

	// glog.Infof("Successfully replicated seq %d", in.Seq)
	return &pb.ReplicationReply{
		Success: true,
	}, nil
}
