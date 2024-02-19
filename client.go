package main

import (
	"fmt"
	"net"
)

type client struct {
	serverIp   string
	serverPort int
	name       string
	conn       net.Conn
}

func NewClient(ip string, port int) *client {
	// 创建client指针
	client := &client{
		serverIp:   ip,
		serverPort: port,
	}
	// 连接服务端
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Printf("net.Dial Error:%v\n", err)
		return nil
	}
	// 补全conn和name
	client.conn = conn
	client.name = conn.LocalAddr().String()
	// 返回指针
	return client
}

func main() {
	client := NewClient("localhost", 8000)
	if client == nil {
		fmt.Println(">>>>> 连接服务器失败...")
		return
	}
	fmt.Println(">>>>> 连接服务器成功~~")
	select {}
}
