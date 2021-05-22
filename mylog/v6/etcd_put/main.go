package main

import (
	"context"
	"fmt"
	"go_demos/my_demos/mylog/v6/utils"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	defer cli.Close()

	// put
	// value := `[{"path":"/tmp/nginx2.log", "topic":"nginx_log"}]`
	value := `[{"path":"/tmp/nginx.log", "topic":"nginx_log"},{"path":"/tmp/redis.log", "topic":"redis_log"},{"path":"/tmp/mysql.log", "topic":"mysql_log"}]`
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	ip, err := utils.GetOutBoundIp()
	if err != nil {
		return
	}
	key := fmt.Sprintf("/logagent/%s/collect_config", ip)
	_, err = cli.Put(ctx, key, value)
	cancel()
	if err != nil {
		fmt.Printf("put to etcd failed, err:%v\n", err)
		return
	}
}
