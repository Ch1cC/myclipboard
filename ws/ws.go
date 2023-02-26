package ws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"myclipboard/clipboard"
	"myclipboard/config"
	"myclipboard/convert"
	"net/http"
	"time"
)

var (
	upgrader         = websocket.Upgrader{}
	pongWait         = 60 * time.Second  //等待时间
	pingPeriod       = 9 * pongWait / 10 //周期54s
	maxMsgSize int64 = 512               //消息最大长度
	writeWait        = 10 * time.Second  //
	send       chan []byte
)

func SocketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}
	data, _ := json.Marshal(convert.BuildJson())
	_ = conn.WriteMessage(websocket.TextMessage, data)
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error during message reading:", err)
		}
		//log.Printf("Received: %s", message)
		if len(message) > 0 {
			put(message)
		}
		data, _ := json.Marshal(convert.BuildJson())
		err = conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Println("Error during message writing:", err)
		}
	}
}
func put(message []byte) {
	unixMicro := time.Now().UnixMicro()
	convert.Put(unixMicro, clipboard.Clipboard{UnixMicro: unixMicro, Msg: string(message)}, config.Duration)
}
