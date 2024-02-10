package server

import (
	"fmt"
	"net"
)

type Server struct {
	ip   string
	port int
}

func NewServer(ip string, port int) *Server {
	return &Server{ip, port}
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%v:%v", s.ip, s.port))
	if err != nil {
		fmt.Printf("Server Listen is failed! Error:%v\n", err)
		return
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Server Accept is failed! Error:%v\n", err)
			continue
		}
		go s.Handler(conn)
	}
}

func (s *Server) Handler(conn net.Conn) {
	fmt.Printf("连接成功，开始处理!\n")
}
