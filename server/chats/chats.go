package chats

import (
	"github.com/awesome-cmd/chat/core/model"
	"github.com/awesome-cmd/chat/core/net"
	"github.com/awesome-cmd/chat/core/protocol"
	"github.com/awesome-cmd/chat/core/util/async"
	"github.com/awesome-cmd/chat/core/util/json"
	"github.com/awesome-cmd/chat/server/cluster"
	"sort"
	"time"
)

var (
	chats = map[int64]*model.Chat{}
	chatClients = map[int64]map[int64]bool{}
	clients = map[int64]*model.Client{}
	clientConn = map[int64]*net.Conn{}
)

func Change(c *model.Client, chatId int64) bool{
	return Join(c, chatId)
}

func Create(c *model.Client, name string) (*model.Chat, error){
	id, err := cluster.NextID()
	if err != nil{
		return nil, err
	}
	chat := &model.Chat{
		ID: id,
		Name: name,
		Creator: c.Name,
		CreateID: c.ID,
		CreateTime: time.Now(),
	}
	chats[chat.ID] = chat
	chatClients[chat.ID] = map[int64]bool{}
	return chat, nil
}

func Delete(c *model.Client, chatId int64) bool{
	chat := chats[chatId]
	if chat != nil && chat.CreateID == c.ID{
		delete(chats, chatId)
		if chatClients[chat.ID] != nil {
			for cid := range chatClients[chat.ID]{
				if clients[cid] != nil {
					Leave(clients[cid])
				}
			}
		}
		delete(chatClients, chatId)
		return true
	}
	return false
}

func Broadcast(c *model.Client, id int64, msg *model.Resp){
	if c.ChatID > 0 && chatClients[c.ChatID] != nil{
		for clientId := range chatClients[c.ChatID]{
			cid := clientId
			if clients[cid] != nil {
				async.Async(func() {
					Reply(clients[cid], id, msg)
				})
			}
		}
	}
}

func Reply(c *model.Client, id int64, msg *model.Resp){
	conn := clientConn[c.ID]
	if conn != nil {
		_ = conn.Write(protocol.Msg{
			ID: id,
			Data: json.Marshal(msg),
		})
	}
}

func BindClient(conn *net.Conn) {
	conn.ID = time.Now().UnixNano()
	clientConn[conn.ID] = conn
	clients[conn.ID] = &model.Client{
		ID: conn.ID,
	}
}

func Chats() []*model.Chat{
	result := make([]*model.Chat, 0)
	for _, v := range chats {
		result = append(result, v)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	return result
}

func GetChats() map[int64]*model.Chat{
	return chats
}

func SetChats(chatList map[int64]*model.Chat){
	for k, v := range chatList {
		chats[k] = v
		if _, ok := chatClients[k]; ! ok{
			chatClients[k] = map[int64]bool{}
		}
	}
}

func Client(c *net.Conn) *model.Client{
	return clients[c.ID]
}

func Leave(c *model.Client) {
	if _, ok := chatClients[c.ChatID]; ok {
		delete(chatClients[c.ChatID], c.ID)
		c.ChatID = 0
	}
}

func Join(c *model.Client, chatId int64) bool{
	if _, ok := chatClients[chatId]; ok {
		Leave(c)
		chatClients[chatId][c.ID] = true
		c.ChatID = chatId
		return true
	}
	return false
}

func Clean(c *net.Conn){
	client := clients[c.ID]
	if client != nil {
		delete(chatClients[client.ChatID], c.ID)
		delete(clients, c.ID)
	}
	delete(clientConn, c.ID)
}
