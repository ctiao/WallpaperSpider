package main

import (
	"./model"
	"fmt"
	// "runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Print("start=====\n")
	taskManager := model.NewInstance(10)
	//task := &model.PrintTask{Text: "hello"}
	task := model.NewFetchTaskInstance(1, 10, "c:\\temp\\")
	taskManager.AddTask(task)
	taskManager.Run()
	var str, str1 string
	fmt.Scanf("%s%s", &str, &str1)
	if str == "q" || str == "Q" {
		taskManager.Stop()
		fmt.Println("taskManager stop")
	}
}
