package balancer

import (
	"fmt"
	"net"
	"shield/utils"
)

func Init() {
	// 监听端口
	server, err := net.Listen("tcp", ":51777")
	if err != nil {
		panic(err)
		return
	}
	fmt.Println("启动监听 51777 端口")
	for {
		//收到请求
		client, err := server.Accept()
		if err == nil {
			go handleClientRequest(client)
		}
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
	fmt.Println("准备与gateway建立连接发送数据")

	//1.和balance建立建立连接
	remote, err := net.Dial("tcp", "101.43.218.210:51777")
	if err != nil {
		panic(err)
	}
	fmt.Println("has dial to balance")

	// 把数据写入到balance
	//todo  如果有bug 则考虑把下面三行去掉
	n, e = remote.Write(buf[:n])
	if e != nil {
		panic(e)
	}

	//todo 将gateway发回来的数据写入到client

	fmt.Println("转发数据")
	// proxyconn - 游戏客户端连接
	// targetconn  - balance连接
	go utils.ProxyRequest(client, remote)
	go utils.ProxyRequest(remote, client)
}
