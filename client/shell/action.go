package shell

import (
	"bytes"
	"errors"
	"fmt"
	xnet "github.com/awesome-cmd/chat/core/net"
	"net"
	"os"
	"strings"
)

type action func(s *shell, inputs []byte) (string, error)
type actions map[string]action

var (
	baseActions = actions{}
	rootActions = actions{}
	serverActions =actions{}
)

func distribute(s *shell, inputs []byte, actions ...actions) (string, error){
	args:= strings.Split(strings.ToLower(string(inputs)), " ")
	for _, as := range actions{
		if action, ok := as[args[0]]; ok {
			res, err := action(s, inputs)
			return res, err
		}
	}
	keys := make([]string, 0)
	for _, as := range actions{
		for k := range as {
			keys = append(keys, k)
		}
	}
	return fmt.Sprintf("unknow command %s, reference %v\n", args[0], keys), nil
}

func (a actions) registerAction(key string, act action) actions{
	a[key] = act
	return a
}

func (a actions) registerActions(keys []string, act action) actions{
	for _, key := range keys {
		a[key] = act
	}
	return a
}

func init(){
	initBaseActions()
	initRootActions()
	initServerActions()
}

func initBaseActions(){
	baseActions.registerAction("cd", func(s *shell, inputs []byte) (string, error) {
		args := strings.Split(string(inputs), " ")
		if len(args) != 2 {
			return "", errors.New("cd $dirname")
		}
		dirs := strings.Split(args[1], "/")
		for _, dir := range dirs {
			if dir == ".." {
				if s.position.parent != nil {
					s.position = s.position.parent
				}
			}else if dir == "."{
				// no things to do!
			}else{
				for _, v := range s.position.child {
					if v.name == dir {
						s.position = v
					}
				}
			}
		}
		return "", nil
	}).registerAction("ls", func(s *shell, inputs []byte) (string, error) {
		buffer := bytes.Buffer{}
		for _, v := range s.position.child {
			buffer.WriteString(fmt.Sprintf("%s ", v.name))
		}
		buffer.WriteString("\n")
		return buffer.String(), nil
	}).registerAction("ll", func(s *shell, inputs []byte) (string, error) {
		buffer := bytes.Buffer{}
		for _, v := range s.position.child {
			buffer.WriteString(fmt.Sprintf("%s %s\n", v.name, v.desc))
		}
		return buffer.String(), nil
	}).registerAction("exit", func(s *shell, inputs []byte) (string, error) {
		os.Exit(0)
		return "", nil
	})
}

func initRootActions(){
	rootActions.registerAction("ls", func(s *shell, inputs []byte) (string, error) {
		err := s.refreshServerList()
		if err != nil {
			return "", errors.New(fmt.Sprintf("ls error: %v", err))
		}
		return baseActions["ls"](s, inputs)
	}).registerAction("ll", func(s *shell, inputs []byte) (string, error) {
		err := s.refreshServerList()
		if err != nil {
			return "", errors.New(fmt.Sprintf("ll error: %v", err))
		}
		return baseActions["ll"](s, inputs)
	}).registerAction("cd", func(s *shell, inputs []byte) (s2 string, e error) {
		err := s.refreshServerList()
		if err != nil {
			return "", errors.New(fmt.Sprintf("ll error: %v", err))
		}
		args := strings.Split(string(inputs), " ")
		if len(args) != 2 {
			return "", errors.New("cd $server")
		}
		for _, v := range s.position.child {
			if v.name == args[1] {
				serverInfo := strings.Split(v.desc, "|")
				addr := serverInfo[0]
				tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
				if err != nil{
					return "", errors.New(fmt.Sprintf("cd server error: %v", err))
				}
				conn, err := net.DialTCP("tcp", nil, tcpAddr)
				if err != nil{
					return "", errors.New(fmt.Sprintf("cd server error: %v", err))
				}
				s.ctx.BindConn(xnet.NewConn(conn))
				_ = s.ctx.Rename(s.ctx.Name)
			}
		}
		if s.ctx.Conn() == nil {
			return "", errors.New("cd fail, not found server")
		}
		return baseActions["cd"](s, inputs)
	})
}

func initServerActions(){
	serverActions.registerActions([]string{"ls", "ll"}, func(s *shell, inputs []byte) (s2 string, e error) {
		chats, err := s.ctx.Chats()
		if err != nil{
			return "", err
		}
		builder := bytes.Buffer{}
		for _, v := range chats {
			builder.WriteString(fmt.Sprintf("%d. %s\n", v.ID, v.Name))
		}
		return builder.String(), nil
	}).registerAction("touch", func(s *shell, inputs []byte) (s2 string, e error) {
		args := strings.Split(string(inputs), " ")
		if len(args) != 2 {
			return "", errors.New("mkdir $chatName")
		}
		chat, err := s.ctx.CreateChat(args[1])
		if err != nil {
			return "", errors.New(fmt.Sprintf("create chat err: %v", err))
		}
		return fmt.Sprintf("Create successful ! The new chat id is %d\n", chat.ID), nil
	}).registerAction("vim", func(s *shell, inputs []byte) (s2 string, e error) {
		args := strings.Split(string(inputs), " ")
		if len(args) != 2 {
			return "", errors.New("vim $chatId")
		}
		err := s.ctx.ChangeChat(args[1])
		if err != nil {
			return "", errors.New(fmt.Sprintf("vim chat err: %v", err))
		}
		for {
			msg, _ := s.readline()
			if string(msg) == ":q" {
				if s.ctx.Leave() == nil {
					break
				}
			}
			err := s.ctx.Broadcast(string(msg))
			if err != nil{
				return "", errors.New(fmt.Sprintf("vim err: %v", err))
			}
		}
		return "", nil
	})
}
