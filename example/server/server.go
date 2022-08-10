package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lainio/err2/try"
	"github.com/shynome/wrtchttp"
	"github.com/shynome/wrtchttp/adapter"
	"github.com/shynome/wrtchttp/example/config"
)

var upgrader = websocket.Upgrader{}

func main() {
	adapter := try.To1(
		adapter.NewAdapter(config.SignalerServer))
	go adapter.Listen()
	l := &wrtchttp.Listener{Adapter: adapter}
	http.HandleFunc("/h", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello world")
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("z1")
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Fprintln(w, "ws upgrade failed")
			return
		}
		defer conn.Close()
		fmt.Println("z2")
		var a = map[string]string{"a": "b"}
		for i := 0; i < 5; i++ {
			time.Sleep(time.Second)
			fmt.Println("w1")
			if err = conn.WriteJSON(a); err != nil {
				fmt.Fprintln(w, "ws write failed. err: ", err)
				return
			}
			fmt.Println("w2")
		}
	})
	http.Serve(l, nil)
}
