package main

import "net"

type User struct {
	Name string
	Addr string
	UsrChan    chan string
	conn net.Conn
}

// 新建一个用户的接口
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User {
		Name: userAddr,
		Addr: userAddr,
		UsrChan: make(chan string),
		conn: conn,
	}

	// 启动一个监听当前user channel消息的goroutine
	go user.ListenMsg()

	return user
}

// 监听当前User的channel的方法，一旦有新消息就发送给
func (user *User) ListenMsg() {
	for {
		msg := <-user.UsrChan
		user.conn.Write([]byte(msg + "\n"))
	}
}