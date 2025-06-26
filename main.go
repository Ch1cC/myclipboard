package main

import (
	"flag"
	"fmt"
	"log"
	"myclipboard/config"
	"myclipboard/convert"
	"myclipboard/logx"
	"myclipboard/ws"
	"net/http"
	"time"
)

var version string // 用于存储版本号
var port int
var crtPath string
var keyPath string

func main() {
	flag.DurationVar(&config.Duration, "duration", time.Minute*15, "过期时间间隔")
	flag.IntVar(&port, "port", 9090, "端口")
	flag.StringVar(&crtPath, "crtPath", "./server.crt", "证书文件")
	flag.StringVar(&keyPath, "keyPath", "./server.key", "证书密钥")
	hub := ws.NewHub()
	go hub.Run()
	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	})
	// 设置访问的路由
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "./html/dist/index.html")
	// })
	// http.Handle("/dist/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	w.Header().Add("Cache-Control", "no-cache")
	// 	if strings.HasSuffix(r.URL.Path, ".wasm") {
	// 		w.Header().Set("content-type", "application/wasm")
	// 	}
	// 	http.StripPrefix("/dist/", http.FileServer(http.Dir("./html/dist"))).ServeHTTP(w, r)
	// }))
	flag.Parse()
	logx.Logger.Printf(":%d\n", port)
	logx.Logger.Printf("过期时间间隔设置为%s\n", config.Duration)
	logx.Logger.Printf("当前版本号%s\n", version)
	convert.KV.Store(time.Now().Unix(), convert.Row{Unix: time.Now().Unix(), Msg: []byte("当前版本号:" + version)})
	convert.KV.Store(time.Now().Unix()+1, convert.Row{Unix: time.Now().Unix(), Msg: []byte("过期时间间隔为:" + config.Duration.String())})
	err := http.ListenAndServeTLS(fmt.Sprintf(":%d", port), crtPath, keyPath, nil) // 设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
