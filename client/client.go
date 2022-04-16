package client

import (
	"fmt"
	"io"
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
		// todo 修复错误的连接地址导致的程序崩溃
		client, err := server.Accept()
		if err != nil {
			fmt.Println("连接出错")
		}
		go handleClientRequest(client)
	}
}

func handleClientRequest(client net.Conn) {
	//读取请求的内容
	// todo 修复连接关闭后导致的程序崩溃
	buf := make([]byte, 1024)
	n, e := client.Read(buf)
	if e != nil {
		defer func() {
			closeErr := client.Close()
			err := client.Close()
			if err == nil {
				err = closeErr
			}
		}()
	}
	fmt.Println("收到来自游戏客户端的数据,大小为:", n)
	fmt.Println("准备与balancer建立连接发送数据")

	//和balance建立建立连接
	remote, err := net.Dial("tcp", "152.67.217.198:51777")
	// todo 修复错误的连接地址导致的程序崩溃
	if err != nil {
		defer func() {
			closeErr := remote.Close()
			err := remote.Close()
			if err == nil {
				err = closeErr
			}
		}()
	}
	fmt.Println("连接到 balance")

	// 把数据写入到balance
	n, e = remote.Write(buf[:n])
	if e != nil {
		defer func() {
			closeErr := remote.Close()
			err := remote.Close()
			if err == nil {
				err = closeErr
			}
		}()
		//todo 增加这里后 好像解决了中途掉线的问题 还有刚开始连接重复掉线
		_, err := io.Copy(client, remote)
		if err != nil {
			panic(err)
		}
	}

	//todo 将balcner发回来的数据写入到游戏客户端
	fmt.Println("转发数据")
	go utils.TcpRequest(client, remote)
	go utils.TcpRequest(remote, client)
}
