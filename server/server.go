package server

import (
	"flag"
	"github.com/awesome-cmd/im/core/model"
	xnet "github.com/awesome-cmd/im/core/net"
	"github.com/awesome-cmd/im/core/protocol"
	"github.com/awesome-cmd/im/core/util/async"
	"github.com/awesome-cmd/im/core/util/json"
	"github.com/awesome-cmd/im/server/chats"
	"github.com/awesome-cmd/im/server/cluster"
	"github.com/awesome-cmd/im/server/events"
	"github.com/awesome-cmd/im/server/task"
	"log"
	"net"
	"strconv"
	"strings"
)

var (
	port         int
	clusterPort  int
	clusterSeeds string
)

func Run() {
	flag.IntVar(&port, "p", 3333, "server port.")
	flag.IntVar(&clusterPort, "cluster-port", 3334, "cluster port.")
	flag.StringVar(&clusterSeeds, "cluster-seeds", "", "cluster seeds.")
	flag.Parse()

	// cluster
	seeds := make([]string, 0)
	if clusterSeeds != "" {
		seeds = strings.Split(clusterSeeds, ",")
	}
	err := cluster.Start(clusterPort, seeds)
	if err != nil {
		log.Fatal(err)
	}

	// task
	task.Start()

	// server
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("listener on %d\n", port)
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
			err = c.Accept(func(msg protocol.Msg, c *xnet.Conn) {
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
		})
	}
}
