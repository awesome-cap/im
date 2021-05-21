package main

import (
	"bufio"
	"fmt"
	"github.com/awesome-cmd/chat/core/model"
	xnet "github.com/awesome-cmd/chat/core/net"
	"net"
	"os"
)

var in = bufio.NewReader(os.Stdin)
func readline() string{
	line, _ := in.ReadString('\n')
	return line[0:len(line) - 1]
}

func main() {
	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:3333")
	conn, _ := net.DialTCP("tcp", nil, tcpAddr)
	
	xconn := xnet.NewConn(conn)

	go func() {
		for {
			str := readline()
			xconn.Write([]byte(str))
		}
	}()
	xconn.Accept(func(msg *model.Msg, c *xnet.Conn) {
		fmt.Println(string(msg.Data))
	})
}