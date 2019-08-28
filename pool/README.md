近期在学习Go语言，自己实现了一个很简单的协程池，记录一下实现过程

## 总体架构
![在这里插入图片描述](https://img-blog.csdnimg.cn/20190826174255306.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3dlc3Ricm9va2xpdQ==,size_16,color_FFFFFF,t_70)
可以看到，架构很简单，就是客户通过`Channel`传一些任务给协程池，然后协程池里的`worker`就会执行
## Task结构体
```go
type Task struct {
	taskId int
	f func() error
}
```
Task结构体很简单，一个`id`，一个要执行的函数

```go
// Task的构造函数
func NewTask(id int, f func() error) *Task {
	return &Task{
		taskId: id,
		f: f,
	}
}
// 执行task
func (t *Task) execute() {
	t.f()
}
```


## 协程池结构体
```go
type Pool struct {
	workerNum int
	EntryChan chan *Task
	workerChan chan *Task
}
```
协程池结构体也很简单
+ `workerNum`：worker的数量，由客户指定
+ `EntryChan`：提供给客户传任务的一个通道
+ `workerChan`：协程池内部传给worker的一个通道

初始化协程池
```go
func NewPool(num int) *Pool {
	return &Pool{
		workerNum: num,
		EntryChan: make(chan *Task),
		workerChan: make(chan *Task),
	}
}
```

worker执行
```go
func (p *Pool) worker(id int) {
	for task := range p.workerChan {
		task.execute()
		fmt.Println("workerId:", id, "taskId:", task.taskId, "is done")
	}
}
```

协程池启动
```go
func (p *Pool) Run() {
	for i := 0; i < p.workerNum; i++ {
		go p.worker(i)
	}

	for task := range p.EntryChan {
		p.workerChan <- task
	}
}
```

## 测试
```go
func task() error {
	fmt.Println(time.Now(), "Do something")
	return nil
}

func main() {
	p := pool.NewPool(3)

	id := 0
	go func() {
		for {
			p.EntryChan <- pool.NewTask(id, task)
			id++
		}
	}()

	p.Run()
}
```

## 完整代码
pool目录下的`pool.go`
```go
package pool

import "fmt"

type Task struct {
	taskId int
	f func() error
}

func NewTask(id int, f func() error) *Task {
	return &Task{
		taskId: id,
		f: f,
	}
}

func (t *Task) execute() {
	t.f()
}

type Pool struct {
	workerNum int
	EntryChan chan *Task
	workerChan chan *Task
}

func NewPool(num int) *Pool {
	return &Pool{
		workerNum: num,
		EntryChan: make(chan *Task),
		workerChan: make(chan *Task),
	}
}

func (p *Pool) worker(id int) {
	for task := range p.workerChan {
		task.execute()
		fmt.Println("workerId:", id, "taskId:", task.taskId, "is done")
	}
}

func (p *Pool) Run() {
	for i := 0; i < p.workerNum; i++ {
		go p.worker(i)
	}

	for task := range p.EntryChan {
		p.workerChan <- task
	}
}
```

`main.go`
```go
package main

import (
	"./pool"
	"fmt"
	"time"
)

func task() error {
	fmt.Println(time.Now(), "Do something")
	return nil
}

func main() {
	p := pool.NewPool(3)

	id := 0
	go func() {
		for {
			p.EntryChan <- pool.NewTask(id, task)
			id++
		}
	}()

	p.Run()
}
```
