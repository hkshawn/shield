package utils

import (
	"fmt"
	"io"
	"net"
	"time"
)

//tcp数据转发
func TcpRequest(r net.Conn, w net.Conn) {
	defer func() {
		if e := r.Close(); e != nil {
			fmt.Println("关闭读取连接,异常:", e)
			return
		}
	}()
	defer func() {
		if e := w.Close(); e != nil {
			fmt.Println("关闭写入连接异常:", e)
			return
		}
	}()

	// todo 把panic(err)修改为break跳出循环,修复了崩溃BUG
	var buffer = make([]byte, 4096000)

	for {
		if e := r.SetDeadline(time.Now().Add(time.Second * 20)); e != nil {
			fmt.Println("设置读取连接超时时间异常:", e)
		}
		if e := w.SetDeadline(time.Now().Add(time.Second * 20)); e != nil {
			fmt.Println("设置写入连接超时时间异常:", e)
		}

	L:
		{
			n, err := r.Read(buffer)
			if err != nil {
				if err == io.EOF {
					//fmt.Println("等待数据中---------------------------------------------------------------------")
					goto L
				}
				fmt.Println(err)
				break
			}
			//fmt.Println("读取成功,大小:", n)
			n, err = w.Write(buffer[:n])
			if err != nil {
				fmt.Println(err)
				break
			}
			//fmt.Println("写入成功,大小:", n)
		}
	}
}

type AppInfo struct {
	AppID string `json:"app_id"`
	IP    string `json:"ip"`
	Port  string `json:"port"`
}

//todo
//1. 再tcp header中加上appID 端口号
//2. 再balancer解析tcp header,拿到appid,选择转发到哪台gateway
//3. gateway解析header,选择转发到b本机的哪个端口号
//*. 对body进行加密后在发送

/*
panic: runtime error: invalid memory address or nil pointer dereference
	panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x18 pc=0xd6f30]

goroutine 5 [running]:
shield/balancer.handleClientRequest.func2({0x0, 0x0})
	/Users/rango/go/src/shield/balancer/balancer.go:55 +0x20
panic({0xef060, 0x1c10a0})
	/usr/local/go/src/runtime/panic.go:1038 +0x224
shield/balancer.handleClientRequest({0x12b6e8, 0x400000e028})
	/Users/rango/go/src/shield/balancer/balancer.go:64 +0x25c
created by shield/balancer.Init
	/Users/rango/go/src/shield/balancer/balancer.go:29 +0x11c
*/
