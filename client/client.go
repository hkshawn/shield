package client

import (
	"fmt"
	"net"
	"shield/utils"
)

func Init() {
	// 监听端口
	server, err := net.Listen("tcp", ":41999")
	if err != nil {
		panic(err)
		return
	}
	fmt.Println("启动监听 51777 端口")
	for {
		//收到请求
		client, err := server.Accept()
		if err != nil {
			panic(err)
		}
		go handleClientRequest(client)
	}
}

func handleClientRequest(client net.Conn) {
	//读取请求的内容
	buf := make([]byte, 1024)
	n, e := client.Read(buf)
	if e != nil {
		panic(e)
	}
	fmt.Println("收到来自游戏客户端的数据,大小为:", n)
	fmt.Println("准备与balancer建立连接发送数据")

	//1.和balance建立建立连接
	remote, err := net.Dial("tcp", "152.67.217.198:51777")
	if err != nil {
		panic(err)
	}
	fmt.Println("has dial to balance")

	// 把数据写入到bbalance
	//todo  如果有bug 则考虑把下面三行去掉
	n, e = remote.Write(buf[:n])
	if e != nil {
		panic(e)
	}
	//todo 将balcner发回来的数据写入到游戏客户端
	fmt.Println("转发数据")
	// proxyconn - 游戏客户端连接
	// targetconn  - balance连接
	go utils.ProxyRequest(client, remote)
	go utils.ProxyRequest(remote, client)
}
