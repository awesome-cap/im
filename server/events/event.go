package events

import (
	"github.com/awesome-cmd/chat/core/model"
	"github.com/awesome-cmd/chat/core/net"
	"github.com/awesome-cmd/chat/server/chats"
	"strconv"
)

var processors = map[string]processor{
	"rename": func(data string, c *model.Client) model.Resp {
		if len(data) < 0 {
			return resp(500, nil, "name can't be empty")
		}
		if len([]rune(data)) > 30{
			return resp(500, nil, "name length must in 30 char")
		}
		c.Name = data
		return resp(0, nil, "success")
	},
	"change": func(data string, c *model.Client) model.Resp {
		chatId, err := strconv.ParseInt(data, 10, 64)
		if err != nil{
			return resp(500, nil, "please input correct chat number.")
		}
		chats.Change(c, chatId)
		return resp(0, nil, "success")
	},
	"create": func(data string, c *model.Client) model.Resp {
		if len(data) < 0 {
			return resp(500, nil, "chat name can't be empty")
		}
		if len([]rune(data)) > 30{
			return resp(500, nil, "chat name length must in 10 char")
		}
		chat := chats.Create(c, data)
		return resp(0, chat, "success")
	},
	"delete": func(data string, c *model.Client) model.Resp {
		chatId, err := strconv.ParseInt(data, 10, 64)
		if err != nil{
			return resp(500, nil, "please input correct chat number.")
		}
		if chats.Delete(c, chatId) {
			return resp(0, nil, "success")
		}
		return resp(500, nil, "fail")
	},
	"broadcast": func(data string, c *model.Client) model.Resp {
		if len(data) < 0 {
			return resp(500, nil, "message can't be empty")
		}
		if len([]rune(data)) > 65535{
			return resp(500, nil, "message length must in 65535 char")
		}
		chats.Broadcast(c, model.Resp{
			Code: 1,
			Data: data,
			Msg: "message",
		})
		return resp(0, nil, "success")
	},
}

type processor func(data string, c *model.Client) model.Resp

func Process(event model.Event, c *net.Conn) model.Resp{
	if processor, ok := processors[event.Type]; ok{
		client := chats.Client(c)
		if client != nil {
			return processor(event.Data, client)
		}
		return resp(500, nil, "connect closed.")
	}
	return resp(500, nil, "unsupported event.")
}

func resp(code int, data interface{}, msg string) model.Resp{
	resp := model.Resp{
		Code: code,
		Msg: msg,
		Data: data,
	}
	return resp
}
