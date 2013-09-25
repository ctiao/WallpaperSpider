package model

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	KB      = 1024
	ext     = ".jpg"
	ext_png = ".png"
	ext_tmp = ".tmp"
	UA      = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/29.0.1547.76 Safari/537.36"
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
	task.savePath = fmt.Sprintf("%s%s%s", dir, coder.EncodeToString([]byte(url)), ext)
	return task
}

func (this *DownloadTask) SetManager(manager *TaskManager) {
	this.manager = manager
}

func (this *DownloadTask) Run() {
	if this.url == "" || this.savePath == "" {
		return
	}
	err := this.downloadFile(this.url, this.savePath)
	if err != nil {
		fmt.Printf("%s-下载失败,重试-----\n", this.url)
		this.url = strings.Replace(this.url, ext, ext_png, 1)
		err1 := this.downloadFile(this.url, this.savePath)
		if err1 != nil {
			fmt.Printf("%s-下载失败-------------\n", this.url)
		}
	}
}

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, time.Minute)
}

func (this *DownloadTask) downloadFile(url string, savePath string) error {
	e, err := exists(savePath)
	if e {
		return nil
	}

	tmpName := fmt.Sprintf("%s%s", savePath, ext_tmp)

	if e, _ := exists(tmpName); e {
		os.Remove(tmpName)
	}

	//删除报错文件
	defer func() {
		if err != nil {
			fmt.Println("remove file", savePath, err)
			os.Remove(tmpName)
		} else {
			os.Rename(tmpName, savePath)
		}
	}()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", UA)

	transport := http.Transport{
		Dial: dialTimeout,
	}

	client := http.Client{
		Transport: &transport,
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	out, err := os.Create(tmpName)
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

	if length < KB*3 || length != res.ContentLength {
		err = errors.New(fmt.Sprintf("%s,length:%d/%d", url, length, res.ContentLength))
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
