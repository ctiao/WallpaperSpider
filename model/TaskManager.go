package model

import (
	"fmt"
	// "runtime"
	// "sync"
)

type Task interface {
	Run()
	Cancel()
	SetManager(manager *TaskManager)
}

type TaskManager struct {
	maxTaskSize int
	channel     chan Task
	pool        chan int
	taskIndex   int
	count       int
}

func NewInstance(size int) *TaskManager {
	manager := &TaskManager{maxTaskSize: size}
	manager.channel = make(chan Task, manager.maxTaskSize)
	manager.pool = make(chan int, manager.maxTaskSize)
	return manager
}

func (this *TaskManager) AddTask(task Task) {
	defer func() {
		if err := recover(); err != nil {
			this.Stop()
		}
	}()
	if this.channel == nil || this.pool == nil {
		return
	}
	this.pool <- this.taskIndex
	task.SetManager(this)
	this.channel <- task
	this.taskIndex++
	this.count++
	fmt.Println("task count", this.count)
}

func (this *TaskManager) Stop() {
	if this.pool != nil {
		close(this.pool)
		this.pool = nil
	}
	if this.channel != nil {
		close(this.channel)
		this.channel = nil
	}
}

func (this *TaskManager) Run() {

	fmt.Println("TaskManager Running")
	go func() {
		for task := range this.channel {
			go func() {
				task.Run()
				this.count--
				<-this.pool
				fmt.Println("task count", this.count)
			}()
		}
	}()
}
