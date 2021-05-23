package server

import (
	"github.com/awesome-cmd/chat/core/model"
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

func Run(args []string) {
	port := "3333"
	if len(args) > 0 {
		port = args[0]
	}
	listener, err := net.Listen("tcp", ":" + port)
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
			defer chats.Clean(c)
			err := c.Write(protocol.Msg{
				ID: 0,
				Data: json.Marshal(model.Resp{
					Code: 0,
					Type: "id",
					Data: []byte(strconv.FormatInt(c.ID, 10)),
				}),
			})
			if err != nil {
				log.Printf("c.Write err %v\n", err)
				return
			}
			err = c.Accept(func(msg protocol.Msg, c *xnet.Conn) {
				resp := events.Process(msg, c)
				if resp != nil {
					err := c.Write(protocol.Msg{
						ID: msg.ID,
						Data: json.Marshal(resp),
					})
					if err != nil{
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