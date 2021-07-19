package ctx

import (
	"errors"
	"fmt"
	"github.com/awesome-cap/im/core/model"
	"github.com/awesome-cap/im/core/network"
	"github.com/awesome-cap/im/core/protocol"
	"github.com/awesome-cap/im/core/util/async"
	"github.com/awesome-cap/im/core/util/json"
	"github.com/gorilla/websocket"
	"math/rand"
	"net"
	"net/url"
	"strconv"
	"time"
)

var (
	NotifyTimeOutErr = errors.New("notify time out of 3s")
)

type ChatContext struct {
	Name string `json:"name"`

	net       string `json:"net"`
	conn      *network.Conn
	broadcast chan []byte
	server    int
	chatId    int64
	servers   []string
	notifies  map[int64]chan []byte
}

func NewContext(name, net string) *ChatContext {
	return &ChatContext{
		Name:      name,
		net:       net,
		broadcast: make(chan []byte),
		notifies:  map[int64]chan []byte{},
	}
}

func (c *ChatContext) Conn() *network.Conn {
	return c.conn
}

func (c *ChatContext) ListenerBroadcast() {
	async.Async(func() {
		serial := 1
		for {
			msg := <-c.broadcast
			if len(msg) == 1 && msg[0] == 0 {
				break
			}
			resp := model.Resp{}
			json.Unmarshal(msg, &resp)
			if serial > 1 {
				fmt.Printf("\n")
			}
			fmt.Printf("%s %s: %s", resp.Time.Format("2006-01-02 15:04:05"), resp.From.Name, resp.Data)
			serial++
		}
	})
}

func (c *ChatContext) OffListenerBroadcast() {
	c.broadcast <- []byte{0}
}

func (c *ChatContext) Connect(servers []string) error {
	c.servers = servers
	c.server = rand.Intn(len(servers))
	err := c.connect(c.servers, c.server)
	if err != nil {
		return err
	}
	async.Async(func() {
		for {
			_ = c.Rename(c.Name)
			err := c.conn.Accept(func(msg protocol.Msg, conn *network.Conn) {
				if ch, ok := c.notifies[msg.ID]; ok {
					ch <- msg.Data
					delete(c.notifies, msg.ID)
					return
				}
				resp := model.Resp{}
				json.Unmarshal(msg.Data, &resp)
				if resp.Type == "broadcast" && resp.From.ID != c.conn.ID {
					c.broadcast <- msg.Data
				} else if resp.Type == "id" {
					id, _ := strconv.ParseInt(string(resp.Data), 10, 64)
					c.conn.ID = id
				}
			})
			if err != nil && c.conn.State() == 0 {
				n := 10
				if len(c.servers) > n {
					n = len(c.servers)
				}
				_ = c.reconnectN(n)
			}
		}
	})
	return nil
}

func (c *ChatContext) connect(servers []string, server int) error {
	if c.net == "tcp" {
		return c.tcpConnect(servers, server)
	}
	return c.websocketConnect(servers, server)
}

func (c *ChatContext) tcpConnect(servers []string, server int) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", servers[server])
	if err != nil {
		return errors.New(fmt.Sprintf("cd server error: %v", err))
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return errors.New(fmt.Sprintf("cd server error: %v", err))
	}
	if c.conn != nil {
		_ = c.conn.Close()
	}
	c.conn = network.NewConn(protocol.NewTcpReadWriteCloser(conn))
	return nil
}

func (c *ChatContext) websocketConnect(servers []string, server int) error {
	u := url.URL{Scheme: "ws", Host: servers[server], Path: "/im"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return errors.New(fmt.Sprintf("cd server error: %v", err))
	}
	c.conn = network.NewConn(protocol.NewWebsocketReadWriteCloser(conn))
	return nil
}

func (c *ChatContext) reconnect() error {
	c.server++
	c.server = c.server % len(c.servers)
	return c.connect(c.servers, c.server)
}

func (c *ChatContext) reconnectN(n int) error {
	var err error
	fmt.Printf("\nReconnecting")
	for i := 0; i < n; i++ {
		time.Sleep(1 * time.Second)
		fmt.Printf(".")
		err = c.reconnect()
		if err == nil {
			fmt.Printf("Connected\n")
			if c.chatId > 0 {
				_ = c.ChangeChat(c.chatId)
			}
			return nil
		}
	}
	return err
}

func (c *ChatContext) Write(data []byte) (int64, error) {
	id, err := c.conn.WriteID(data)
	return id, err
}

func (c *ChatContext) wait(id int64) ([]byte, error) {
	c.notifies[id] = make(chan []byte)
	defer delete(c.notifies, id)
	select {
	case data := <-c.notifies[id]:
		return data, nil
	case <-time.After(time.Second * 3):
		return nil, NotifyTimeOutErr
	}
}

func (c *ChatContext) request(requestData []byte) (*model.Resp, error) {
	id, err := c.Write(requestData)
	if err != nil {
		return nil, err
	}
	data, err := c.wait(id)
	if err != nil {
		return nil, err
	}
	resp := model.Resp{}
	json.Unmarshal(data, &resp)
	if resp.Code > 100 {
		return nil, errors.New(resp.Msg)
	}
	return &resp, nil
}

func (c *ChatContext) Chats() ([]*model.Chat, error) {
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

func (c *ChatContext) CreateChat(name string) (*model.Chat, error) {
	resp, err := c.request(json.Marshal(model.Event{
		Type: "create",
		Data: name,
	}))
	if err != nil {
		return nil, err
	}
	chat := &model.Chat{}
	json.Unmarshal(resp.Data, chat)
	return chat, nil
}

func (c *ChatContext) DeleteChat(chatId int64) (*model.Chat, error) {
	resp, err := c.request(json.Marshal(model.Event{
		Type: "delete",
		Data: strconv.FormatInt(chatId, 10),
	}))
	if err != nil {
		return nil, err
	}
	chat := &model.Chat{}
	json.Unmarshal(resp.Data, chat)
	return chat, nil
}

func (c *ChatContext) ChangeChat(chatId int64) error {
	_, err := c.request(json.Marshal(model.Event{
		Type: "change",
		Data: strconv.FormatInt(chatId, 10),
	}))
	if err != nil {
		return err
	}
	c.chatId = chatId
	return nil
}

func (c *ChatContext) Rename(name string) error {
	_, err := c.Write(json.Marshal(model.Event{
		Type: "rename",
		Data: name,
	}))
	if err != nil {
		return err
	}
	return nil
}

func (c *ChatContext) Broadcast(msg string) error {
	_, err := c.Write(json.Marshal(model.Event{
		Type: "broadcast",
		Data: msg,
	}))
	if err != nil {
		return err
	}
	return nil
}

func (c *ChatContext) Leave() error {
	_, err := c.Write(json.Marshal(model.Event{
		Type: "leave",
	}))
	if err != nil {
		return err
	}
	c.chatId = 0
	return nil
}
