package network

import (
	"github.com/awesome-cap/im/core/model"
	"github.com/awesome-cap/im/core/network"
	"github.com/awesome-cap/im/core/protocol"
	"github.com/awesome-cap/im/core/util/json"
	"github.com/awesome-cap/im/server/chats"
	"github.com/awesome-cap/im/server/cluster"
	"github.com/awesome-cap/im/server/events"
	"log"
	"strconv"
)

// Network is interface of all kinds of network.
type Network interface {
	Serve() error
}

func handle(rwc protocol.ReadWriteCloser) {
	c := network.NewConn(rwc)
	chats.BindClient(c)
	defer chats.Clean(c)

	// distribute id
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
	err = c.Accept(func(msg protocol.Msg, c *network.Conn) {
		event := model.Event{}
		json.Unmarshal(msg.Data, &event)
		event.From = chats.Client(c)
		resp := events.Process(msg.ID, event, cluster.NextID)
		if cluster.BroadcastEvents[event.Type] {
			cluster.Broadcast(json.Marshal(event))
		}
		if resp != nil {
			err := c.Write(protocol.Msg{
				ID:   msg.ID,
				Data: json.Marshal(resp),
			})
			if err != nil {
				log.Printf("c.Write err %v\n", err)
			}
		}
	})
	if err != nil {
		log.Printf("c.Accept err %v\n", err)
	}
}
