package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type server struct {
	ip        string           // 服务端ip地址
	port      int              // 服务端端口号
	onlineMap map[string]*user // 在线的用户
	mapSync   sync.RWMutex     // map的锁，map不是线程安全的
	message   chan string      // 发送信息通道
}

func NewServer(ip string, port int) *server {
	return &server{
		ip:        ip,
		port:      port,
		onlineMap: make(map[string]*user),
		message:   make(chan string),
	}
}

func (s *server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%v:%v", s.ip, s.port))
	if err != nil {
		fmt.Printf("Server Listen is failed! Error:%v\n", err)
		return
	}
	defer listener.Close()
	go s.ListenMessage() // S端监听用户上线，发送上线信息给C端
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Server Accept is failed! Error:%v\n", err)
			continue
		}
		go s.Handler(conn)
	}
}

func (s *server) Handler(conn net.Conn) {
	fmt.Printf("连接成功，开始处理!\n")
	u := NewUser(conn, s)
	u.Online()
	isAliveCh := make(chan bool)
	isnAliveCh := make(chan bool)
	go func() { // 接收C端发送来的数据
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				if err == io.EOF {
					u.Offline()
				} else {
					fmt.Printf("Conn Read Error:%v\n", err)
				}
				isnAliveCh <- true
				return
			}
			msg := string(buf[:n-1])
			u.DoMessage(msg) // 将信息交给S端的user模块进行处理
			isAliveCh <- true
		}
	}()

	for {
		select {
		case <-isAliveCh:
		case <-time.After(120 * time.Second):
			// 10s超时，强制下线
			u.SendMessage("你已被超时强踢！")
			u.Offline()
			return
		case <-isnAliveCh:
			return
		}
	}
}

func (s *server) Broadcast(u *user, msg string) {
	s.message <- fmt.Sprintf("[%v]%v:%v", u.addr, u.name, msg)
}

func (s *server) ListenMessage() {
	fmt.Println("开始监听用户上线")
	for {
		msg := <-s.message
		s.mapSync.RLock()
		for _, v := range s.onlineMap {
			v.ch <- msg
		}
		s.mapSync.RUnlock()
	}
}
