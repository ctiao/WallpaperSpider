package main

import (
	"./model"
	"fmt"
	"os"
	"runtime"
	"strconv"
)

func convertToInteger(str string, defVal int) int {
	if str == "" {
		return defVal
	}
	result, err := strconv.Atoi(str)
	if err != nil {
		result = defVal
	}
	return result
}

func main() {

	args := os.Args[1:]
	var startPage = 1
	var endPage = startPage
	if len(args) > 0 {
		startPage = convertToInteger(args[0], 1)
		endPage = startPage
	}
	if len(args) > 1 {
		endPage = convertToInteger(args[1], startPage)
	}

	saveDir := "./pics/"
	os.Mkdir(saveDir, 0777)
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Print("start=====\n")
	taskManager := model.NewInstance(10)
	task := model.NewFetchTaskInstance(startPage, endPage, saveDir)
	taskManager.AddTask(task)
	taskManager.Run()
	var str, str1 string
	fmt.Scanf("%s%s", &str, &str1)
	if str == "q" || str == "Q" {
		taskManager.Stop()
		fmt.Println("taskManager stop")
	}
}