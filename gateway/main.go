package gateway

import (
	"fmt"
	"io"
	"net"
)

func Init() {
	server, err := net.Listen("tcp", ":51777")
	if err != nil {
		panic(err)
		return
	}
	fmt.Println("启动监听 51777 端口")
	for {
		client, err := server.Accept()
		if err == nil {
			go handleClientRequest(client)
		}
	}
}

func handleClientRequest(client net.Conn) {
	defer func(client net.Conn) {
		err := client.Close()
		if err != nil {
			panic(err)
		}
	}(client)

	remote, err := net.Dial("tcp", "127.0.0.1:3389")
	if err != nil {
		return
	}
	defer func(remote net.Conn) {
		err := remote.Close()
		if err != nil {
			panic(err)
		}
	}(remote)

	go func() {
		_, err := io.Copy(remote, client)
		if err != nil {
			panic(err)
		}
	}()
	_, err = io.Copy(client, remote)
	if err != nil {
		panic(err)
		return
	}
}
