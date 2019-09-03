package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

type client struct {
	ch chan<- string
	name string
}

var (
	enterChan = make(chan client)
	leaveChan = make(chan client)
	msgChan = make(chan string)
)

func broadcaster() {
	// 客户集合
	clients := make(map[string]client)
	for {
		select {
		case msg := <-msgChan:
			// 广播给所有客户
			for _, cli := range clients {
				cli.ch <- msg
			}
		case cli := <- enterChan:
			clients[cli.name] = cli
			// 通知当前用户现在有哪些用户在线
			if len(clients) > 1 {
				cli.ch <- "All clients:"
				for _, c := range clients {
					cli.ch <- c.name
				}
			}
		case cli := <- leaveChan:
			delete(clients, cli.name)
			close(cli.ch)
		}
	}
}

// 超过指定时间没有发送过消息的客户会被剔除
func countDown(c net.Conn, cli client, notIdleCh <-chan interface{}) {
	ticker := time.NewTicker(time.Second)
	counter := 0
	max := 10
	for {
		select {
		case <- ticker.C:
			counter++
			if counter == max {
				msg := cli.name + " idle too long. Kicked out"
				msgChan <- msg
				fmt.Fprintln(c, msg)
				ticker.Stop()
				c.Close()
				return
			}
		case <- notIdleCh:
			counter = 0
		}
	}
}

func handleConn(c net.Conn) {
	// 输入名字
	input := bufio.NewScanner(c)
	var who string
	fmt.Fprint(c, "Input your name: ")
	if input.Scan() {
		who = input.Text()
	}

	ch := make(chan string)
	me := client {
		ch: ch,
		name: who,
	}
	go output(c, ch)

	ch <- "You're " + me.name
	msgChan <- me.name + " has arrived"
	enterChan <- me

	notIdleCh := make(chan interface{})
	go countDown(c, me, notIdleCh)

	for input.Scan() {
		notIdleCh <- struct{}{}
		msgChan <- me.name + ": " + input.Text() + "\t" + time.Now().Format("15:04:05")
	}

	leaveChan <- me
	msgChan <- me.name + " has left"
	c.Close()
}

func output(c net.Conn, ch chan string) {
	// 直到通道关闭
	for msg := range ch {
		fmt.Fprintln(c, msg)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal(err)
	}

	// broadcast在后台处理客户发送的消息
	go broadcaster()
	// 主goroutine循环接收客户连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}

		// 新开一个goroutine处理客户消息
		go handleConn(conn)
	}
}
