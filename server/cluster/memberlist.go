/**
2 * @Author: Nico
3 * @Date: 2021/5/25 21:40
4 */
package cluster

import (
	"github.com/awesome-cap/im/core/model"
	"github.com/awesome-cap/im/core/util/json"
	"github.com/awesome-cap/im/server/chats"
	"github.com/awesome-cap/im/server/events"
	"github.com/hashicorp/memberlist"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

var (
	ml *memberlist.Memberlist
	broadcasts *memberlist.TransmitLimitedQueue
	did = &DID{}

	BroadcastEvents = map[string]bool{
		"broadcast": true,
		"delete": true,
	}
	didEvents = map[string]bool{
		apply: true,
		applyAccess: true,
		applyRefuse: true,
		increment: true,
		incremented: true,
	}
)

type broadcast struct {
	msg    []byte
	notify chan<- struct{}
}

func (b *broadcast) Invalidates(other memberlist.Broadcast) bool {
	return false
}

func (b *broadcast) Message() []byte {
	return b.msg
}

func (b *broadcast) Finished() {
	if b.notify != nil {
		close(b.notify)
	}
}

type delegate struct{
	mtx        sync.RWMutex
}

func (d *delegate) NodeMeta(limit int) []byte {
	return []byte{}
}

func (d *delegate) NotifyMsg(b []byte) {
	event := model.Event{}
	json.Unmarshal(b, &event)
	if BroadcastEvents[event.Type] {
		go events.Process(0, event, NextID)
	}else if didEvents[event.Type] {
		go did.process(event)
	}
}

func (d *delegate) GetBroadcasts(overhead, limit int) [][]byte {
	return broadcasts.GetBroadcasts(overhead, limit)
}

func (d *delegate) LocalState(join bool) []byte {
	d.mtx.RLock()
	defer d.mtx.RUnlock()
	return json.Marshal(chats.GetChats())
}

func (d *delegate) MergeRemoteState(buf []byte, join bool) {
	chatList := map[int64]*model.Chat{}
	json.Unmarshal(buf, &chatList)
	d.mtx.RLock()
	defer d.mtx.RUnlock()
	chats.MergeRemoteChats(chatList)
}

func Broadcast(data []byte) {
	if broadcasts != nil{
		broadcasts.QueueBroadcast(&broadcast{
			msg: data,
		})
	}
}

func LocalNode() *memberlist.Node{
	return ml.LocalNode()
}

type eventDelegate struct{}

func (ed *eventDelegate) NotifyJoin(node *memberlist.Node) {
	log.Printf("A node has joined: %s\n", node.String())
}

func (ed *eventDelegate) NotifyLeave(node *memberlist.Node) {
	log.Printf("A node has left: %s\n", node.String())
}

func (ed *eventDelegate) NotifyUpdate(node *memberlist.Node) {
	log.Printf("A node was updated: %s\n", node.String())
}

func localName(port int) string{
	interfaces, _ := net.Interfaces()
	return interfaces[0].HardwareAddr.String() + "-" + strconv.Itoa(port) + "-" + strconv.FormatInt(time.Now().UnixNano(), 10)
}

func Start(port int, seeds []string) error{
	config := memberlist.DefaultLocalConfig()
	config.Name = localName(port)
	config.BindPort = port
	config.AdvertisePort = port
	config.PushPullInterval = time.Second * 3
	var err error
	ml, err = memberlist.Create(config)
	if err != nil{
		return err
	}
	broadcasts = &memberlist.TransmitLimitedQueue{
		NumNodes: func() int {
			return ml.NumMembers()
		},
		RetransmitMult: 1,
	}
	delegate := &delegate{}
	config.Delegate = delegate

	if len(seeds) > 0 {
		suc, err := ml.Join(seeds)
		if err != nil {
			return err
		}
		log.Printf("Joined successful cluster num: %d \n", suc)
	}
	return nil
}

