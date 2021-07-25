package main

import (
	"fmt"
	"go_demos/my_demos/rpc/cmd"
	"net/rpc"
)

func main() {
	client, err := rpc.DialHTTP("tcp", ":1234")
	if err != nil {
		fmt.Println("链接rpc服务器失败:", err)
	}
	var reply float64
	var args = cmd.DivArgs{
		A: 1,
		B: 2,
	}
	err = client.Call("DemoService.Div", args, &reply)
	if err != nil {
		fmt.Println("调用远程服务失败", err)
	}
	fmt.Println("远程服务返回结果：", reply)
}
