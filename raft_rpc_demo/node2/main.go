package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"sync"
	"time"

	"log"
)

//对每个节点id和端口的封装类型
type nodeInfo struct {
	id   string
	port string //代表ip+port
}

//声明节点对象类型Raft
type Raft struct {
	node nodeInfo
	mu   sync.Mutex
	//当前节点编号
	me            int
	currentTerm   int
	votedFor      int
	state         int
	timeout       int
	currentLeader int
	//该节点最后一次处理数据的时间
	lastHeartBeatFromLeader int64
	message                 chan bool
	electCh                 chan bool
	heartbeat               chan bool
	//子节点给主节点返回心跳信号
	heartbeatRe chan bool
}

//声明leader对象
type Leader struct {
	//任期
	Term int
	//leader 编号
	LeaderId int
}

//设置节点个数
const raftCount = 2

var leader = Leader{0, -1}

//存储缓存信息
var bufferMessage = make(map[string]string)

//处理数据库信息
var mysqlMessage = make(map[string]string)

//操作消息数组下标
var messageId = 1

//用nodeTable存储每个节点中的键值对
var nodeTable map[string]string

func main() {
	//终端接收来的是数组
	if len(os.Args) > 1 {
		//接收终端输入的信息
		userId := os.Args[1]
		//字符串转换整型
		id, _ := strconv.Atoi(userId)
		fmt.Println(id)

		/**
			  定义节点id和端口号
		         "节点id"："ip:port"
			eg. "1":"192.168.12.2:8080"
			也可以写成配置文件读进来
		*/
		nodeTable = map[string]string{
			"1": ":8000",
			"2": ":8001",
		}
		//封装nodeInfo对象
		node := nodeInfo{id: userId, port: nodeTable[userId]}
		//创建节点对象
		rf := Make(id)
		//确保每个新建立的节点都有端口对应
		//127.0.0.1:8000
		rf.node = node
		//注册rpc
		go func() {
			//注册rpc，为了实现远程链接
			rf.raftRegisterRPC(node.port)
		}()
		if userId == "1" {
			go func() {
				//回调方法
				http.HandleFunc("/req", rf.getRequest)
				fmt.Println("监听8080")
				if err := http.ListenAndServe(":8080", nil); err != nil {
					fmt.Println(err)
					return
				}
			}()
		}
	}
	for {
	}
}

var clientWriter http.ResponseWriter

func (rf *Raft) getRequest(writer http.ResponseWriter, request *http.Request) {
	_ = request.ParseForm()
	if len(request.Form["age"]) > 0 {
		clientWriter = writer
		fmt.Println("主节点广播客户端请求age:", request.Form["age"][0])

		param := Param{Msg: request.Form["age"][0], MsgId: strconv.Itoa(messageId)}
		messageId++
		if leader.LeaderId == rf.me {
			rf.sendMessageToOtherNodes(param)
		} else {
			/*
					将消息转发给leader 这步是因为我写死了sever 1 为接受消息的server 所以需要
				转发，raft论文里实际上是All changes to the system now go through the leader.
				其实这里的实现有些问题，有机会重写一下
			*/
			leaderId := nodeTable[strconv.Itoa(leader.LeaderId)]
			//连接远程rpc服务
			Rpc, err := rpc.DialHTTP("tcp", "127.0.0.1"+leaderId)
			if err != nil {
				log.Fatal("\nrpc转发连接server错误:", leader.LeaderId, err)
			}
			var bo = false
			//首先给leader传递
			err = Rpc.Call("Raft.ForwardingMessage", param, &bo)
			if err != nil {
				log.Fatal("\nrpc转发调用server错误:", leader.LeaderId, err)
			}
		}
	}
}

func (rf *Raft) sendMessageToOtherNodes(param Param) {
	bufferMessage[param.MsgId] = param.Msg
	// 只有领导才能给其它服务器发送消息
	if rf.currentLeader == rf.me {
		var successCount = 0
		fmt.Printf("领导者发送数据中 。。。\n")
		go func() {
			rf.broadcast(param, "Raft.LogDataCopy", func(ok bool) {
				//需要其它服务端回应
				rf.message <- ok
			})
		}()

		for i := 0; i < raftCount-1; i++ {
			fmt.Println("等待其它服务端回应")
			select {
			case ok := <-rf.message:
				if ok {
					successCount++
					if successCount >= raftCount/2 {
						rf.mu.Lock()
						mysqlMessage[param.MsgId] = bufferMessage[param.MsgId]
						delete(bufferMessage, param.MsgId)
						if clientWriter != nil {
							_, _ = fmt.Fprintf(clientWriter, "OK")
						}
						fmt.Printf("\n领导者发送数据结束\n")
						rf.mu.Unlock()
					}
				}
			}
		}
	}
}

//注册Raft对象，注册后的目的为确保每个节点（raft) 可以远程接收
func (rf *Raft) raftRegisterRPC(port string) {
	//注册一个服务器
	_ = rpc.Register(rf)
	//把服务绑定到http协议上
	rpc.HandleHTTP()
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("注册rpc服务失败", err)
	}
}

//创建节点对象
func Make(me int) *Raft {
	rf := &Raft{}
	rf.me = me
	rf.votedFor = -1
	//0 follower ,1 candidate ,2 leader
	rf.state = 0
	rf.timeout = 0
	rf.currentLeader = -1
	rf.setTerm(0)

	//初始化通道
	rf.message = make(chan bool)
	rf.heartbeat = make(chan bool)
	rf.heartbeatRe = make(chan bool)
	rf.electCh = make(chan bool)

	//每个节点都有选举权
	go rf.election()
	//每个节点都有心跳功能
	go rf.sendLeaderHeartBeat()

	return rf
}

//leader选举成功后，应该广播所有的节点，本节点成为了leader
//顺便完成数据同步
//看其他节点挂没挂
func (rf *Raft) sendLeaderHeartBeat() {
	for {
		select {
		case <-rf.heartbeat:
			rf.sendAppendEntriesImpl()
		}
	}
}

func (rf *Raft) sendAppendEntriesImpl() {
	//是leader
	if rf.currentLeader == rf.me {
		//记录确认节点个数
		var successCount = 0
		go func() {
			param := Param{Msg: "leader heartbeat",
				Arg: Leader{rf.currentTerm, rf.me}}
			rf.broadcast(param, "Raft.Heartbeat", func(ok bool) {
				//写
				rf.heartbeatRe <- ok
			})
		}()
		for i := 0; i < raftCount-1; i++ {
			select {
			//收 并统计heartbeatRe
			case ok := <-rf.heartbeatRe:
				if ok {
					successCount++
					if successCount >= raftCount/2 {
						rf.mu.Lock()
						//	rf.lastMessageTime = milliseconds()
						fmt.Println("接收到了子节点们的返回信息")
						rf.mu.Unlock()
					}
				}
			}
		}
	}
}

func timeSleepRandomRange(min, max int64) int64 {
	//设置随机时间
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(max-min) + min
}

//获得当前时间（毫秒）
func milliseconds() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func (rf *Raft) election() {
	//标记是否选出了leader
	var result bool
	//每隔一段时间发一次心跳  死循环 不成功就一直选举
	for {
		//延时时间
		timeout := timeSleepRandomRange(1500, 3000)

		//超过十秒没有收到leader的心跳 则转换为candidate继续选举
		if rf.lastHeartBeatFromLeader <= milliseconds()+10000 {

			select {
			//间隔时间为1500-3000ms的随机值
			case <-time.After(time.Duration(timeout) * time.Millisecond):
				fmt.Println("当前节点状态：", rf.state)
			}

			result = false
			for !result {
				//选择leader
				result = rf.electionOneRound(&leader)
			}
		}
	}
}

func (rf *Raft) electionOneRound(args *Leader) bool {
	//已经有了leader，并且不是自己，那么return
	if args.LeaderId > -1 && args.LeaderId != rf.me {
		fmt.Printf("%d已是leader，终止%d选举\n", args.LeaderId, rf.me)
		return true
	}

	var timeout int64
	var vote int
	//是否开始产生心跳信号
	var triggerHeartbeat bool
	timeout = 2000
	last := milliseconds()
	success := false
	rf.mu.Lock()
	rf.becomeCandidate()
	rf.mu.Unlock()
	fmt.Printf("candidate=%d start electing leader\n", rf.me)
	for {
		fmt.Printf("candidate=%d send request vote to server\n", rf.me)
		go func() {
			rf.broadcast(Param{Msg: "send request vote"}, "Raft.ElectingLeader", func(ok bool) {
				//无论成功失败都需要发送到通道 避免堵塞
				rf.electCh <- ok
			})
		}()

		vote = 0
		triggerHeartbeat = false
		for i := 0; i < raftCount-1; i++ {
			fmt.Printf("candidate=%d waiting for select for i=%d\n", rf.me, i)
			//计算投票数量
			select {
			case ok := <-rf.electCh:
				if ok {
					vote++
					success = vote >= raftCount/2 || rf.currentLeader > -1
					if success && !triggerHeartbeat {
						//选主成功，触发心跳信号检测
						fmt.Println("okok", args)
						triggerHeartbeat = true
						rf.mu.Lock()
						rf.becomeLeader()
						args.Term = rf.currentTerm + 1
						args.LeaderId = rf.me
						rf.mu.Unlock()
						fmt.Printf("candidate=%d becomes leader\n", rf.currentLeader)
						//leader 向其他节点发送信号
						rf.heartbeat <- true
					}
				}
			}
			fmt.Printf("candidate=%d complete for select for i=%d\n", rf.me, i)
		}
		//校验 不超时 且票数大于一半 则成功
		if (timeout+last < milliseconds()) || (vote >= raftCount/2 || rf.currentLeader > -1) {
			break
		} else {
			//等一会 继续选举
			select {
			case <-time.After(time.Duration(5000) * time.Millisecond):
			}
		}
	}
	fmt.Printf("candidate=%d receive votes status=%t\n", rf.me, success)
	return success
}

func (rf *Raft) becomeLeader() {
	rf.state = 2
	fmt.Println(rf.me, "成为了leader")
	rf.currentLeader = rf.me
}

//设置发送参数的数据类型  RPC
type Param struct {
	Msg   string
	MsgId string
	Arg   Leader
}

//rpc
func (rf *Raft) broadcast(msg Param, path string, fun func(ok bool)) {
	//设置不要自己给自己广播
	for nodeID, port := range nodeTable {
		if nodeID == rf.node.id {
			continue
		}
		//链接远程rpc
		rp, err := rpc.DialHTTP("tcp", "127.0.0.1"+port)
		if err != nil {
			fun(false)
			continue
		}
		//接受返回信息 bo
		var bo = false
		err = rp.Call(path, msg, &bo)
		if err != nil {
			fun(false)
			continue
		}
		fun(bo)
	}
}

func (rf *Raft) becomeCandidate() {
	if rf.state == 0 || rf.currentLeader == -1 {
		//候选人状态
		rf.state = 1
		rf.votedFor = rf.me
		rf.setTerm(rf.currentTerm + 1)
		rf.currentLeader = -1

	}
}

func (rf *Raft) setTerm(term int) {
	rf.currentTerm = term
}

//Rpc处理
func (rf *Raft) ElectingLeader(param Param, a *bool) error {
	//给leader投票
	fmt.Println(param)
	*a = true
	//	rf.lastMessageTime = milliseconds()
	return nil
}

func (rf *Raft) Heartbeat(param Param, a *bool) error {
	fmt.Println("\nrpc:heartbeat:", rf.me, param.Msg)
	if param.Arg.Term < rf.currentTerm {
		*a = false
	} else {
		leader = param.Arg
		fmt.Printf("%d收到leader%d心跳\n", rf.currentLeader, leader.LeaderId)
		*a = true
		rf.mu.Lock()
		rf.currentLeader = leader.LeaderId
		rf.votedFor = leader.LeaderId
		rf.lastHeartBeatFromLeader = milliseconds()
		rf.state = 0
		//	rf.lastMessageTime = milliseconds()
		fmt.Printf("server = %d learned that leader = %d\n", rf.me, rf.currentLeader)
		rf.mu.Unlock()
	}
	return nil
}

//连接到leader节点
func (rf *Raft) ForwardingMessage(param Param, a *bool) error {
	fmt.Println("\nrpc:forwardingMessage:", rf.me, param.Msg)

	rf.sendMessageToOtherNodes(param)

	*a = true
	//	rf.lastMessageTime = milliseconds()

	return nil
}

//接收leader传过来的日志
func (rf *Raft) LogDataCopy(param Param, a *bool) error {
	fmt.Println("\nrpc:LogDataCopy:", rf.me, param.Msg)
	bufferMessage[param.MsgId] = param.Msg
	*a = true
	return nil
}
