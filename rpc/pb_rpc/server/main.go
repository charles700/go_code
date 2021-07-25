package main

import (
	"errors"
	message "go_demos/my_demos/rpc/pb_rpc/pb"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

//订单服务
type OrderService struct {
}

func (os *OrderService) GetOrderInfo(request message.OrderRequest, response *message.OrderInfo) error {

	//201907310003

	orderMap := map[string]message.OrderInfo{

		"201907300001": message.OrderInfo{OrderId: "201907300001", OrderName: "衣服", OrderStatus: "已付款"},

		"201907310001": message.OrderInfo{OrderId: "201907310001", OrderName: "零食", OrderStatus: "已付款"},

		"201907310002": message.OrderInfo{OrderId: "201907310002", OrderName: "食品", OrderStatus: "未付款"},
	}

	current := time.Now().Unix()

	if request.TimeStamp > current {

		*response = message.OrderInfo{OrderId: "0", OrderName: "", OrderStatus: "订单信息异常"}

	} else {

		result := orderMap[request.OrderId] //201907310003

		if result.OrderId != "" {

			*response = orderMap[request.OrderId]

		} else {

			return errors.New("server error")

		}

	}

	return nil

}

func main() {
	orderService := new(OrderService)
	rpc.Register(orderService)
	rpc.HandleHTTP()
	listen, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err.Error())
	}
	http.Serve(listen, nil)
}
