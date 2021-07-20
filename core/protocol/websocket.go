package protocol

import (
	"github.com/awesome-cap/im/core/util/json"
	"github.com/gorilla/websocket"
)

type WebsocketReadWriteCloser struct {
	conn *websocket.Conn
}

func NewWebsocketReadWriteCloser(conn *websocket.Conn) WebsocketReadWriteCloser {
	return WebsocketReadWriteCloser{conn: conn}
}

func (w WebsocketReadWriteCloser) Read() (*Msg, error) {
	_, b, err := w.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	msg := &Msg{}
	json.Unmarshal(b, msg)
	return msg, nil
}

func (w WebsocketReadWriteCloser) Write(msg Msg) error {
	return w.conn.WriteMessage(websocket.TextMessage, json.Marshal(msg))
}

func (w WebsocketReadWriteCloser) Close() error {
	return w.conn.Close()
}
