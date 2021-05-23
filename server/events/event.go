package events

import (
	"github.com/awesome-cmd/chat/core/model"
	"github.com/awesome-cmd/chat/core/net"
	"github.com/awesome-cmd/chat/core/protocol"
	"github.com/awesome-cmd/chat/core/util/json"
	"github.com/awesome-cmd/chat/server/chats"
	"strconv"
)

var processors = map[string]processor{
	"rename": func(id int64, event model.Event) *model.Resp {
		if len(event.Data) < 0 {
			return event.Resp(500, nil,  "name can't be empty")
		}
		if len([]rune(event.Data)) > 30{
			return event.Resp(500, nil, "name length must in 30 char")
		}
		event.From.Name = event.Data
		return event.Resp(0, nil, "success")
	},
	"chats": func(id int64, event model.Event) *model.Resp {
		return event.Resp(0, json.Marshal(chats.Chats()), "success")
	},
	"change": func(id int64, event model.Event) *model.Resp {
		chatId, err := strconv.ParseInt(event.Data, 10, 64)
		if err != nil{
			return event.Resp(500, nil, "please input correct chat number.")
		}
		chats.Change(event.From, chatId)
		return event.Resp(0, nil, "success")
	},
	"leave": func(id int64, event model.Event) *model.Resp {
		chats.Leave(event.From)
		return event.Resp(0, nil, "success")
	},
	"create": func(id int64, event model.Event) *model.Resp {
		if len(event.Data) < 0 {
			return event.Resp(500, nil, "chat name can't be empty")
		}
		if len([]rune(event.Data)) > 30{
			return event.Resp(500, nil, "chat name length must in 10 char")
		}
		chat := chats.Create(event.From, event.Data)
		return event.Resp(0, json.Marshal(chat), "success")
	},
	"delete": func(id int64, event model.Event) *model.Resp {
		chatId, err := strconv.ParseInt(event.Data, 10, 64)
		if err != nil{
			return event.Resp(500, nil, "please input correct chat number.")
		}
		if chats.Delete(event.From, chatId) {
			return event.Resp(0, nil, "success")
		}
		return event.Resp(500, nil, "fail")
	},
	"broadcast": func(id int64, event model.Event) *model.Resp {
		if len(event.Data) < 0 {
			return event.Resp(500, nil, "message can't be empty")
		}
		if len([]rune(event.Data)) > 65535{
			return event.Resp(500, nil, "message length must in 65535 char")
		}
		chats.Broadcast(event.From, id, event.Resp(0, []byte(event.Data), "success"))
		return nil
	},
}

type processor func(id int64, event model.Event) *model.Resp

func Process(msg protocol.Msg, c *net.Conn) *model.Resp{
	event := model.Event{}
	json.Unmarshal(msg.Data, &event)
	event.From = chats.Client(c)
	if processor, ok := processors[event.Type]; ok{
		return processor(msg.ID, event)
	}
	return event.Resp(500, nil, "unsupported event.")
}
