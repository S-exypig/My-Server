package main

import (
	"fmt"
	"net"
	"strings"
)

type user struct {
	name   string      // 用户姓名
	addr   string      // IP地址
	ch     chan string // 接收信息通道
	conn   net.Conn    // C端和S端建立的socket连接
	server *server     // 访问server类的句柄
}

func NewUser(conn net.Conn, s *server) *user {
	userAddress := conn.RemoteAddr().String() // 获取ip地址
	u := &user{
		name:   userAddress,
		addr:   userAddress,
		ch:     make(chan string, 1),
		conn:   conn,
		server: s,
	}
	go u.ListenMessage() // 客户端开始监听通道是否有信息

	return u
}

func (u *user) ListenMessage() {
	for {
		msg := <-u.ch      // 当user的通道有信息，则取出
		u.SendMessage(msg) 
	}
}

func (u *user) SendMessage(msg string) {
	// 将信息通过socket发送给Client端
	u.conn.Write([]byte(msg + "\n"))
}

func (u *user) Online() {
	u.server.mapSync.Lock()
	u.server.onlineMap[u.name] = u
	u.server.mapSync.Unlock()
	u.server.Broadcast(u, "已上线")
}

func (u *user) Offline() {
	u.server.mapSync.Lock()
	delete(u.server.onlineMap, u.name)
	u.server.mapSync.Unlock()
	u.server.Broadcast(u, "已下线")
}

func (u *user) DoMessage(msg string) {
	if msg == "who" {
		// 查询在线用户
		u.server.mapSync.RLock()
		for _, v := range u.server.onlineMap {
			m := fmt.Sprintf("[%v]%v:在线...\n", v.addr, v.name)
			u.SendMessage(m)
		}
		u.server.mapSync.RUnlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// 重命名：rename|xxx
		newName := strings.Split(msg, "|")[1]
		u.server.mapSync.RLock()
		_, exist := u.server.onlineMap[newName]
		u.server.mapSync.RUnlock()
		if exist {
			m := fmt.Sprintf("用户名%v已存在!", newName)
			u.SendMessage(m)
		} else {
			u.server.mapSync.Lock()
			delete(u.server.onlineMap, u.name)
			u.server.onlineMap[newName] = u
			u.server.mapSync.Unlock()
			u.name = newName
			u.SendMessage(fmt.Sprintf("您的用户名已更新:%v", newName))
		}
	} else {
		u.server.Broadcast(u, msg)
	}
}
