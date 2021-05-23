package main

import (
	"github.com/awesome-cmd/chat/client/shell"
	"log"
	"os"
)

func main() {
	os.Args = append(os.Args, "nico")
	if len(os.Args) == 1 {
		log.Fatal("Please input your name !")
	}
	name := os.Args[1]
	if len(name) < 0 {
		log.Fatal("name can't be empty")
	}
	if len([]rune(name)) > 30{
		log.Fatal("name length must in 30 char")
	}
	shell.New(name).Start()

}