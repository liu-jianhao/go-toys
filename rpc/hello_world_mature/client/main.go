package main

import (
	"fmt"

	"github.com/liu-jianhao/go-toys/rpc/hello_world_mature/protocol"
)

func main() {
	client, err := protocol.DialHelloWorldService("tcp", ":12345")
	if err != nil {
		panic("rpc dial error")
	}

	var reply string
	err = client.HelloWorld("world", &reply)
	if err != nil {
		panic("call rpc service error")
	}

	fmt.Println(reply)
}
