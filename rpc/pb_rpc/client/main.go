package main

import (
	"fmt"
	message "go_demos/my_demos/rpc/pb_rpc/pb"
	"net/rpc"
	"time"
)

func main() {
	client, err := rpc.DialHTTP("tcp", "localhost:8081")

	if err != nil {
		panic(err.Error())
	}

	timeStamp := time.Now().Unix()

	request := message.OrderRequest{OrderId: "201907310001", TimeStamp: timeStamp}

	var response *message.OrderInfo

	// 调用 rpc方法
	err = client.Call("OrderService.GetOrderInfo", request, &response)

	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("%+v\n", response)
}
