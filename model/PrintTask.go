package model

import (
	"fmt"
)

type PrintTask struct {
	manager *TaskManager
	Text    string
}

func (this *PrintTask) SetManager(manager *TaskManager) {
	this.manager = manager
}

func (this *PrintTask) Run() {
	str := fmt.Sprintf("PrintTask.Run()  %v", &this)
	fmt.Println(str)
	if this.manager == nil {
		return
	}
	newTask := &PrintTask{Text: str}
	for i := 1; i < 5; i++ {
		this.manager.AddTask(newTask)
	}
}

func (this *PrintTask) Cancel() {

}
