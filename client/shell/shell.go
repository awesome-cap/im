package shell

import (
	"bufio"
	"fmt"
	"github.com/awesome-cmd/chat/client/ctx"
	"github.com/awesome-cmd/chat/core/util/http"
	"github.com/awesome-cmd/chat/core/util/json"
	"os"
	"strings"
)

type shell struct {
	in *bufio.Reader
	ctx *ctx.ChatContext
	position *directory
}

func New(name string) *shell{
	return &shell{
		in: bufio.NewReader(os.Stdin),
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
	resp, err := http.Get("https://gitee.com/ainilili/test/raw/master/serverlist.json")
	if err != nil && resp == ""{
		return err
	}
	serverList := make([]string, 0)
	json.Unmarshal([]byte(resp), &serverList)
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