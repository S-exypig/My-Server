package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "localhost", "设置客户端所要连接的服务器的ip")
	flag.IntVar(&serverPort, "port", 8000, "设置客户端所要连接的服务器的端口")
}

func main() {
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>> 连接服务器失败...")
		return
	}
	fmt.Println(">>>>> 连接服务器成功~~")
	go client.DealResponse()
	client.Run()
}

type client struct {
	serverIp   string
	serverPort int
	// name       string
	conn net.Conn
	mode int
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
	// client.name = conn.LocalAddr().String()
	// 返回指针
	return client
}

func (c *client) Menu() bool {
	fmt.Println("---------\n0.退出")
	fmt.Println("1.私聊模式")
	fmt.Println("2.公聊模式")
	fmt.Println("3.修改用户名")
	fmt.Println("4.查看在线用户")
	var mode int
	_, err := fmt.Scanln(&mode)
	if err != nil || mode < 0 || mode > 4 {
		fmt.Println("输入有误!")
		return false
	}
	c.mode = mode
	return true
}

func (c *client) Run() {
LOOP:
	for {
		for c.Menu() != true {
		}
		switch c.mode {
		case 0:
			fmt.Println("退出!")
			break LOOP
		case 1:
			fmt.Println("私聊模式...")
		case 2:
			c.PublicChat()
		case 3:
			if c.UpdateName() {
				fmt.Println("改名成功!")
			} else {
				fmt.Println("改名未成功!")
			}
		case 4:
			fmt.Println("查看在线用户...")
		}
	}
}

func (c *client) UpdateName() bool {
	fmt.Println(">>>>>请输入用户名:")
	var name string
	// fmt.Scanln(&c.name)
	// sendMsg := fmt.Sprintf("rename|%v\n", c.name)
	fmt.Scanln(&name)
	sendMsg := fmt.Sprintf("rename|%v\n", name)
	_, err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Printf("conn.Write Error:%v\n", err)
		return false
	}
	return true
}

// 处理S端回应的信息，直接打印到标准输出
func (c *client) DealResponse() {
	// 一旦客户端的conn收到数据，直接copy到标准输出流
	io.Copy(os.Stdout, c.conn)
}

func (c *client) PublicChat() {
	var chatMsg string
	fmt.Println(">>>>>公聊模式--请输入聊天内容，exit退出.")
	fmt.Scanln(&chatMsg)
	for chatMsg != "exit" {
		if chatMsg != "" {
			_, err := c.conn.Write([]byte(chatMsg + "\n"))
			chatMsg = ""
			if err != nil {
				fmt.Printf("Conn.Write Error:%v\n", err)
				break
			}
		}
		fmt.Println(">>>>>公聊模式--请输入聊天内容，exit退出.")
		fmt.Scanln(&chatMsg)
	}
}
