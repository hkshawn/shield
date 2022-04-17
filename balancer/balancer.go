package balancer

import (
	"fmt"
	"net"
	"shield/utils"
	"time"
)

func Init() {
	// 监听端口
	server, err := net.Listen("tcp", ":51777")
	if err != nil {
		err := server.Close()
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("启动监听 51777 端口")
	for {
		//收到请求
		client, err := server.Accept()
		if err != nil {
			fmt.Println("h获取y游戏客户端l连接y异常,", err)
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
	fmt.Println("准备与gateway建立连接发送数据")

	//和gateway建立建立连接
	remote, err := net.DialTimeout("tcp", "101.43.218.210:51777", 2*time.Second)
	if err != nil {
		fmt.Println("与gateway建立连接异常", err)
		return
	}
	fmt.Println("连接到 gateway")

	// 把数据写入到gateway
	n, e = remote.Write(buf[:n])
	if e != nil {
		fmt.Println("数据写入gateway异常", err)
		return
	}

	//todo 将gateway发回来的数据写入到client
	fmt.Println("转发数据")
	go utils.TcpRequest(client, remote)
	go utils.TcpRequest(remote, client)
}
