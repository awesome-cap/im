package net

import (
	"github.com/awesome-cmd/chat/core/protocol"
	"github.com/awesome-cmd/chat/core/util/async"
	"net"
	"sync/atomic"
)

var id int64 = 0

type Conn struct {
	ID int64 `json:"id"`

	conn net.Conn
	streams []byte
}

func NewConn(conn net.Conn) *Conn{
	return &Conn{
		ID: atomic.AddInt64(&id, 1),
		conn: conn,
		streams: make([]byte, 0),
	}
}

func (c *Conn) read() error{
	buffered := make([]byte, 1024)
	n, err := c.conn.Read(buffered)
	if err != nil{
		return err
	}
	c.streams = append(c.streams, buffered[0:n]...)
	return nil
}

func (c *Conn) parse() []byte {
	data, index, err := protocol.Decode(c.streams)
	if err != nil{
		return nil
	}
	c.streams = c.streams[index:]
	return data
}

func (c *Conn) Accept(apply func(data []byte, c *Conn)){
	for {
		err := c.read()
		if err != nil {
			_ = c.conn.Close()
			break
		}
		msg := c.parse()
		if msg != nil {
			async.Async(func() {
				apply(msg, c)
			})
		}
	}
}

func (c *Conn) Write(data []byte) error{
	_, err := c.conn.Write(protocol.Encode(data))
	return err
}

