package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var kV sync.Map

// Set 缓存过期功能实现 类Redis
func Set(key interface{}, value interface{}, exp time.Duration) {
	kV.Store(key, value)
	time.AfterFunc(exp, func() {
		kV.Delete(key)
	})
}

func postData(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	Set(time.Now().UnixNano(), r.Form.Get("data"), time.Minute*5)
	w.Header().Set("content-type", "text/json")
	jsonobj := make(map[int64]interface{})
	kV.Range(func(k interface{}, v interface{}) bool {
		jsonobj[k.(int64)] = v
		return true
	})
	msg, _ := json.Marshal(jsonobj)
	fmt.Fprintf(w, string(msg)) // 这个写入到 w 的是输出到客户端的
}
func getData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/json")
	jsonobj := make(map[int64]interface{})
	kV.Range(func(k interface{}, v interface{}) bool {
		jsonobj[k.(int64)] = v
		return true
	})
	msg, _ := json.Marshal(jsonobj)
	fmt.Fprintf(w, string(msg)) // 这个写入到 w 的是输出到客户端的
}

func main() {
	http.HandleFunc("/post", postData) // 设置访问的路由
	http.HandleFunc("/get", getData)   // 设置访问的路由
	http.Handle("/", http.FileServer(http.Dir("./static")))
	fmt.Println("启动在127.0.0.1:9090")
	err := http.ListenAndServe("127.0.0.1:9090", nil) // 设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
