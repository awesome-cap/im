package main

import (
	"flag"
	xnet "github.com/awesome-cmd/chat/core/net"
	"github.com/awesome-cmd/chat/core/protocol"
	"github.com/awesome-cmd/chat/core/util/async"
	"github.com/awesome-cmd/chat/core/util/json"
	"github.com/awesome-cmd/chat/server/chats"
	"github.com/awesome-cmd/chat/server/events"
	"log"
	"net"
	"strconv"
)

var port int

func init(){
	flag.IntVar(&port, "p", 3333, "Port of Server")
}

func main() {
	listener, err := net.Listen("tcp", ":" + strconv.Itoa(port))
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("listener.Accept err %v\n", err)
			continue
		}
		async.Async(func() {
			c := xnet.NewConn(conn)
			chats.BindClient(c)
			err := c.Accept(func(msg protocol.Msg, c *xnet.Conn) {
				resp := events.Process(msg, c)
				if resp != nil {
					err := c.Write(protocol.Msg{
						ID: msg.ID,
						Data: json.Marshal(resp),
					})
					if err != nil{
						chats.Clean(c)
						log.Printf("c.Write err %v\n", err)
					}
				}
			})
			if err != nil {
				log.Printf("c.Accept err %v\n", err)
			}
		})
	}
}