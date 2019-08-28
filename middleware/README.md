## 中间件

中间件作用就是将业务和非业务代码功能解耦

### 用法
```go
func middlewareUsingHandlerFunc(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// the middlerware's logic here...
		f(w, r) // equivalent to f.ServeHTTP(w, r)
	}
}

func middlewareUsingHander(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// the middlerware's logic here...
		next.ServeHTTP(w, r)
	})
}
```
可以看到中间件有两种用法，一种返回`http.Handler`，一种返回`http.HandlerFunc`，
这是因为http有两种用法：
```go
http.Handle("/", http.HandlerFunc(f))
http.HandleFunc("/", f)
```

不难看出，中间件就是在业务代码外层套一层非业务逻辑的代码


### 例子
运行`main.go`：
```go
go run main.go
```
在另外一个终端运行：
```go
> curl 127.0.0.1:8080/stat -i
HTTP/1.1 200 OK
Date: Wed, 28 Aug 2019 04:49:44 GMT
Content-Length: 18
Content-Type: text/plain; charset=utf-8

Request Served: 0

```
然后在运行`main.go`的终端可以看到：
```
2019/08/28 12:49:41 Starting server...
2019/08/28 12:49:44 Logger >> start GET "/stat"
2019/08/28 12:49:44 Stats provided
2019/08/28 12:49:44 Logger >> end GET "/stat" (16.226µs)
```