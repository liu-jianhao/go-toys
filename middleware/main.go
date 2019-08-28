package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

var requestServed uint64

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello")
	log.Println("Greeted")
}

type statsHandler struct {}

func (s *statsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Request Served: %d\n", atomic.LoadUint64(&requestServed))
	log.Println("Stats provided")
}

/*
func counter(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		atomic.AddUint64(&requestServed, 1)
		log.Println("Counter >> Counted")
	}
}
*/
func counter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		atomic.AddUint64(&requestServed, 1)
		log.Println("Counter >> Counted")
	})
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r* http.Request) {
		log.Printf("Logger >> start %s %q\n", r.Method, r.URL.String())
		t := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("Logger >> end %s %q (%v)\n", r.Method, r.URL.String(), time.Now().Sub(t))
	})
}

func use(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for _, mw := range middleware {
		h = mw(h)
	}
	return h
}

func main() {
	//http.Handle("/hello", logger(counter(helloHandler)))
	http.Handle("/hello", use(http.HandlerFunc(helloHandler), counter, logger))

	s := &statsHandler{}
	http.Handle("/stat", use(s, logger))

	log.Println("Starting server...")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
