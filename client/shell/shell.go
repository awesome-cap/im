package shell

import (
	"errors"
	"fmt"
	"github.com/awesome-cmd/chat/client/ctx"
	"github.com/awesome-cmd/chat/client/render"
	"github.com/awesome-cmd/chat/core/util/http"
	"github.com/awesome-cmd/chat/core/util/json"
	"strings"
)

var servers = []string{
	"https://raw.githubusercontent.com/awesome-cmd/chat/main/servers.json",
	"https://gitee.com/ainilili/chat/raw/main/servers.json",
}

type shell struct {
	ctx *ctx.ChatContext
	position *directory
}

func New(name string) *shell{
	return &shell{
		ctx: ctx.NewContext(name),
		position: root,
	}
}

func (s *shell) Start(){
	for{
		fmt.Printf("[%s@chat %s]# ", s.ctx.Name, s.position.name)
		inputs, err := render.Readline()
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
	serverList = []string{"127.0.0.1:3333|nico"}
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

