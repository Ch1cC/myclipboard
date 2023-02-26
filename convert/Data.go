package convert

import (
	"myclipboard/clipboard"
	"sort"
	"sync"
	"time"
)

var KV sync.Map

func Put(key interface{}, value interface{}, exp time.Duration) {
	KV.Store(key, value)
	time.AfterFunc(exp, func() {
		KV.Delete(key)
	})
}
func BuildJson() interface{} {
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
	return jsonObj
}
