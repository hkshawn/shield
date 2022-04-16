package balancer

import (
	"fmt"
	"io"
	"net"
	"shield/utils"
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
	fmt.Println("准备与gateway建立连接发送数据")

	//和gateway建立建立连接
	remote, err := net.Dial("tcp", "101.43.218.210:51777")
	// todo 修复连接关闭后导致的程序崩溃
	if err != nil {
		defer func() {
			closeErr := remote.Close()
			err := remote.Close()
			if err == nil {
				err = closeErr
			}
		}()
	}
	fmt.Println("连接到 gateway")

	// 把数据写入到gateway
	// todo 关闭gateway后这里还是会导致balancer崩溃
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

	//todo 将gateway发回来的数据写入到client
	fmt.Println("转发数据")
	go utils.TcpRequest(client, remote)
	go utils.TcpRequest(remote, client)
}
