package events

import (
	"github.com/awesome-cmd/chat/core/model"
	"github.com/awesome-cmd/chat/core/util/json"
	"github.com/awesome-cmd/chat/server/chats"
	"strconv"
)

var processors = map[string]processor{
	"rename": func(id int64, event model.Event, producer func() (int64, error)) *model.Resp {
		if len(event.Data) == 0 {
			return event.Resp(500, nil,  "name can't be empty")
		}
		if len([]rune(event.Data)) > 30{
			return event.Resp(500, nil, "name length must in 30 char")
		}
		event.From.Name = event.Data
		return event.Resp(0, nil, "success")
	},
	"chats": func(id int64, event model.Event, producer func() (int64, error)) *model.Resp {
		return event.Resp(0, json.Marshal(chats.Chats()), "success")
	},
	"change": func(id int64, event model.Event, producer func() (int64, error)) *model.Resp {
		chatId, err := strconv.ParseInt(event.Data, 10, 64)
		if err != nil{
			return event.Resp(500, nil, "please input correct chat number.")
		}
		suc := chats.Change(event.From, chatId)
		if ! suc {
			return event.Resp(500, nil, "chat not exist.")
		}
		return event.Resp(0, nil, "success")
	},
	"leave": func(id int64, event model.Event, producer func() (int64, error)) *model.Resp {
		chats.Leave(event.From)
		return event.Resp(0, nil, "success")
	},
	"create": func(id int64, event model.Event, producer func() (int64, error)) *model.Resp {
		if len(event.Data) == 0 {
			return event.Resp(500, nil, "chat name can't be empty")
		}
		if len([]rune(event.Data)) > 30{
			return event.Resp(500, nil, "chat name length must in 10 char")
		}
		id, err := producer()
		if err != nil {
			return event.Resp(500, nil, err.Error())
		}
		chat, err := chats.Create(event.From, event.Data, id)
		if err != nil{
			return event.Resp(500, nil, err.Error())
		}
		return event.Resp(0, json.Marshal(chat), "success")
	},
	"delete": func(id int64, event model.Event, producer func() (int64, error)) *model.Resp {
		chatId, err := strconv.ParseInt(event.Data, 10, 64)
		if err != nil{
			return event.Resp(500, nil, "please input correct chat number.")
		}
		if chats.Delete(event.From, chatId) {
			return event.Resp(0, nil, "success")
		}
		return event.Resp(500, nil, "fail")
	},
	"broadcast": func(id int64, event model.Event, producer func() (int64, error)) *model.Resp {
		if len(event.Data) == 0 {
			return event.Resp(500, nil, "message can't be empty")
		}
		if len([]rune(event.Data)) > 65535{
			return event.Resp(500, nil, "message length must in 65535 char")
		}
		if ! chats.Exist(event.From.ChatID) {
			return event.Resp(500, nil, "chat not exist")
		}
		chats.Broadcast(event.From, id, event.Resp(0, []byte(event.Data), "success"))
		return nil
	},
}

type processor func(id int64, event model.Event, producer func() (int64, error)) *model.Resp

func Process(msgId int64, event model.Event, producer func() (int64, error)) *model.Resp{
	if processor, ok := processors[event.Type]; ok {
		return processor(msgId, event, producer)
	}
	return event.Resp(500, nil, "unsupported event.")
}
