package model

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

type FetchTask struct {
	startPage int
	endPage   int
	manager   *TaskManager
}

func (this *FetchTask) SetManager(manager *TaskManager) {
	this.manager = manager
}

func (this *FetchTask) Run() {
	for p := this.startPage; p <= this.endPage; p++ {
		fmt.Println("======================fetch=====================")
		url := this.nextUrl(p)
		fmt.Println(url)
		content, err := this.getHTML(url)
		if err != nil {
			continue
		}
		//保存处理内容
		//fmt.Println(content)
		this.getThumbImgUrls(content)

		//开启图片下载任务
	}
	fmt.Println("Fetch end")
}

func (this *FetchTask) getThumbImgUrls(content string) (urls []string, err error) {
	re, err := regexp.Compile("<img[^>]+?data-original=\"([^\"]+)\"")
	if err != nil {
		return
	}
	urlArr := re.FindAllSubmatch([]byte(content), -1)
	//fmt.Println(urlArr)
	if urlArr == nil || len(urlArr) == 0 {
		return
	}
	length := len(urlArr)
	urls = make([]string, length)
	for k, v := range urlArr {
		//fmt.Printf("%d,%v\n\n", k, string(v[1]))
		urls[k] = string(v[1])
	}
	fmt.Println(urls)
	return
}

// 得到内容
func (this *FetchTask) getHTML(url string) (content string, err error) {
	var resp *http.Response
	resp, err = http.Get(url)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	} else {
		fmt.Println("ERROR " + url + " 返回为空 ")
	}
	if resp == nil || resp.Body == nil || err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("ERROR " + url)
		fmt.Println(err)
		return
	}

	var buf []byte
	buf, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	content = string(buf)
	return
}

var URL_FORMAT = "http://wallbase.cc/toplist/index/%d?section=wallpapers&q=&res_opt=eqeq&res=0x0&thpp=32&purity=100&board=1&aspect=0.00&ts=3d"
var PAGE_SIZE = 32

func (this *FetchTask) nextUrl(page int) string {
	url := fmt.Sprintf(URL_FORMAT, PAGE_SIZE*(page-1))
	return url
}

func NewFetchTaskInstance(sp int, ep int) *FetchTask {
	task := &FetchTask{startPage: sp, endPage: ep}
	return task
}
