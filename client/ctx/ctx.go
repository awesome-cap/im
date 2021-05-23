package ctx

import (
	"errors"
	"fmt"
	"github.com/awesome-cmd/chat/core/model"
	"github.com/awesome-cmd/chat/core/net"
	"github.com/awesome-cmd/chat/core/protocol"
	"github.com/awesome-cmd/chat/core/util/async"
	"github.com/awesome-cmd/chat/core/util/json"
	"log"
	"time"
)

var(
	NotifyTimeOutErr = errors.New("notify time out of 3s")
)

type ChatContext struct {
	Name string `json:"name"`

	conn *net.Conn
	notifies map[int64]chan []byte
}

func NewContext(name string) *ChatContext{
	return &ChatContext{
		Name: name,
		notifies: map[int64]chan []byte{},
	}
}

func (c *ChatContext) Conn() *net.Conn{
	return c.conn
}

func (c *ChatContext) BindConn(conn *net.Conn){
	c.conn = conn
	async.Async(func() {
		err := c.conn.Accept(func(msg protocol.Msg, conn *net.Conn) {
			if ch, ok := c.notifies[msg.ID]; ok{
				ch <- msg.Data
				return
			}
			resp := model.Resp{}
			json.Unmarshal(msg.Data, &resp)
			if resp.Type == "broadcast" && resp.From.ID != c.conn.ID{
				fmt.Printf("%s: %s\n", resp.From.Name, resp.Data)
			}
		})
		if err != nil {
			log.Printf("conn err: %v\n", err)
		}
	})
}

func (c *ChatContext) wait(id int64) ([]byte, error){
	c.notifies[id] = make(chan []byte)
	defer delete(c.notifies, id)
	select {
	case data := <- c.notifies[id]:
		return data, nil
	case <- time.After(time.Second * 3):
		return nil, NotifyTimeOutErr
	}
}

func (c *ChatContext) request(requestData []byte) (*model.Resp, error){
	id, err := c.conn.WriteID(requestData)
	if err != nil{
		return nil, err
	}
	data, err := c.wait(id)
	if err != nil{
		return nil, err
	}
	resp := model.Resp{}
	json.Unmarshal(data, &resp)
	if resp.Code > 100 {
		return nil, errors.New(resp.Msg)
	}
	return &resp, nil
}

func (c *ChatContext) Chats() ([]*model.Chat, error){
	resp, err := c.request(json.Marshal(model.Event{
		Type: "chats",
	}))
	if err != nil {
		return nil, err
	}
	chats := make([]*model.Chat, 0)
	json.Unmarshal(resp.Data, &chats)
	return chats, nil
}

func (c *ChatContext) CreateChat(name string) (*model.Chat, error){
	resp, err := c.request(json.Marshal(model.Event{
		Type: "create",
		Data: name,
	}))
	if err != nil{
		return nil, err
	}
	chat := &model.Chat{}
	json.Unmarshal(resp.Data, chat)
	return chat, nil
}

func (c *ChatContext) ChangeChat(chatId string) error{
	_, err := c.request(json.Marshal(model.Event{
		Type: "change",
		Data: chatId,
	}))
	if err != nil{
		return err
	}
	return nil
}

func (c *ChatContext) Rename(name string) error{
	_, err := c.request(json.Marshal(model.Event{
		Type: "rename",
		Data: name,
	}))
	if err != nil{
		return err
	}
	return nil
}

func (c *ChatContext) Broadcast(msg string) error{
	_, err := c.conn.WriteID(json.Marshal(model.Event{
		Type: "broadcast",
		Data: msg,
	}))
	if err != nil{
		return err
	}
	return nil
}

func (c *ChatContext) Leave() error{
	_, err := c.conn.WriteID(json.Marshal(model.Event{
		Type: "leave",
	}))
	if err != nil{
		return err
	}
	return nil
}