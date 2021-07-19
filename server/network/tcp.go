package network

import (
	"github.com/awesome-cap/im/core/protocol"
	"github.com/awesome-cap/im/core/util/async"
	"log"
	"net"
)

type Tcp struct {
	addr string
}

func NewTcpServer(addr string) Tcp {
	return Tcp{addr: addr}
}

func (t Tcp) Serve() error {
	listener, err := net.Listen("tcp", t.addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Tcp server listener on %s\n", t.addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("listener.Accept err %v\n", err)
			continue
		}
		async.Async(func() {
			handle(protocol.NewTcpReadWriteCloser(conn))
		})
	}
}
