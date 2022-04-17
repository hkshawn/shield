package gateway

import (
	"flag"
	"fmt"
	"io"
	"net"
	"shield/utils"
)

func Init() {
	// 监听端口
	listenip := flag.String("l", "0.0.0.0:40000", "listen")
	remoteip := flag.String("r", "127.0.0.1:40000", "remote")
	flag.Parse()
	server, err := net.Listen("tcp", *listenip)
	if err != nil {
		fmt.Println("监听异常,err:", err)
		return
	}
	fmt.Println(*listenip, *remoteip)
	for {
		//收到请求
		client, err := server.Accept()
		if err != nil {
			fmt.Println("创建连接异常", err)
			return
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
			fmt.Println("读取balancer数据为空")
			return
		}
		fmt.Println("读取balancer数据异常", e)
		return
	}
	fmt.Println("收到来自游戏客户端的数据,大小为:", n)
	fmt.Println("准备与game-server建立连接发送数据")

	//和game-server建立建立连接
	remote, err := net.Dial("tcp", remoteIp)
	if err != nil {
		fmt.Println("连接错误:", err)
		return
	}
	fmt.Println("连接到 game-server")

	// 把数据写入到game-server
	n, e = remote.Write(buf[:n])
	if e != nil {
		fmt.Println("写入到balancer错误", err)
		return
	}

	//todo 将game-server发回来的数据写入到balancer
	fmt.Println("转发数据")
	go utils.TcpRequest(client, remote)
	go utils.TcpRequest(remote, client)
}
