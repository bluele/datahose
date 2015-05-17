package tasks

import (
	"github.com/bluele/datahose/lib/connections"
	"github.com/bluele/datahose/lib/plugins"
	"sync"
	"time"
)

type TaskManager struct {
	tasks    map[string]*Task
	channels map[string]chan [][]byte
	mutex    sync.RWMutex
}

var Manager *TaskManager

func NewTaskManager() *TaskManager {
	return &TaskManager{tasks: make(map[string]*Task)}
}

func (tm *TaskManager) CreateTask(pl plugins.Plugin) {
	tk := NewTask(pl)
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	cmap := make(map[string]chan []byte)
	connections.Pool.SetConnectionMap(pl.Name(), cmap)
	pl.Init(cmap)
	tm.tasks[pl.Name()] = tk
}

func (tm *TaskManager) Execute() {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()
	for _, task := range tm.tasks {
		go task.Execute()
		go func(task *Task) {
			ch := task.DataChannel()
			for {
				if task.IsRunning() {
					select {
					case outs := <-ch:
						for _, out := range outs {
							connections.Pool.PutMessage(task.Name(), out)
						}
					default:
					}
				} else {
					time.Sleep(task.Interval())
				}
			}
		}(task)
	}
}

type Task struct {
	plugins.Plugin
}

func NewTask(pl plugins.Plugin) *Task {
	return &Task{pl}
}

func init() {
	Manager = NewTaskManager()
}
