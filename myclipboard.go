package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"myclipboard/clipboard"
	"myclipboard/config"
	"myclipboard/convert"
	"myclipboard/ws"
	"net/http"
	"time"
)

func main() {
	flag.DurationVar(&config.Duration, "duration", time.Minute*15, "过期时间间隔,默认15分钟")
	hub := ws.NewHub()
	go hub.Run()
	http.HandleFunc("/ws", func(rw http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, rw, r)
	}) // 设置访问的路由
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		t1, err := template.ParseFiles("static/index.html")
		if err != nil {
			panic(err)
		}
		config.ConfigRandom()
		t1.Execute(rw, config.Token.String())
	})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Println("启动在127.0.0.1:9090")
	flag.Parse()
	fmt.Printf("过期时间间隔设置为%s\n", config.Duration)
	convert.KV.Store(time.Now().UnixMicro(), clipboard.Clipboard{UnixMicro: time.Now().UnixMicro(), Msg: "欢迎使用"})
	convert.KV.Store(time.Now().UnixMicro()+1, clipboard.Clipboard{UnixMicro: time.Now().UnixMicro(), Msg: "过期时间间隔为" + config.Duration.String()})
	err := http.ListenAndServe("127.0.0.1:9090", nil) // 设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
