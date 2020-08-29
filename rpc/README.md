# RPC
RPC是远程过程调用的缩写，简单来说就是调用远程的一个函数。

Go语言的标准库中也有RPC包：`net/rpc`，下面就实现一个RPC版本的Hello World

## naive版本
先实现一个把协议相关的代码放到一起的版本。
### 服务端
service.go
```go
type HelloWorldService struct{}
func (s *HelloWorldService) HelloWorld(request string, reply *string) error {
	*reply = "hello:" + request
	return nil
}
```

main.go
```go
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
```

### 客户端
```go
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
```

## mature版本
涉及RPC的应用一般分为三部分：服务端实现RPC方法，客户端实现调用RPC方法，最后是服务端和客户端对该RPC方法的协议。

### 协议
```go
package protocol

import "net/rpc"

const HelloWorldServiceName = "HelloWorldService"

type HelloWorldServiceInterface interface {
	HelloWorld(request string, reply *string) error
}

func RegisterHelloWorldService(svc HelloWorldServiceInterface) error {
	return rpc.RegisterName(HelloWorldServiceName, svc)
}

type HelloWorldServiceClient struct {
	*rpc.Client
}

func DialHelloWorldService(network, address string) (*HelloWorldServiceClient, error) {
	c, err := rpc.Dial(network, address)
	if err != nil {
		return nil, nil
	}

	return &HelloWorldServiceClient{c}, nil
}

func (p *HelloWorldServiceClient) HelloWorld(request string, reply *string) error {
	return p.Client.Call(HelloWorldServiceName+".HelloWorld", request, reply)
}
```

### 服务端
```go
const HelloWorldServiceName = "HelloWorldService"

type HelloWorldServiceInterface interface {
	HelloWorld(request string, reply *string) error
}

func RegisterHelloWorldService(svc HelloWorldServiceInterface) error {
	return rpc.RegisterName(HelloWorldServiceName, svc)
}

type HelloWorldServiceClient struct {
	*rpc.Client
}

func DialHelloWorldService(network, address string) (*HelloWorldServiceClient, error) {
	c, err := rpc.Dial(network, address)
	if err != nil {
		return nil, nil
	}

	return &HelloWorldServiceClient{c}, nil
}

func (p *HelloWorldServiceClient) HelloWorld(request string, reply *string) error {
	return p.Client.Call(HelloWorldServiceName+".HelloWorld", request, reply)
}
```

### 客户端
```go
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
```
