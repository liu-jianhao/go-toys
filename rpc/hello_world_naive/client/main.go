package main

import (
	"fmt"
	"net/rpc"
)

func main() {
	client, err := rpc.Dial("tcp", ":12345")
	if err != nil {
		panic("rpc dial error")
	}

	var reply string
	err = client.Call("HelloWorldService.HelloWorld", "world", &reply)
	if err != nil {
		panic("call rpc service error")
	}

	fmt.Println(reply)
}
