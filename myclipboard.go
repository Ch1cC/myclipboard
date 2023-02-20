package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"
)

var (
	kV       sync.Map
	duration time.Duration
)

type clipboard struct {
	UnixMicro int64  `json:"unixMicro"`
	Msg       string `json:"msg"`
}

// Set 缓存过期功能实现 类Redis
func Set(key interface{}, value interface{}, exp time.Duration) {
	kV.Store(key, value)
	time.AfterFunc(exp, func() {
		kV.Delete(key)
	})
}

func postData(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	unixMicro := time.Now().UnixMicro()
	Set(unixMicro, clipboard{UnixMicro: unixMicro, Msg: r.Form.Get("data")}, duration)
	w.Header().Set("content-type", "text/json")
	msg, _ := json.Marshal(buildJson())
	fmt.Fprintf(w, string(msg)) // 这个写入到 w 的是输出到客户端的
}
func getData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/json")
	msg, _ := json.Marshal(buildJson())
	fmt.Fprintf(w, string(msg)) // 这个写入到 w 的是输出到客户端的
}

func buildJson() interface{} {
	var jsonObj []clipboard
	var keys []int64
	//无序的
	kV.Range(func(k, v interface{}) bool {
		//取出所有的key
		keys = append(keys, k.(int64))
		return true
	})
	//转化int64 to int
	sortKeys := make([]int, len(keys))
	for i := range sortKeys {
		sortKeys[i] = int(keys[i])
	}
	//倒序
	sort.Sort(sort.Reverse(sort.IntSlice(sortKeys)))
	for _, k := range sortKeys {
		v, _ := kV.Load(int64(k))
		jsonObj = append(jsonObj, v.(clipboard))
	}
	return jsonObj
}
func main() {
	flag.DurationVar(&duration, "duration", time.Minute*15, "过期时间间隔,默认15分钟")
	http.HandleFunc("/post", postData) // 设置访问的路由
	http.HandleFunc("/get", getData)   // 设置访问的路由
	http.Handle("/", http.FileServer(http.Dir("./static")))
	fmt.Println("启动在127.0.0.1:9090")
	flag.Parse()
	fmt.Printf("过期时间间隔设置为%s\n", duration)
	err := http.ListenAndServe("127.0.0.1:9090", nil) // 设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
