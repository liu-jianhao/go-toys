package main

import (
	"net"
	"net/rpc"

	"github.com/liu-jianhao/go-toys/rpc/hello_world_mature/protocol"
)

func main() {
	err := protocol.RegisterHelloWorldService(new(HelloWorldService))
	if err != nil {
		panic("register error")
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
