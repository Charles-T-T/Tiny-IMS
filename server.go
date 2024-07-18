package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	IP        string
	Port      int
	OnlineMap map[string]*User // 在线用户列表
	mapLock   sync.RWMutex
	MsgChan   chan string // 消息该广播的channel

}

// 新建一个服务器的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		IP:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		MsgChan:   make(chan string),
	}

	return server
}

// 监听message广播消息channel的goroutine，一旦有消息就发送给所有在线用户
func (server *Server) ListenMessager() {
	for {
		msg := <-server.MsgChan

		// 将msg发送给所有在线用户
		server.mapLock.Lock()
		for _, cli := range server.OnlineMap {
			cli.UsrChan <- msg
		}
		server.mapLock.Unlock()
	}
}

// 广播消息的方法
func (server *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	server.MsgChan <- sendMsg
}

func (server *Server) Handler(conn net.Conn) {
	// todo 当前连接的业务
	fmt.Println("Connection built successfully!")

	user := NewUser(conn)

	// 用户上线，将用户加入到OnlineMap中
	server.mapLock.Lock()
	server.OnlineMap[user.Name] = user
	server.mapLock.Unlock()

	// 广播当前用户上线信息
	server.BroadCast(user, "is online now")

	// 当前handler阻塞
	select {}
}

// 启动服务器的接口
func (server *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.IP, server.Port))
	if err != nil {
		fmt.Println("net.Listen error:", err)
		return
	}

	// close listen socket
	defer listener.Close()

	// 启动监听用户消息管道UsrChan的goroutine
	go server.ListenMessager()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept error:", err)
			continue
		}

		// do handler
		go server.Handler(conn)
	}
}
