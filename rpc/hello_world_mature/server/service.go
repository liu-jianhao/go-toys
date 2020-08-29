package main

type HelloWorldService struct{}

func (s *HelloWorldService) HelloWorld(request string, reply *string) error {
	*reply = "hello:" + request
	return nil
}
