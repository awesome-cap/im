package client

import (
	"flag"
	"github.com/awesome-cmd/chat/client/shell"
	"log"
)

var (
	name string
)

func Run() {
	flag.Bool("c", true, "")
	flag.StringVar(&name, "n", "doge", "your name.")
	flag.Parse()

	if len(name) <= 0 {
		log.Fatal("name can't be empty")
	}
	if len([]rune(name)) > 30{
		log.Fatal("name length must in 30 char")
	}
	shell.New(name).Start()
}