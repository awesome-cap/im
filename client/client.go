package client

import (
	"flag"
	"github.com/awesome-cap/im/client/shell"
	"log"
)

var (
	name  string
	addrs string
	net   string
)

func Run() {
	flag.StringVar(&name, "n", "doge", "your name.")
	flag.StringVar(&addrs, "addrs", "", "server addrs.")
	flag.StringVar(&net, "net", "tcp", "network.")
	flag.Parse()

	if len(name) <= 0 {
		log.Fatal("name can't be empty")
	}
	if len([]rune(name)) > 30 {
		log.Fatal("name length must in 30 char")
	}
	shell.New(name, net, addrs).Start()
}
