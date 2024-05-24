package main

import (
	"flag"
	"fmt"
	"log"
	"myclipboard/config"
	"myclipboard/convert"
	"myclipboard/ws"
	"net/http"
	"strings"
	"time"
)

var version string // 用于存储版本号
func main() {
	var port int
	flag.DurationVar(&config.Duration, "duration", time.Minute*15, "过期时间间隔")
	flag.IntVar(&port, "port", 9090, "端口")
	hub := ws.NewHub()
	go hub.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	})
	// 设置访问的路由
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})
	http.Handle("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		if strings.HasSuffix(r.URL.Path, ".wasm") {
			w.Header().Set("content-type", "application/wasm")
		}
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))).ServeHTTP(w, r)
	}))
	flag.Parse()
	fmt.Printf("127.0.0.1:%d\n", port)
	fmt.Printf("过期时间间隔设置为%s\n", config.Duration)
	convert.KV.Store(time.Now().Unix(), convert.Row{Unix: time.Now().Unix(), Msg: []byte("当前版本号:" + version)})
	convert.KV.Store(time.Now().Unix()+1, convert.Row{Unix: time.Now().Unix(), Msg: []byte("过期时间间隔为" + config.Duration.String())})
	err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), nil) // 设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
