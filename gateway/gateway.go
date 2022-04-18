package gateway

import (
	"flag"
	"fmt"
	"io"
	"net"
	"shield/utils"
	"time"
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

var hasRetryTimes = 0

func handleClientRequest(client net.Conn, remoteIp string) {
	//读取请求的内容
	buf := make([]byte, 1024)
	n, e := client.Read(buf)
	if e != nil {
		//todo by@老王 客户端建立了连接，但未发送任何数据，此时我们每隔2秒累计重试3此
		// todo 从io.copy代码测试来看 游戏是进入时发送请求,然后立即断开整条连接,最后登录后再保持长连接
		if e == io.EOF {
			fmt.Println("读取client数据为空")
			if hasRetryTimes > 3 {
				hasRetryTimes++
				time.Sleep(2 * time.Second)
				handleClientRequest(client, remoteIp)
			} else {
				return
			}
		}
		fmt.Println("读取balancer数据异常", e)
	}
	fmt.Println("收到来自balancer的数据,大小为:", n)
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

	fmt.Println("转发数据")
	go utils.TcpRequest(client, remote)
	go utils.TcpRequest(remote, client)
}
