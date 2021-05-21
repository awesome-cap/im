package main

import (
	"flag"
	"github.com/awesome-cmd/chat/core/model"
	xnet "github.com/awesome-cmd/chat/core/net"
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
			c.Accept(func(data []byte, c *xnet.Conn) {
				event := model.Event{}
				json.Unmarshal(data, event)
				err := c.Write(json.Marshal(events.Process(event, c)))
				if err != nil{
					log.Printf("c.Write err %v\n", err)
				}
			})
		})
	}
}