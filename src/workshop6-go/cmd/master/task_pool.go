package main

import (
	"errors"
	"fmt"
	"sync"

	pb "workshop/protos"

	glog "github.com/golang/glog"
)

const ( // Stage
	StageMap = iota
	StageReduce
)

const ( // Task status
	TaskStatusIdle = iota
	TaskStatusInProgress
	TaskStatusDone
)

type Task struct {
	Id     int32
	Input  []*pb.BlobInfo
	Output []*pb.BlobInfo
	Status int
}

type TaskPool struct {
	MapTasks         map[int32]*Task
	ReduceTasks      map[int32]*Task
	ActiveTasks      map[string]*Task
	MapTasksMutex    sync.Mutex
	ReduceTasksMutex sync.Mutex
	ActiveTasksMutex sync.Mutex
}

func initTaskPool(m int, r int) *TaskPool {
	taskPool := TaskPool{}
	taskPool.Init(m, r)
	return &taskPool
}

func (taskPool *TaskPool) Init(m int, r int) {
	taskPool.MapTasks = make(map[int32]*Task)
	taskPool.ReduceTasks = make(map[int32]*Task)
	taskPool.ActiveTasks = make(map[string]*Task)
	for i := 0; i < m; i++ {
		taskPool.MapTasks[int32(i)] = &Task{
			Id:     int32(i),
			Input:  []*pb.BlobInfo{},
			Output: []*pb.BlobInfo{},
			Status: TaskStatusIdle,
		}
	}
	for i := 0; i < r; i++ {
		taskPool.ReduceTasks[int32(i)] = &Task{
			Id:     int32(i),
			Input:  []*pb.BlobInfo{},
			Output: []*pb.BlobInfo{},
			Status: TaskStatusIdle,
		}
	}
}

func (taskPool *TaskPool) InitMapTask(taskNum int32, input []*pb.BlobInfo) {
	taskPool.MapTasksMutex.Lock()
	defer taskPool.MapTasksMutex.Unlock()
	mapTask := taskPool.MapTasks[taskNum]
	mapTask.Input = input
	taskPool.MapTasks[taskNum] = mapTask
}

func (taskPool *TaskPool) HandleMapSuccess(workerAddress string, taskNum int32, output []*pb.BlobInfo) {
	// Record output
	taskPool.MapTasksMutex.Lock()
	defer taskPool.MapTasksMutex.Unlock()
	taskPool.ActiveTasksMutex.Lock()
	defer taskPool.ActiveTasksMutex.Unlock()
	// Confirm it is the correct task
	if mapTask, ok := taskPool.MapTasks[taskNum]; !ok || mapTask.Status != TaskStatusInProgress {
		return
	}
	if task, ok := taskPool.ActiveTasks[workerAddress]; !ok || task.Id != taskNum {
		return
	}
	glog.Infof("Handling map success for %s, %d", workerAddress, taskNum)
	mapTask := taskPool.MapTasks[taskNum]
	mapTask.Output = output
	mapTask.Status = TaskStatusDone
	taskPool.MapTasks[taskNum] = mapTask
	// Remove from active tasks
	if task, ok := taskPool.ActiveTasks[workerAddress]; ok {
		if task.Id == taskNum {
			delete(taskPool.ActiveTasks, workerAddress)
		}
	}
}

func (taskPool *TaskPool) HandleReduceSuccess(workerAddress string, taskNum int32, output []*pb.BlobInfo) {
	// Record output
	taskPool.ReduceTasksMutex.Lock()
	defer taskPool.ReduceTasksMutex.Unlock()
	taskPool.ActiveTasksMutex.Lock()
	defer taskPool.ActiveTasksMutex.Unlock()
	// Confirm it is the correct task
	if reduceTask, ok := taskPool.ReduceTasks[taskNum]; !ok || reduceTask.Status != TaskStatusInProgress {
		return
	}
	if task, ok := taskPool.ActiveTasks[workerAddress]; !ok || task.Id != taskNum {
		return
	}
	glog.Infof("Handling reduce success for %s, %d", workerAddress, taskNum)
	reduceTask := taskPool.ReduceTasks[taskNum]
	reduceTask.Output = output
	reduceTask.Status = TaskStatusDone
	taskPool.ReduceTasks[taskNum] = reduceTask
	// Remove from active tasks
	if task, ok := taskPool.ActiveTasks[workerAddress]; ok {
		if task.Id == taskNum {
			delete(taskPool.ActiveTasks, workerAddress)
		}
	}
}

func (taskPool *TaskPool) HandleMapFailure(workerAddress string, taskNum int32) {
	taskPool.MapTasksMutex.Lock()
	defer taskPool.MapTasksMutex.Unlock()
	taskPool.ActiveTasksMutex.Lock()
	defer taskPool.ActiveTasksMutex.Unlock()
	// Confirm it is the correct task
	if mapTask, ok := taskPool.MapTasks[taskNum]; !ok || mapTask.Status != TaskStatusInProgress {
		return
	}
	if task, ok := taskPool.ActiveTasks[workerAddress]; !ok || task.Id != taskNum {
		return
	}
	glog.Infof("Handling map failure for %s, %d", workerAddress, taskNum)
	mapTask := taskPool.MapTasks[taskNum]
	mapTask.Status = TaskStatusIdle
	taskPool.MapTasks[taskNum] = mapTask
	// Remove from active tasks
	if task, ok := taskPool.ActiveTasks[workerAddress]; ok {
		if task.Id == taskNum {
			delete(taskPool.ActiveTasks, workerAddress)
		}
	}
}

func (taskPool *TaskPool) HandleReduceFailure(workerAddress string, taskNum int32) {
	taskPool.ReduceTasksMutex.Lock()
	defer taskPool.ReduceTasksMutex.Unlock()
	taskPool.ActiveTasksMutex.Lock()
	defer taskPool.ActiveTasksMutex.Unlock()
	// Confirm it is the correct task
	if reduceTask, ok := taskPool.ReduceTasks[taskNum]; !ok || reduceTask.Status != TaskStatusInProgress {
		return
	}
	if task, ok := taskPool.ActiveTasks[workerAddress]; !ok || task.Id != taskNum {
		return
	}
	glog.Infof("Handling reduce failure for %s, %d", workerAddress, taskNum)
	reduceTask := taskPool.ReduceTasks[taskNum]
	reduceTask.Status = TaskStatusIdle
	taskPool.ReduceTasks[taskNum] = reduceTask
	// Remove from active tasks
	if task, ok := taskPool.ActiveTasks[workerAddress]; ok {
		if task.Id == taskNum {
			delete(taskPool.ActiveTasks, workerAddress)
		}
	}
}

func (task *TaskPool) AssignIdleTask(workerAddress string) (*Task, error) {
	stage := task.GetStage()
	idleTask := int32(-1)
	if stage == StageMap {
		task.MapTasksMutex.Lock()
		defer task.MapTasksMutex.Unlock()
		task.ActiveTasksMutex.Lock()
		defer task.ActiveTasksMutex.Unlock()
		for _, mapTask := range task.MapTasks {
			// Clean up done tasks
			if mapTask.Status == TaskStatusDone {
				glog.Infof("Cleaning up map task %d", mapTask.Id)
				task.updateReduceTaskInputs(mapTask.Id)
				delete(task.MapTasks, mapTask.Id)
				continue
			}
			// Assign idle task
			if mapTask.Status == TaskStatusIdle {
				mapTask.Status = TaskStatusInProgress
				task.ActiveTasks[workerAddress] = mapTask
				idleTask = mapTask.Id
				// Assign worker to task
				task.ActiveTasks[workerAddress] = task.MapTasks[idleTask]
				break
			}
		}
		if idleTask == -1 {
			return nil, errors.New("no idle map task")
		}
	} else if stage == StageReduce {
		task.ReduceTasksMutex.Lock()
		defer task.ReduceTasksMutex.Unlock()
		task.ActiveTasksMutex.Lock()
		defer task.ActiveTasksMutex.Unlock()
		for _, reduceTask := range task.ReduceTasks {
			// Clean up done tasks
			if reduceTask.Status == TaskStatusDone {
				glog.Infof("Cleaning up reduce task %d", reduceTask.Id)
				delete(task.ReduceTasks, reduceTask.Id)
				continue
			}
			// Assign idle task
			if reduceTask.Status == TaskStatusIdle {
				reduceTask.Status = TaskStatusInProgress
				task.ActiveTasks[workerAddress] = reduceTask
				idleTask = reduceTask.Id
				// Assign worker to task
				task.ActiveTasks[workerAddress] = task.ReduceTasks[idleTask]
				break
			}
		}
		if idleTask == -1 {
			return nil, errors.New("no idle reduce task")
		}
	} else {
		return nil, errors.New("invalid stage")
	}

	glog.Infof("Assigned task %d to worker %s", idleTask, workerAddress)
	return task.ActiveTasks[workerAddress], nil
}

func (task *TaskPool) GetActiveTaskForWorker(workerAddress string) (*Task, error) {
	task.ActiveTasksMutex.Lock()
	defer task.ActiveTasksMutex.Unlock()
	if task, ok := task.ActiveTasks[workerAddress]; ok {
		return task, nil
	}
	return nil, errors.New("no active task for worker")
}

func (task *TaskPool) updateReduceTaskInputs(mapTaskId int32) {
	// Used when Map mutex already locked
	task.ReduceTasksMutex.Lock()
	defer task.ReduceTasksMutex.Unlock()
	for i, blobInfo := range task.MapTasks[mapTaskId].Output {
		reduceTask := task.ReduceTasks[int32(i)]
		reduceTask.Input = append(reduceTask.Input, blobInfo)
		task.ReduceTasks[int32(i)] = reduceTask
		// glog.Infof("Updated reduce task %d with input %s/%s", reduceTask.Id, blobInfo.ContainerName, blobInfo.BlobName)
	}
}

func (task *TaskPool) MarkWorkerFailure(workerAddress string) {
	if task.GetStage() == StageMap {
		// Lock and delete from map tasks
		task.MapTasksMutex.Lock()
		defer task.MapTasksMutex.Unlock()
		task.ActiveTasksMutex.Lock()
		defer task.ActiveTasksMutex.Unlock()
		if activeTask, ok := task.ActiveTasks[workerAddress]; ok {
			taskId := activeTask.Id
			delete(task.ActiveTasks, workerAddress)
			// Mark map task as idle
			if mapTask, ok := task.MapTasks[taskId]; ok {
				glog.Infof("Marking map task %d as idle due to worker %s failure", taskId, workerAddress)
				mapTask.Status = TaskStatusIdle
				task.MapTasks[taskId] = mapTask
			}
		}
	} else if task.GetStage() == StageReduce {
		// Lock and delete from reduce tasks
		task.ReduceTasksMutex.Lock()
		defer task.ReduceTasksMutex.Unlock()
		task.ActiveTasksMutex.Lock()
		defer task.ActiveTasksMutex.Unlock()
		if activeTask, ok := task.ActiveTasks[workerAddress]; ok {
			taskId := activeTask.Id
			delete(task.ActiveTasks, workerAddress)
			// Mark reduce task as idle
			if reduceTask, ok := task.ReduceTasks[taskId]; ok {
				glog.Infof("Marking reduce task %d as idle due to worker %s failure", taskId, workerAddress)
				reduceTask.Status = TaskStatusIdle
				task.ReduceTasks[taskId] = reduceTask
			}
		}
	}
}

func (task *TaskPool) GetStage() int {
	// Check size of map tasks and reduce tasks
	task.MapTasksMutex.Lock()
	defer task.MapTasksMutex.Unlock()
	task.ReduceTasksMutex.Lock()
	defer task.ReduceTasksMutex.Unlock()
	if len(task.MapTasks) > 0 {
		return StageMap
	}
	if len(task.ReduceTasks) > 0 {
		return StageReduce
	}
	return -1
}

func (task *TaskPool) GetMapTasksRemaining() int {
	// task.MapTasksMutex.Lock()
	// defer task.MapTasksMutex.Unlock()
	return len(task.MapTasks)
}

func (task *TaskPool) GetReduceTasksRemaining() int {
	// task.ReduceTasksMutex.Lock()
	// defer task.ReduceTasksMutex.Unlock()
	return len(task.ReduceTasks)
}

func (task *TaskPool) PrettyPrintMapTasks() (string, error) {
	// task.MapTasksMutex.Lock()
	// defer task.MapTasksMutex.Unlock()
	if len(task.MapTasks) == 0 {
		return "", errors.New("no map tasks")
	}
	result := ""
	for _, mapTask := range task.MapTasks {
		mapTaskString := ""
		mapTaskString += "Map Task " + fmt.Sprintf("%d", mapTask.Id) + ": "
		mapTaskString += "Input: "
		for _, blobInfo := range mapTask.Input {
			mapTaskString += blobInfo.ContainerName + "/" + blobInfo.BlobName + ", "
		}
		mapTaskString += "Output: "
		for _, blobInfo := range mapTask.Output {
			mapTaskString += blobInfo.ContainerName + "/" + blobInfo.BlobName + ", "
		}
		mapTaskString += "Status: " + fmt.Sprintf("%d", mapTask.Status)
		// Query worker it is assigned to
		for worker, task := range task.ActiveTasks {
			if task.Id == mapTask.Id {
				result += "Worker: " + worker + ", "
				break
			}
		}
		result += mapTaskString + "\n"
	}
	return result, nil
}

func (task *TaskPool) PrettyPrintReduceTasks() (string, error) {
	// task.ReduceTasksMutex.Lock()
	// defer task.ReduceTasksMutex.Unlock()
	if len(task.ReduceTasks) == 0 {
		return "", errors.New("no reduce tasks")
	}
	result := ""
	for _, reduceTask := range task.ReduceTasks {
		reduceTaskString := ""
		reduceTaskString += "Reduce Task " + fmt.Sprintf("%d", reduceTask.Id) + ": "
		reduceTaskString += "Input: "
		for _, blobInfo := range reduceTask.Input {
			reduceTaskString += blobInfo.ContainerName + "/" + blobInfo.BlobName + ", "
		}
		reduceTaskString += "Output: "
		for _, blobInfo := range reduceTask.Output {
			reduceTaskString += blobInfo.ContainerName + "/" + blobInfo.BlobName + ", "
		}
		reduceTaskString += "Status: " + fmt.Sprintf("%d", reduceTask.Status)
		// Query worker it is assigned to
		for worker, task := range task.ActiveTasks {
			if task.Id == reduceTask.Id {
				result += "Worker: " + worker + ", "
				break
			}
		}
		result += reduceTaskString + "\n"
	}
	return result, nil
}

func (task *TaskPool) PrettyPrintActiveTasks() (string, error) {
	// task.ActiveTasksMutex.Lock()
	// defer task.ActiveTasksMutex.Unlock()
	if len(task.ActiveTasks) == 0 {
		return "", errors.New("no active tasks")
	}
	result := ""
	for worker, task := range task.ActiveTasks {
		result += "Worker: " + worker + ", "
		result += "Task: " + fmt.Sprintf("%d", task.Id) + "\n"
	}
	return result, nil
}
