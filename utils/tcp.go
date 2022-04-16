package utils

import (
	"fmt"
	"net"
)

//tcp数据转发
func TcpRequest(r net.Conn, w net.Conn) {
	defer func() {
		//todo 优化defer 连接关闭后不会导致程序崩溃
		closeErr := r.Close()
		err := r.Close()
		if err == nil {
			err = closeErr
		}
	}()
	defer func() {
		//todo 优化defer 连接关闭后不会导致程序崩溃
		closeErr := w.Close()
		err := w.Close()
		if err == nil {
			err = closeErr
		}
	}()

	// todo 把panic(err)修改为break跳出循环,修复了崩溃BUG
	var buffer = make([]byte, 4096000)
	for {
		n, err := r.Read(buffer)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println("读取成功,大小:", n)
		n, err = w.Write(buffer[:n])
		if err != nil {
			fmt.Println(err)
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Println("写入成功,大小:", n)
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
