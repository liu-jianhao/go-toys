package main

import (
	"net"
	"net/rpc"
)

func main() {
	err := rpc.RegisterName("HelloWorldService", new(HelloWorldService))
	if err != nil {
		panic("register rpc error")
	}

	listener, err := net.Listen("tcp", ":12345")
	if err != nil {
		panic("listen error")
	}

	conn, err := listener.Accept()
	if err != nil {
		panic("accept error")
	}

	rpc.ServeConn(conn)
}
