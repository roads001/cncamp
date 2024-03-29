package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	started          = time.Now()
	requestDurations = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "httpserver",
			Name:      "http_request_duration_seconds",
			Help:      "A histogram of the HTTP request durations in seconds.",
			//Buckets: []float64{0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
			Buckets: prometheus.ExponentialBuckets(0.001, 2, 15),
		}, []string{"step"})
)

type ExecutionTimer struct {
	histo *prometheus.HistogramVec
}

func (et ExecutionTimer) Observe(value float64) {
	(et.histo).WithLabelValues("total").Observe(value)
}

func index(w http.ResponseWriter, r *http.Request) {

	log.Println("get", r.RequestURI)

	// 启动一个计时器
	exeTime := ExecutionTimer{
		histo: requestDurations,
	}
	timer := prometheus.NewTimer(exeTime)
	// 停止计时器
	defer timer.ObserveDuration()

	// sleep 0-2s
	randInt := rand.Intn(2000)
	// 将持续时间放进 requestDurations 的直方图指标中去
	//exeTime.Observe(float64(randInt / 1000))
	time.Sleep(time.Millisecond * time.Duration(randInt))
	log.Println("sleep ", randInt)

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
	//w.WriteHeader(200)

}

// 当访问 localhost/healthz 时，应返回 200
func healthz(w http.ResponseWriter, r *http.Request) {
	log.Println("get", r.RequestURI)
	fmt.Fprintf(w, "Status OK!")
	//w.WriteHeader(200)
}

// 当访问 localhost/healthz10
func healthz10(w http.ResponseWriter, r *http.Request) {
	// code aus kubernetes.io
	duration := time.Now().Sub(started)
	// at the first 10 second, status ok
	// then, error
	if duration.Seconds() > 10 {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("error: %v", duration.Seconds())))
	} else {
		//w.WriteHeader(200)
		w.Write([]byte("Status OK!"))
	}
}

func main() {

	prometheus.MustRegister(requestDurations)

	mux := http.NewServeMux() //声明多路复用mux对象
	// use mux.HandleFunc define routing
	mux.HandleFunc("/", index)
	mux.HandleFunc("/healthz", healthz)
	mux.HandleFunc("/healthz10", healthz10)

	// register a new handler for the /metrics endpoint
	mux.Handle("/metrics", promhttp.Handler())

	http.ListenAndServe("0.0.0.0:8765", mux)

}
