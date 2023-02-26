package main

import (
	"flag"
	"fmt"
	"log"
	"myclipboard/clipboard"
	"myclipboard/config"
	"myclipboard/convert"
	"myclipboard/ws"
	"net/http"
	"time"
)

// Put 缓存过期功能实现 类Redis

/*func postData(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	unixMicro := time.Now().UnixMicro()
	convert.Put(unixMicro, clipboard.Clipboard{UnixMicro: unixMicro, Msg: r.Form.Get("data")}, config.Duration)
	w.Header().Set("content-type", "text/json")
	msg, _ := json.Marshal(convert.BuildJson())
	fmt.Fprintf(w, string(msg)) // 这个写入到 w 的是输出到客户端的
}
func getData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/json")
	msg, _ := json.Marshal(convert.BuildJson())
	fmt.Fprintf(w, string(msg)) // 这个写入到 w 的是输出到客户端的
}*/

func main() {
	flag.DurationVar(&config.Duration, "duration", time.Minute*15, "过期时间间隔,默认15分钟")
	/*http.HandleFunc("/post", postData)       // 设置访问的路由
	http.HandleFunc("/get", getData)         // 设置访问的路由*/
	http.HandleFunc("/ws", ws.SocketHandler) // 设置访问的路由
	http.Handle("/", http.FileServer(http.Dir("./static")))
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
