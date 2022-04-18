package balancer

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
	remoteip := flag.String("r", "119.188.247.208:40001", "remote")
	flag.Parse()
	server, err := net.Listen("tcp", *listenip)
	if err != nil {
		fmt.Println("接收Client数据异常", err)
		return
	}
	fmt.Println(*listenip, *remoteip)
	for {
		//收到请求
		client, err := server.Accept()
		if err != nil {
			fmt.Println("h获取y游戏客户端l连接y异常,", err)
			break
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
		fmt.Println("读取client数据异常", e)
	}
	fmt.Println("收到来自游戏客户端的数据,大小为:", n)
	fmt.Println("准备与gateway建立连接发送数据")

	//和gateway建立建立连接
	remote, err := net.Dial("tcp", remoteIp)
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

	fmt.Println("转发数据")
	go utils.TcpRequest(client, remote)
	go utils.TcpRequest(remote, client)
}
