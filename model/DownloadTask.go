package model

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	ext = "jpg"
	UA  = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/29.0.1547.76 Safari/537.36"
)

var coder = base64.StdEncoding

type DownloadTask struct {
	url      string
	savePath string
	manager  *TaskManager
}

func NewDownloadTaskInstance(url string, dir string) *DownloadTask {
	task := &DownloadTask{}
	task.url = url
	task.savePath = fmt.Sprintf("%s%s.%s", dir, coder.EncodeToString([]byte(url)), ext)
	return task
}

func (this *DownloadTask) SetManager(manager *TaskManager) {
	this.manager = manager
}

func (this *DownloadTask) Run() {
	if this.url == "" || this.savePath == "" {
		return
	}
	fmt.Printf("========%s-下载开始\n", this.url)
	err := this.downloadFile(this.url, this.savePath)
	if err == nil {
		fmt.Printf("%s-下载成功\n", this.url)
	} else {
		fmt.Printf("%s-下载失败-----\n", this.url)
	}
}

func (this *DownloadTask) downloadFile(url string, savePath string) error {
	exists, err := exists(savePath)
	if exists {
		return err
	}

	//删除报错文件
	defer func() {
		if err != nil {
			fmt.Println("remove file")
			err1 := os.Remove(savePath)
			if err1 != nil {
				fmt.Println(err1)
			}
		}
	}()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", UA)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	out, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer func() {
		out.Close()
		out = nil
	}()

	var length int64 = 0
	length, err = io.Copy(out, res.Body)

	if err != nil {
		fmt.Println("io Copy error", length, res.ContentLength)
		return err
	}

	return err
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
