package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)


func index(w http.ResponseWriter, r *http.Request) {

	log.Println("get", r.RequestURI)

	//接收客户端 request，并将 request 中带的 header 写入 response header
	for k, v := range r.Header { // return string, []string
		fmt.Fprintf(w, "%s=%s\n", k, v)
		for _, vv := range v {
			w.Header().Set(k, vv)
		}
	}

	//读取当前系统的环境变量中的 VERSION 配置，并写入 response header
	version := os.Getenv("VERSION")
	fmt.Fprintf(w, "Enviroment VERSION is: %s\n", version)
	w.Header().Set("VERSION", version)

	//Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
	realIP := r.Header.Get("X-Real-IP")
	remoteIP := strings.Split(r.RemoteAddr, ":")[0]
	log.Println("The Client Header X-Real-IP: ", realIP)
	log.Println("The Client RemoteAddr: ", remoteIP)
	w.WriteHeader(200)

}

// 当访问 localhost/healthz 时，应返回 200
func healthz(w http.ResponseWriter, r *http.Request) {
	log.Println("get", r.RequestURI)
	fmt.Fprintf(w, "Status OK!")
	w.WriteHeader(200)
}

func main() {

	mux := http.NewServeMux() //声明多路复用mux对象
	// use mux.HandleFunc define routing
	mux.HandleFunc("/", index)
	mux.HandleFunc("/healthz", healthz)

	http.ListenAndServe("0.0.0.0:8080", mux)

}
