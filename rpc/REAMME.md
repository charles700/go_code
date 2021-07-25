golang 中的 rpc
------

## 三种rpc方式

#### net/rpc
golang官方的net/rpc库使用encoding/gob进行编解码，支持tcp或http数据传输方式.

##### 缺点：
由于其他语言不支持gob编解码方式，所以使用net/rpc库实现的RPC方法没办法进行跨语言调用。


#### net/rpc/jsonrpc
JSON RPC采用JSON进行数据编解码，因而支持跨语言调用。但目前的jsonrpc库是基于tcp协议实现的，暂时不支持使用http进行数据传输。

#### 第三方 protorpc
protobuf进行数据编解码，根据protobuf声明文件自动生成rpc方法定义与服务注册代码，在golang中可以很方便的进行rpc服务调用。
1. 编写 message.proto 文件
2. 在 pb 目录下 执行`protoc ./message.proto --go_out=./` 生成 pb.go 文件
3. 编写 server/main.go， client/main.go
4. go mod tidy
5. 运行 server 和 client


#### 第三方 grpc
