package main

import (
	"go_demos/my_demos/rpc/cmd"
	"net"
	"net/http"
	"net/rpc"
)

func main() {
	/*将服务对象进行注册*/
	err := rpc.Register(&cmd.DemoService{})
	if err != nil {
		panic(err)
	}
	/*将函数的服务注册到HTTP上*/
	rpc.HandleHTTP()

	/*在特定的端口进行监听*/
	listen, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}

	http.Serve(listen, nil)
}
