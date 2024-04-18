package convert

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"myclipboard/config"
	"sort"
	"sync"
	"time"
)

var KV sync.Map

// 创建一个字节缓冲区来保存压缩后的数据
var compressed bytes.Buffer

func put(key int64, value interface{}, exp time.Duration) {
	KV.Store(key, value)
	time.AfterFunc(exp, func() {
		KV.Delete(key)
	})
}

type Row struct {
	UnixMicro int64  `json:"unixMicro"`
	Msg       []byte `json:"msg"`
}

/*
*
存储数据
*/
func Put(message []byte) {
	unixMicro := time.Now().UnixMicro()
	put(unixMicro, Row{UnixMicro: unixMicro, Msg: message}, config.Duration)
}

func BuildJson() []byte {
	var rowArray []Row
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
		rowArray = append(rowArray, v.(Row))
	}
	jsonByte, _ := json.Marshal(rowArray)
	compressed.Reset()
	// 创建一个gzip写入器，将数据写入到压缩缓冲区
	if gzipWriter, err := gzip.NewWriterLevel(&compressed, gzip.NoCompression); err != nil {
		fmt.Println("压缩数据时发生错误：", err)
	} else {
		if _, err := gzipWriter.Write(jsonByte); err != nil {
			fmt.Println("压缩数据时发生错误：", err)
		}
		// 关闭gzip写入器，这样会将剩余的数据刷新到缓冲区
		if err := gzipWriter.Close(); err != nil {
			fmt.Println("关闭gzip写入器时发生错误：", err)
		}
	}
	return compressed.Bytes()
}
