package main

import (
	"./model"
	"fmt"
	// "runtime"
)

func main() {
	//runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Print("start=====\n")
	taskManager := model.NewInstance(5)
	//task := &model.PrintTask{Text: "hello"}
	task := model.NewFetchTaskInstance(1, 10, "c:\\temp\\")
	taskManager.AddTask(task)
	fmt.Printf("end=====%v\n", taskManager)
	taskManager.Run()
	var str, str1 string
	fmt.Scanf("%s%s", &str, &str1)
	if str == "q" || str == "Q" {
		taskManager.Stop()
		fmt.Println("taskManager stop")
	}
	//fmt.Scanln(&str1)
}
