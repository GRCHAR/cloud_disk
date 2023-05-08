package server_tasks

import (
	"crypto/md5"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"io"
	"sync"
	"time"
)

type serverTaskManager struct {
	tasks map[string]*serverTask
	lock  sync.RWMutex
}

type serverTask struct {
	taskId string
	task   func()
	stop   chan bool
	loop   bool
	//second
	afterTime time.Duration
	logger    *logrus.Logger
}

func newServerTaskManager() *serverTaskManager {
	return &serverTaskManager{tasks: make(map[string]*serverTask)}
}

var serverTaskMgr = newServerTaskManager()

func GetServerTaskManager() *serverTaskManager {
	return serverTaskMgr
}

func init() {

}

// AddTask Add a task to the task list
func (stm *serverTaskManager) AddTask(task func(), stop chan bool, loop bool, afterTime int) string {
	uuid := getUUID()
	newServerTask := &serverTask{
		taskId:    uuid,
		task:      task,
		stop:      stop,
		loop:      loop,
		afterTime: time.Duration(afterTime),
	}
	stm.runTask(newServerTask)
	return uuid
}

// RunTask Run a task
func (*serverTaskManager) runTask(task *serverTask) {
	serverTaskMgr.lock.Lock()
	defer serverTaskMgr.lock.Unlock()
	serverTaskMgr.tasks[task.taskId] = task
	go func() {
		for {
			select {
			case <-task.stop:
				delete(serverTaskMgr.tasks, task.taskId)
				return
			case <-time.After(task.afterTime * time.Second):
				task.task()
				if !task.loop {
					delete(serverTaskMgr.tasks, task.taskId)
					return
				}
			}
		}
	}()

}

// StopTask Stop a task
func (*serverTaskManager) stopTask(uuid string) {
	if serverTaskMgr.tasks[uuid] != nil {
		serverTaskMgr.tasks[uuid].stop <- true
	}

}

func (*serverTaskManager) deleteTask() {

}

func strToMd5(data string) string {
	t := md5.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

func getUUID() string {
	return uuid.New().String()
}
