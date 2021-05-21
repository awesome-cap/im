package chats

import (
	"github.com/awesome-cmd/chat/core/net"
	"github.com/awesome-cmd/chat/core/model"
	"github.com/awesome-cmd/chat/core/util/json"
	"sync/atomic"
)

var (
	id int64 = 0
	chats = map[int64]*model.Chat{}
	chatClients = map[int64]map[int64]bool{}
	clients = map[int64]*model.Client{}
	clientConn = map[int64]*net.Conn{}
)

func Change(c *model.Client, chatId int64){
	if c.ChatID > 0{
		delete(chatClients[c.ChatID], c.ID)
	}
	chatClients[chatId][c.ID] = true
}

func Create(c *model.Client, name string) *model.Chat{
	chat := &model.Chat{
		ID: atomic.AddInt64(&id, 1),
		Name: name,
		Creator: c.Name,
		CreateID: c.ID,
	}
	chats[chat.ID] = chat
	chatClients[chat.ID] = map[int64]bool{}
	return chat
}

func Delete(c *model.Client, chatId int64) bool{
	chat := chats[chatId]
	if chat != nil && chat.CreateID == c.ID{
		delete(chats, chatId)
		if chatClients[chat.ID] != nil {
			for cid := range chatClients[chat.ID]{
				if clients[cid] != nil {
					clients[cid].ChatID = 0
				}
			}
		}
		delete(chatClients, chatId)
		return true
	}
	return false
}

func Broadcast(c *model.Client, msg model.Resp){
	if c.ChatID > 0 && chatClients[c.ChatID] != nil{
		for cid := range chatClients[c.ChatID]{
			if clients[cid] != nil {
				Reply(clients[cid], msg)
			}
		}
	}
}

func Reply(c *model.Client, msg model.Resp){
	conn := clientConn[c.ID]
	if conn != nil {
		_ = conn.Write(json.Marshal(msg))
	}
}

func BindClient(conn *net.Conn) {
	clientConn[conn.ID] = conn
	clients[conn.ID] = &model.Client{
		ID: conn.ID,
	}
}

func Client(conn *net.Conn) *model.Client{
	return clients[conn.ID]
}