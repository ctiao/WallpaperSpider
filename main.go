package main

import (
	"./model"
	"fmt"
	"os"
	"runtime"
)

func main() {
	saveDir := "./pics/"
	os.Mkdir(saveDir, 0777)
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Print("start=====\n")
	taskManager := model.NewInstance(10)
	task := model.NewFetchTaskInstance(1, 10, saveDir)
	taskManager.AddTask(task)
	taskManager.Run()
	var str, str1 string
	fmt.Scanf("%s%s", &str, &str1)
	if str == "q" || str == "Q" {
		taskManager.Stop()
		fmt.Println("taskManager stop")
	}
}
