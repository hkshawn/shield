package client

import (
	"flag"
	"fmt"
	"io"
	"net"
	"shield/utils"
)

func Init() {
	// 监听端口
	listenip := flag.String("l", "127.0.0.2:40000", "listen")
	remoteip := flag.String("r", "49.232.26.65:40000", "remote")
	flag.Parse()
	server, err := net.Listen("tcp", *listenip)
	if err != nil {
		fmt.Println("接收游戏数据异常", err)
		return
	}
	fmt.Println(*listenip, *remoteip)
	for {
		//收到请求
		client, err := server.Accept()
		if err != nil {
			fmt.Println("接收游戏客户端数据异常", err)
			break
		}
		go handleClientRequest(client, *remoteip)
	}
}

func handleClientRequest(client net.Conn, remoteIp string) {
	//读取请求的内容
	buf := make([]byte, 1024)
	n, e := client.Read(buf)
	if e != nil {
		if e == io.EOF {
			fmt.Println("读取游戏客户端数据为空")
			return
		}
		fmt.Println("读取游戏客户端数据异常", e)
	}
	fmt.Println("收到来自游戏客户端的数据,大小为:", n)
	fmt.Println("准备与balancer建立连接发送数据")

	//和balance建立建立连接
	remote, err := net.Dial("tcp", remoteIp)
	if err != nil {
		fmt.Println("与balancer建立连接异常", err)
		return
	}
	fmt.Println("连接到 balance")

	// 把数据写入到balance
	n, e = remote.Write(buf[:n])
	if e != nil {
		fmt.Println("写入balancer异常", err)
		return
	}

	//todo 将balcner发回来的数据写入到游戏客户端
	fmt.Println("转发数据")
	go utils.TcpRequest(client, remote)
	go utils.TcpRequest(remote, client)
}
