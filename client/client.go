package client

import (
	"fmt"
	"net"
	"shield/utils"
)

func Init() {
	// 监听端口
	server, err := net.Listen("tcp", ":3153")
	if err != nil {
		panic(err)
		return
	}
	fmt.Println("启动监听 3153 端口")
	for {
		//收到请求
		client, err := server.Accept()
		if err != nil {
			fmt.Println(err)
			break
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

	//和balance建立建立连接
	remote, err := net.Dial("tcp", "152.67.217.198:51777")
	if err != nil {
		err := remote.Close()
		if err != nil {
			panic(e)
		}
	}
	fmt.Println("连接到 balance")

	// 把数据写入到balance
	n, e = remote.Write(buf[:n])
	if e != nil {
		panic(e)
	}

	//todo 将balcner发回来的数据写入到游戏客户端
	fmt.Println("转发数据")
	go utils.TcpRequest(client, remote)
	go utils.TcpRequest(remote, client)
}
