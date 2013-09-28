package task

import (
	"fmt"
	"math/rand"
	"time"
)

type PrintTask struct {
	ImplTask
	manager *TaskManager
	Text    string
}

func (this *PrintTask) SetManager(manager *TaskManager) {
	this.manager = manager
}

func (this *PrintTask) Run() {
	s := time.Duration(rand.Intn(120))
	str := fmt.Sprintf("PrintTask.Run() %ds  %v", s, &this)
	fmt.Println(str)
	if this.manager == nil {
		return
	}

	go func() {
		for i := 1; i < 32; i++ {
			fmt.Println(i)
			newTask := new(PrintTask)
			newTask.Text = fmt.Sprintln("task:", i)
			this.manager.AddTask(newTask)
		}
	}()

	// select {
	// case <-time.After(3 * time.Second):
	// 	//fmt.Println("time out")
	// }
}

func (this *PrintTask) Cancel() {

}
