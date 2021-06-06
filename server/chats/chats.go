package chats

import (
	"github.com/awesome-cmd/chat/core/model"
	"github.com/awesome-cmd/chat/core/net"
	"github.com/awesome-cmd/chat/core/protocol"
	"github.com/awesome-cmd/chat/core/util/async"
	"github.com/awesome-cmd/chat/core/util/json"
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

func Create(c *model.Client, name string, id int64) (*model.Chat, error){
	chat := &model.Chat{
		ID: id,
		Name: name,
		Creator: c.Name,
		CreateID: c.ID,
		CreateTime: time.Now(),
		LastActiveTime: time.Now(),
	}
	chats[chat.ID] = chat
	chatClients[chat.ID] = map[int64]bool{}
	return chat, nil
}

func Delete(c *model.Client, chatId int64) bool{
	chat := chats[chatId]
	if chat != nil && chat.CreateID == c.ID{
		DeleteChatLogically(chatId)
		return true
	}
	return false
}

func DeleteChatLogically(chatId int64){
	defer flushLastActiveTime(chatId)
	if chat, ok := chats[chatId]; ok {
		chat.Deleted = true
		deleteChatClients(chatId)
	}
}

func DeleteChatPhysically(chatId int64){
	defer flushLastActiveTime(chatId)
	delete(chats, chatId)
	deleteChatClients(chatId)
}

func deleteChatClients(chatId int64){
	if chatClients[chatId] != nil {
		for cid := range chatClients[chatId]{
			if clients[cid] != nil {
				Leave(clients[cid])
			}
		}
	}
	delete(chatClients, chatId)
}

func Broadcast(c *model.Client, id int64, msg *model.Resp){
	defer flushLastActiveTime(c.ChatID)
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
		if ! v.Deleted {
			result = append(result, v)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	return result
}

func Exist(chatId int64) bool{
	return chats[chatId] != nil
}

func GetChats() map[int64]*model.Chat{
	return chats
}

func Client(c *net.Conn) *model.Client{
	return clients[c.ID]
}

func Leave(c *model.Client) {
	defer flushLastActiveTime(c.ChatID)
	if _, ok := chatClients[c.ChatID]; ok {
		delete(chatClients[c.ChatID], c.ID)
		c.ChatID = 0
	}
}

func Join(c *model.Client, chatId int64) bool{
	defer flushLastActiveTime(c.ChatID)
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

func MergeRemoteChats(chatList map[int64]*model.Chat){
	for k, v := range chatList {
		chat := chats[v.ID]
		if chat != nil {
			if v.LastActiveTime.After(chat.LastActiveTime) {
				chat.LastActiveTime = v.LastActiveTime
			}
			if v.Deleted {
				chat.Deleted = true
			}
		}else{
			chats[v.ID] = v
			if _, ok := chatClients[k]; ! ok{
				chatClients[k] = map[int64]bool{}
			}
		}
	}
}

func flushLastActiveTime(chatId int64){
	if chat, ok := chats[chatId]; ok {
		if ! chat.Deleted{
			chat.LastActiveTime = time.Now()
		}
	}
}
