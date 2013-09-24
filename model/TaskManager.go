package model

import (
	"fmt"
	"runtime"
	"sync"
)

type Task interface {
	Run()
	SetManager(manager *TaskManager)
}

type TaskManager struct {
	maxTaskSize    int
	channel        chan Task
	aliveTaskCount int
	lock           *sync.Mutex
	w              sync.WaitGroup
}

func NewInstance(size int) *TaskManager {
	manager := &TaskManager{maxTaskSize: size, aliveTaskCount: 0, lock: &sync.Mutex{}}
	manager.channel = make(chan Task, manager.maxTaskSize)
	return manager
}

func (this *TaskManager) AddTask(task Task) {
	defer func() {
		if err := recover(); err != nil {
			this.Stop()
		}
	}()
	if this.channel == nil {
		return
	}
	task.SetManager(this)
	this.channel <- task
}

func (this *TaskManager) Stop() {
	if this.channel != nil {
		close(this.channel)
		this.channel = nil
	}
}

func (this *TaskManager) Run() {
	go func() {
		fmt.Println("TaskManager Run")
		var task Task
		get := true
	Loop:
		for get {
			select {
			case task, get = <-this.channel:
				if get {
					runtime.GC()
					go task.Run()
				} else {
					fmt.Println("channel closed!")
					break Loop
				}

			}
		}
	}()

}
