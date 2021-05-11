package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

//  处理POST请求 -- JSON
func postHandlerJson(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// 请求类型是application/json时从r.Body读取数据
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("read request.Body failed, err:%v\n", err)
		return
	}

	fmt.Println(string(b))

	answer := `{"status": "ok"}`
	w.Write([]byte(answer))
}

// 处理POST请求 -- Form
func postHandlerForm(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	//请求类型是application/x-www-form-urlencoded时解析form数据
	r.ParseForm()
	fmt.Println(r.PostForm) // 打印form数据
	fmt.Println(r.PostForm.Get("name"), r.PostForm.Get("age"))

	answer := `{"status": "ok"}`
	w.Write([]byte(answer))
}

func main() {
	http.HandleFunc("/form", postHandlerForm)
	http.HandleFunc("/json", postHandlerJson)
	http.ListenAndServe(":8080", nil)
}
