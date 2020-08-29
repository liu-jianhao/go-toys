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
