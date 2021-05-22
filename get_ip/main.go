package main

import (
	"fmt"
	"net"
	"strings"
)

func GetOutBoundIp() (ip string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return
	}

	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(localAddr.String())

	ip = strings.Split(localAddr.IP.String(), ":")[0]
	return
}

func main() {
	ip, err := GetOutBoundIp()
	if err != nil {
		fmt.Println("GetOutBoundIp faild, err: ", err)
		return
	}
	fmt.Println("本地IP", ip)
}