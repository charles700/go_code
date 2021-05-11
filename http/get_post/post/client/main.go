package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// simpleGet 简单的Get请求示例
func simpleGet() {
	resp, err := http.Get("http://www.baidu.com/")
	if err != nil {
		fmt.Printf("get failed, err:%v\n", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read from resp.Body failed, err:%v\n", err)
		return
	}
	fmt.Println(string(body))
}

// withParmasGet 带有参数的 get请求
func withParmasGet() {
	apiUrl := "http://www.baidu.com/s"
	// URL param
	data := url.Values{}
	data.Set("wd", "小王子")
	data.Set("age", "18")
	u, err := url.ParseRequestURI(apiUrl)
	if err != nil {
		fmt.Printf("parse url requestUrl failed, err:%v\n", err)
	}
	u.RawQuery = data.Encode() // URL encode
	fmt.Println(u.String())

	// 发送请求
	resp, err := http.Get(u.String())

	if err != nil {
		fmt.Printf("post failed, err:%v\n", err)
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("get resp failed, err:%v\n", err)
		return
	}
	fmt.Println(string(b))
}

func main() {
	simpleGet()

	withParmasGet()
}
