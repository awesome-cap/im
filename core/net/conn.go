package net

import (
	"github.com/awesome-cmd/chat/core/protocol"
	"github.com/awesome-cmd/chat/core/util/async"
	"net"
	"sync/atomic"
)

var connID int64 = 0

type Conn struct {
	ID int64 `json:"id"`

	msgID int64
	conn net.Conn
	streams []byte
}

func NewConn(conn net.Conn) *Conn{
	return &Conn{
		ID: atomic.AddInt64(&connID, 1),
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

func (c *Conn) parse() (*protocol.Msg, error) {
	msg, index, err := protocol.Decode(c.streams)
	if err != nil{
		return nil, err
	}
	c.streams = c.streams[index:]
	return msg, nil
}

func (c *Conn) Close() error {
	return c.conn.Close()
}

func (c *Conn) Accept(apply func(msg protocol.Msg, c *Conn)) error{
	for {
		err := c.read()
		if err != nil {
			return err
		}
		msg, err := c.parse()
		if err != nil {
			return err
		}
		if msg != nil {
			async.Async(func() {
				apply(*msg, c)
			})
		}
	}
}

func (c *Conn) nextMsgID() int64{
	return atomic.AddInt64(&c.msgID, 1)
}

func (c *Conn) Write(msg protocol.Msg) error{
	_, err := c.conn.Write(protocol.Encode(msg))
	return err
}

func (c *Conn) WriteID(data []byte) (int64, error){
	msgID := c.nextMsgID()
	return msgID, c.Write(protocol.Msg{
		ID: msgID,
		Data: data,
	})
}

