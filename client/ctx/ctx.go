package ctx

import "github.com/awesome-cmd/chat/core/net"

type ChatContext struct {
	ID int64 `json:"id"`
	Name string `json:"name"`
	Conn *net.Conn
}