package main

import "net"

type user struct {
	name string      // 用户姓名
	addr string      // IP地址
	c    chan string // 接收信息通道
	conn net.Conn    // C端和S端建立的socket连接
}

func NewUser(conn net.Conn) *user {
	userAddress := conn.RemoteAddr().String() // 获取ip地址
	u := &user{
		name: userAddress,
		addr: userAddress,
		c:    make(chan string, 1),
		conn: conn,
	}
	go u.ListenMessage() // 客户端开始监听通道是否有信息

	return u
}

func (u *user) ListenMessage() {
	for {
		msg := <-u.c                     // 当user的通道有信息，则取出
		u.conn.Write([]byte(msg + "\n")) // 将信息通过socket发送给Client端
	}
}
