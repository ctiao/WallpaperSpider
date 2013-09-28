package taskpool

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
	g_dialTimeout = time.Duration(time.Second * 120)
	KB            = 1024
	BUFFER_SIZE   = KB * 20
	ext           = ".jpg"
	ext_png       = ".png"
	ext_tmp       = ".tmp"
	UA            = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/29.0.1547.76 Safari/537.36"
)

var coder = base64.StdEncoding

type DownloadTask struct {
	url        string
	savePath   string
	manager    *TaskManager
	cancelFlag bool
}

func NewDownloadTaskInstance(u string, dir string) *DownloadTask {
	task := new(DownloadTask)
	task.cancelFlag = false
	task.url = u
	task.savePath = fmt.Sprintf("%s%s%s", dir, coder.EncodeToString([]byte(task.url)), ext)
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

	conn, _ := net.DialTimeout(network, addr, g_dialTimeout)
	//conn.SetDeadline(time.Now().Add(time.Minute * 2))
	return conn, nil
}

func (this *DownloadTask) downloadFile(url string, savePath string) error {

	defer func() {
		if ret := recover(); ret != nil {
			fmt.Println("somethins wrong,recovered from downloading file")
		}
	}()

	e := Exists(savePath)
	if e {
		fmt.Println("文件已存在,跳过")
		return nil
	}

	fmt.Println(url)

	tmpName := fmt.Sprintf("%s%s", savePath, ext_tmp)

	if Exists(tmpName) {
		os.Remove(tmpName)
	}

	var err error
	var out *os.File
	var in io.ReadCloser

	//删除报错文件
	defer func() {
		if in != nil {
			in.Close()
		}
		if out != nil {
			out.Close()
		}
		if err != nil {
			if Exists(tmpName) {
				err1 := os.Remove(tmpName)
				fmt.Println("下载失败,remove file", tmpName, err1)
			}
		} else {
			err1 := os.Rename(tmpName, savePath)
			fmt.Println("下载成功", savePath, err1)
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
	in = res.Body
	defer func() {
		if in != nil {
			in.Close()
			in = nil
		}
	}()

	if res.StatusCode != 200 || res.ContentLength == -1 {
		err = errors.New("网络文件不存在")
		return err
	}

	//创建临时文件
	out, err = os.OpenFile(tmpName, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	defer func() {
		if out != nil {
			out.Close()
			out = nil
		}
	}()

	// length, err := io.Copy(out, res.Body)
	pbytes := make([]byte, BUFFER_SIZE)

	var length int64 = 0
	var readed int
	for {
		readed, err = in.Read(pbytes)
		if err != nil && err != io.EOF {
			break
		}
		if readed > 0 {
			written, err1 := out.Write(pbytes[:readed])
			length += int64(written)
			if err1 != nil {
				fmt.Println(err1)
				break
			}
		}
		if readed == 0 || err == io.EOF {
			err = nil
			break
		}
		if this.cancelFlag {
			break
		}
	}

	if err != nil || length < KB*3 || length != res.ContentLength {
		err = errors.New(fmt.Sprintf("%s,length:%d/%d", url, length, res.ContentLength))
	}

	if in != nil {
		in.Close()
		in = nil
	}
	if out != nil {
		out.Close()
		out = nil
	}

	return err
}

// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func (this *DownloadTask) Cancel() {
	this.cancelFlag = true
}
