/**
2 * @Author: Nico
3 * @Date: 2021/5/23 17:08
4 */
package main

import (
	"github.com/awesome-cmd/dchat/client"
	"github.com/awesome-cmd/dchat/server"
	"log"
	"os"
)

func main()  {
	if len(os.Args) < 2 {
		log.Fatal("please input chat -s or -c")
	}
	if os.Args[1] == "-s" {
		server.Run()
	}else if os.Args[1] == "-c" {
		client.Run()
	}else{
		log.Fatal("-s or -c required")
	}
}


