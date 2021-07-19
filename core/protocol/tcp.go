package protocol

import (
	"net"
)

type TcpReadWriteCloser struct {
	conn net.Conn
}

func NewTcpReadWriteCloser(conn net.Conn) TcpReadWriteCloser {
	return TcpReadWriteCloser{conn: conn}
}

func (t TcpReadWriteCloser) Read() (*Msg, error) {
	return Decode(t.conn)
}

func (t TcpReadWriteCloser) Write(msg Msg) error {
	_, err := t.conn.Write(Encode(msg))
	return err
}

func (t TcpReadWriteCloser) Close() error {
	return t.conn.Close()
}
