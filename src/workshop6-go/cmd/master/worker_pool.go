package main

import (
	"errors"
	"sync"
	"time"

	glog "github.com/golang/glog"
)

// WORKER STATE //
type WorkerInfo struct {
	Address   string
	Heartbeat bool
	Status    int
}

const ( // Worker status
	WorkerStatusIdle = iota
	WorkerStatusBusy
	WorkerStatusDead
)

type WorkerPool struct {
	Queue                []string
	State                map[string]WorkerInfo
	DeadWorkerQueue      []string
	QueueMutex           sync.Mutex
	StateMutex           sync.Mutex
	DeadWorkerQueueMutex sync.Mutex
}

func initWorkerPool() *WorkerPool {
	workerPool := WorkerPool{}
	workerPool.Init()
	return &workerPool
}

func (pool *WorkerPool) Init() {
	pool.Queue = make([]string, 0)
	pool.State = make(map[string]WorkerInfo)
	pool.DeadWorkerQueue = make([]string, 0)
}

func (pool *WorkerPool) RegisterWorker(address string) {
	pool.StateMutex.Lock()
	if _, ok := pool.State[address]; ok {
		defer pool.StateMutex.Unlock()
		glog.Infof("Worker %s already registered", address)
		// Mark as dead if not dead
		workerInfo := pool.State[address]
		if workerInfo.Status != WorkerStatusDead {
			workerInfo.Status = WorkerStatusDead
			pool.State[address] = workerInfo
			pool.DeadWorkerQueueMutex.Lock()
			defer pool.DeadWorkerQueueMutex.Unlock()
			pool.AppendDeadWorker(address)
		}
	} else {
		pool.State[address] = WorkerInfo{
			Address:   address,
			Heartbeat: true,
			Status:    WorkerStatusIdle,
		}
		glog.Infof("Worker %s registered", address)
		pool.StateMutex.Unlock()
		pool.AddIdleWorker(address)
	}
}

func (pool *WorkerPool) AddIdleWorker(address string) {
	pool.QueueMutex.Lock()
	pool.StateMutex.Lock()
	defer pool.StateMutex.Unlock()
	defer pool.QueueMutex.Unlock()
	// glog.Infof("Adding worker %s to idle queue", address)
	pool.Queue = append(pool.Queue, address)
	workerInfo := pool.State[address]
	workerInfo.Status = WorkerStatusIdle
	pool.State[address] = workerInfo
}

func (pool *WorkerPool) GetIdleWorker() (string, error) {
	pool.QueueMutex.Lock()
	pool.StateMutex.Lock()
	defer pool.StateMutex.Unlock()
	defer pool.QueueMutex.Unlock()

	for {
		if len(pool.Queue) == 0 {
			return "", errors.New("no idle workers")
		}
		workerAddress := pool.Queue[0]
		workerInfo := pool.State[workerAddress]
		if workerInfo.Status == WorkerStatusIdle {
			pool.Queue = pool.Queue[1:]
			workerInfo.Status = WorkerStatusBusy
			pool.State[workerAddress] = workerInfo
			// glog.Infof("Remove idle worker from queue %s", workerAddress)
			return workerAddress, nil
		} else {
			// Remove busy/dead worker from queue
			pool.Queue = pool.Queue[1:]
		}
	}
}

func (pool *WorkerPool) PollForIdleWorker() string {
	counter := 0
	for {
		workerAddress, err := pool.GetIdleWorker()
		if err != nil {
			counter++
			if counter%10 == 0 {
				glog.Info("Wait for workers to become alive...")
			}
			time.Sleep(1 * time.Second)
			continue
		}
		return workerAddress
	}
}

func (pool *WorkerPool) MarkWorkerHeartbeat(address string) {
	pool.StateMutex.Lock()
	defer pool.StateMutex.Unlock()
	workerInfo := pool.State[address]
	workerInfo.Heartbeat = true
	pool.State[address] = workerInfo
}

func (pool *WorkerPool) MarkWorkersDeadIfNoHeartbeat() {
	pool.StateMutex.Lock()
	defer pool.StateMutex.Unlock()
	for address, workerInfo := range pool.State {
		if !workerInfo.Heartbeat {
			glog.Infof("Worker %s is dead", address)
			workerInfo.Status = WorkerStatusDead
			pool.DeadWorkerQueueMutex.Lock()
			defer pool.DeadWorkerQueueMutex.Unlock()
			pool.AppendDeadWorker(address)
		}
		workerInfo.Heartbeat = false
		pool.State[address] = workerInfo
	}
}

func (pool *WorkerPool) AppendDeadWorker(deadWorkerAddress string) {
	pool.DeadWorkerQueue = append(pool.DeadWorkerQueue, deadWorkerAddress)
}

func (pool *WorkerPool) GetDeadWorkerAddresses() []string {
	pool.DeadWorkerQueueMutex.Lock()
	defer pool.DeadWorkerQueueMutex.Unlock()
	deadWorkerAddresses := make([]string, len(pool.DeadWorkerQueue))
	copy(deadWorkerAddresses, pool.DeadWorkerQueue)
	pool.DeadWorkerQueue = pool.DeadWorkerQueue[:0]
	return deadWorkerAddresses
}

func (pool *WorkerPool) GetWorkerAddresses() []string {
	pool.StateMutex.Lock()
	defer pool.StateMutex.Unlock()
	addresses := make([]string, 0)
	for address := range pool.State {
		addresses = append(addresses, address)
	}
	return addresses
}

func (pool *WorkerPool) PrettyPrintWorkerStatuses() string {
	// pool.StateMutex.Lock()
	// defer pool.StateMutex.Unlock()
	statuses := ""
	for address, workerInfo := range pool.State {
		statuses += address + ": "
		switch workerInfo.Status {
		case WorkerStatusIdle:
			statuses += "idle"
		case WorkerStatusBusy:
			statuses += "busy"
		case WorkerStatusDead:
			statuses += "dead"
		}
		statuses += "\n"
	}
	return statuses
}
