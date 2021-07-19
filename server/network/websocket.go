package network

import (
	"github.com/awesome-cap/im/core/protocol"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Websocket struct {
	addr string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewWebsocketServer(addr string) Websocket {
	return Websocket{addr: addr}
}

func (w Websocket) Serve() error {
	http.HandleFunc("/im", serveWs)
	log.Printf("Websocket server listener on %s\n", w.addr)
	return http.ListenAndServe(w.addr, nil)
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	handle(protocol.NewWebsocketReadWriteCloser(conn))
}
