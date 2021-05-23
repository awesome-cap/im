package shell

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/awesome-cmd/chat/client/ctx"
	"github.com/awesome-cmd/chat/core/util/http"
	"github.com/awesome-cmd/chat/core/util/json"
	"os"
	"strings"
)

var servers = []string{
	"https://raw.githubusercontent.com/awesome-cmd/chat/main/servers.json",
	"https://gitee.com/ainilili/chat/raw/main/servers.json",
}

type shell struct {
	in *bufio.Reader
	out *bufio.Writer
	ctx *ctx.ChatContext
	position *directory
}

func New(name string) *shell{
	return &shell{
		in: bufio.NewReader(os.Stdin),
		out: bufio.NewWriter(os.Stdout),
		ctx: ctx.NewContext(name),
		position: root,
	}
}

func (s *shell) Start(){
	for{
		fmt.Printf("[%s@chat %s]# ", s.ctx.Name, s.position.name)
		inputs, err := s.readline()
		if err != nil{
			fmt.Println(err.Error())
			continue
		}
		res, err := s.position.action(s, inputs)
		if err != nil{
			fmt.Println(err.Error())
			continue
		}
		fmt.Print(res)
	}
}

func (s *shell) refreshServerList() error{
	serverList := make([]string, 0)
	for _, server := range servers {
		resp, err := http.Get(server)
		if err != nil && resp == ""{
			continue
		}
		json.Unmarshal([]byte(resp), &serverList)
	}
	if len(serverList) == 0 {
		return errors.New("no available server. ")
	}
	s.position.reset()
	for _, v := range serverList{
		serverInfo := strings.Split(v, "|")
		s.position.add(newDirectory(strings.ToLower(serverInfo[1]), strings.ToLower( serverInfo[0]), serverActions, baseActions))
	}
	return nil
}

func (s *shell) readline() ([]byte, error){
	inputs, err := s.in.ReadBytes('\n')
	return []byte(strings.TrimSpace(string(inputs))), err
}

func (s *shell) eraseLine() {
	fmt.Printf("\033[1A")
	fmt.Printf("\r\r")
}