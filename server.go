package main

import (
	"fmt"
	"net"
	"sync"
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
	u := NewUser(conn)
	s.mapSync.Lock()
	s.onlineMap[u.name] = u
	s.mapSync.Unlock()
	s.SendMessage(u, "已上线")
}

func (s *server) SendMessage(u *user, msg string) {
	s.message <- fmt.Sprintf("[%v]%v:%v", u.addr, u.name, msg)
}

func (s *server) ListenMessage() {
	fmt.Println("开始监听用户上线")
	for {
		msg := <-s.message
		s.mapSync.Lock()
		for _, v := range s.onlineMap {
			v.c <- msg
		}
		s.mapSync.Unlock()
	}
}
