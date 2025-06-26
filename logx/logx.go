// logx/logx.go
package logx

import (
	"fmt"
	"log"
	"os"
	"time"
)

var Logger *log.Logger

func init() {
	dir := "./log"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0766)
	}
	file := "./log/" + time.Now().Format("2006-01-02") + ".log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	Logger = log.New(logFile, "[myclipboard]", log.LstdFlags|log.Lshortfile)
}
