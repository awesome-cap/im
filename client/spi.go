package client

import (
	"github.com/awesome-cmd/chat/client/shell"
	"log"
)

func Run(args []string) {
	if len(args) == 0 {
		log.Fatal("Please input your name !")
	}
	name := args[0]
	if len(name) < 0 {
		log.Fatal("name can't be empty")
	}
	if len([]rune(name)) > 30{
		log.Fatal("name length must in 30 char")
	}
	shell.New(name).Start()
}