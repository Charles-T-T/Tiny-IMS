package main

import (
	"fmt"
	"net"
)

type Server struct {
	IP   string
	Port int
}

// 新建一个服务器的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		IP:   ip,
		Port: port,
	}
	return server
}

func (server *Server) Handler(conn net.Conn) {
	// todo 当前连接的业务
	fmt.Println("Connection built successfully!")
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
