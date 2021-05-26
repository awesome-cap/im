/**
2 * @Author: Nico
3 * @Date: 2021/5/25 21:40
4 */
package cluster

import (
	"github.com/awesome-cmd/chat/core/model"
	"github.com/awesome-cmd/chat/core/util/json"
	"github.com/awesome-cmd/chat/server/events"
	"github.com/hashicorp/memberlist"
	"log"
	"os"
	"strconv"
	"sync"
)

var (
	broadcasts *memberlist.TransmitLimitedQueue
	BroadcastEvents = map[string]bool{
		"create": true,
		"delete": true,
		"broadcast": true,
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
	items      map[string]string
}

func (d *delegate) NodeMeta(limit int) []byte {
	return []byte{}
}

func (d *delegate) NotifyMsg(b []byte) {
	event := model.Event{}
	json.Unmarshal(b, event)
	events.Process(0, event)
}

func (d *delegate) GetBroadcasts(overhead, limit int) [][]byte {
	return broadcasts.GetBroadcasts(overhead, limit)
}

func (d *delegate) LocalState(join bool) []byte {
	d.mtx.RLock()
	m := d.items
	d.mtx.RUnlock()
	return json.Marshal(m)
}

func (d *delegate) MergeRemoteState(buf []byte, join bool) {
	if len(buf) == 0 {
		return
	}
	if !join {
		return
	}
	var m map[string]string
	json.Unmarshal(buf, &m)
	d.mtx.Lock()
	for k, v := range m {
		d.items[k] = v
	}
	d.mtx.Unlock()
}

func Broadcast(data []byte) {
	if broadcasts != nil{
		broadcasts.QueueBroadcast(&broadcast{
			msg: data,
		})
	}
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

func Init(port int, seeds []string) error{
	hostname, _ := os.Hostname()
	config := memberlist.DefaultLocalConfig()
	config.Name = hostname + "-" + strconv.Itoa(port)
	config.BindPort = port
	config.AdvertisePort = port
	delegate := &delegate{}
	config.Delegate = delegate
	config.Events = &eventDelegate{}
	m, err := memberlist.Create(config)
	if err != nil{
		return err
	}
	if len(seeds) > 0 {
		suc, err := m.Join(seeds)
		if err != nil {
			return err
		}
		log.Printf("Joined successful cluster num: %d \n", suc)
	}
	broadcasts = &memberlist.TransmitLimitedQueue{
		NumNodes: func() int {
			return m.NumMembers()
		},
		RetransmitMult: 1,
	}
	return nil
}

