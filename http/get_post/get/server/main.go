package main

import (
	"fmt"
	"net/http"
)

// 处理Get请求
func getHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data := r.URL.Query()
	fmt.Println(data.Get("name"))
	fmt.Println(data.Get("age"))
	answer := `{"status": "ok"}`
	w.Write([]byte(answer))
}

func main() {
	http.HandleFunc("/", getHandler)
	http.ListenAndServe(":8080", nil)
}
