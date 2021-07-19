package server

import (
	"flag"
	"github.com/awesome-cap/im/core/util/async"
	"github.com/awesome-cap/im/server/cluster"
	"github.com/awesome-cap/im/server/network"
	"github.com/awesome-cap/im/server/task"
	"log"
	"strconv"
	"strings"
	"sync"
)

var (
	port         int
	wsPort       int
	clusterPort  int
	clusterSeeds string
)

func Run() {
	flag.IntVar(&port, "p", 3333, "server port.")
	flag.IntVar(&wsPort, "ws-port", 0, "server port.")
	flag.IntVar(&clusterPort, "cluster-port", 0, "cluster port.")
	flag.StringVar(&clusterSeeds, "cluster-seeds", "", "cluster seeds.")
	flag.Parse()

	// cluster
	if clusterPort > 0 {
		seeds := make([]string, 0)
		if clusterSeeds != "" {
			seeds = strings.Split(clusterSeeds, ",")
		}
		err := cluster.Start(clusterPort, seeds)
		if err != nil {
			log.Fatal(err)
		}
	}

	// task
	task.Start()

	// server
	networks := make([]network.Network, 0)
	networks = append(networks, network.NewTcpServer(":"+strconv.Itoa(port)))
	if wsPort > 0 {
		networks = append(networks, network.NewWebsocketServer(":"+strconv.Itoa(wsPort)))
	}
	wg := sync.WaitGroup{}
	wg.Add(len(networks))
	for _, n := range networks {
		network := n
		async.Async(func() {
			defer wg.Add(-1)
			err := network.Serve()
			if err != nil {
				log.Println(err)
			}
		})
	}
	wg.Wait()
}
