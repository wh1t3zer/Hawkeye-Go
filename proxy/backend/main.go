package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// RealServer ...
type RealServer struct {
	Addr string
}

// HHandler ...
func (r *RealServer) HHandler(w http.ResponseWriter, req *http.Request) {
	upath := fmt.Sprintf("http://%s%s\n", r.Addr, req.URL.Path)
	realIP := fmt.Sprintf("RemoteAddr=%s,X-Forwarded-For=%v,X-Real-Ip=%v\n", req.RemoteAddr, req.Header.Get("X-Forwarded-For"), req.Header.Get("X-Real-Ip"))
	header := fmt.Sprintf("headers =%v\n", req.Header)
	io.WriteString(w, upath)
	io.WriteString(w, realIP)
	io.WriteString(w, header)
}

// EHandler ...
func (r *RealServer) EHandler(w http.ResponseWriter, req *http.Request) {
	upath := "error handler\n"
	w.WriteHeader(500)
	io.WriteString(w, upath)
}

// TimeoutHandler ...
func (r *RealServer) TimeoutHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("ddd")
	time.Sleep(2 * time.Second)
	time.Sleep(4 * time.Second)
	upath := "timeout handler\n" // 超时时间小于http.Server.WriteTimeout才会返回内容
	w.WriteHeader(200)           // 不然会返回 "Empty reply from server"
	io.WriteString(w, upath)
}

// Run ...
func (r *RealServer) Run() {
	log.Println("Starting httpserver at " + r.Addr)
	mux := http.NewServeMux()
	mux.HandleFunc("/", r.HHandler)
	mux.HandleFunc("/base/error", r.EHandler)
	mux.HandleFunc("/test_http_string/test_http_string/aaa", r.TimeoutHandler)
	server := &http.Server{
		Addr:         r.Addr,
		WriteTimeout: time.Second * 3,
		Handler:      mux,
	}
	go func() {
		log.Fatal(server.ListenAndServe())
	}()
}

func main() {
	rs1 := &RealServer{Addr: "127.0.0.1:2003"}
	rs2 := &RealServer{Addr: "127.0.0.1:2004"}
	rs1.Run()
	rs2.Run()

	//监听关闭信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
