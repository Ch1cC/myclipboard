package convert

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"myclipboard/clipboard"
	"sort"
	"sync"
	"time"
)

var KV sync.Map

// 创建一个字节缓冲区来保存压缩后的数据
var compressed bytes.Buffer

func Put(key int64, value interface{}, exp time.Duration) {
	KV.Store(key, value)
	fmt.Println("put key:", key)
	fmt.Println("put value:", value)
	time.AfterFunc(exp, func() {
		KV.Delete(key)
	})
}
func BuildJson() []byte {
	var jsonObj []clipboard.Clipboard
	var keys []int64
	//无序的
	KV.Range(func(k, v interface{}) bool {
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
		v, _ := KV.Load(int64(k))
		jsonObj = append(jsonObj, v.(clipboard.Clipboard))
	}
	jsonByte, _ := json.Marshal(jsonObj)
	compressed.Reset()
	// 创建一个gzip写入器，将数据写入到压缩缓冲区
	gzipWriter, _ := gzip.NewWriterLevel(&compressed, 2)
	_, err := gzipWriter.Write(jsonByte)
	if err != nil {
		fmt.Println("压缩数据时发生错误：", err)
	}

	// 关闭gzip写入器，这样会将剩余的数据刷新到缓冲区
	err = gzipWriter.Close()
	if err != nil {
		fmt.Println("关闭gzip写入器时发生错误：", err)
	}
	return compressed.Bytes()
}
