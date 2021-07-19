package network

import "C"
import (
	"github.com/awesome-cap/im/core/protocol"
	"sync/atomic"
)

var connID int64 = 0
var msgID int64 = 0

type Conn struct {
	ID int64 `json:"id"`

	state   int
	conn    protocol.ReadWriteCloser
	streams []byte
}

func NewConn(conn protocol.ReadWriteCloser) *Conn {
	return &Conn{
		ID:      atomic.AddInt64(&connID, 1),
		conn:    conn,
		streams: make([]byte, 0),
	}
}

func (c *Conn) Close() error {
	c.state = 1
	return c.conn.Close()
}

func (c *Conn) State() int {
	return c.state
}

func (c *Conn) Accept(apply func(msg protocol.Msg, c *Conn)) error {
	for {
		msg, err := c.conn.Read()
		if err != nil {
			return err
		}
		apply(*msg, c)
	}
}

func nextMsgID() int64 {
	return atomic.AddInt64(&msgID, 1)
}

func (c *Conn) Write(msg protocol.Msg) error {
	return c.conn.Write(msg)
}

func (c *Conn) WriteID(data []byte) (int64, error) {
	msgID := nextMsgID()
	return msgID, c.Write(protocol.Msg{
		ID:   msgID,
		Data: data,
	})
}
