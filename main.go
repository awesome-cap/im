/**
2 * @Author: Nico
3 * @Date: 2021/5/23 17:08
4 */
package main

import (
	"github.com/awesome-cmd/chat/client"
	"github.com/awesome-cmd/chat/server"
	"log"
	"os"
)


func main()  {
	if len(os.Args) < 2 {
		log.Fatal("please input chat -s or -c")
	}
	if os.Args[1] == "-s" {
		server.Run(os.Args[2:])
	}else if os.Args[1] == "-c" {
		client.Run(os.Args[2:])
	}else{
		log.Fatal("-s or -c required")
	}
}


