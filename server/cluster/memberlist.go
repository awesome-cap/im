/**
2 * @Author: Nico
3 * @Date: 2021/5/25 21:40
4 */
package cluster

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/memberlist"
	"log"
	"os"
	"strconv"
	"sync"
)

var(
	mtx        sync.RWMutex
	items      = map[string]string{}
	broadcasts *memberlist.TransmitLimitedQueue
)

type delegate struct{}

func (d *delegate) NodeMeta(limit int) []byte {
	return []byte{}
}

func (d *delegate) NotifyMsg(b []byte) {
	//if len(b) == 0 {
	//	return
	//}
	//
	//switch b[0] {
	//case 'd': // data
	//	mtx.Lock()
	//	mtx.Unlock()
	//}
	fmt.Println(string(b))
}

func (d *delegate) GetBroadcasts(overhead, limit int) [][]byte {
	return broadcasts.GetBroadcasts(overhead, limit)
}

func (d *delegate) LocalState(join bool) []byte {
	mtx.RLock()
	m := items
	mtx.RUnlock()
	b, _ := json.Marshal(m)
	return b
}

func (d *delegate) MergeRemoteState(buf []byte, join bool) {
	if len(buf) == 0 {
		return
	}
	if !join {
		return
	}
	var m map[string]string
	if err := json.Unmarshal(buf, &m); err != nil {
		return
	}
	mtx.Lock()
	for k, v := range m {
		items[k] = v
	}
	mtx.Unlock()
}

type eventDelegate struct{}

func (ed *eventDelegate) NotifyJoin(node *memberlist.Node) {
	fmt.Println("A node has joined: " + node.String())
}

func (ed *eventDelegate) NotifyLeave(node *memberlist.Node) {
	fmt.Println("A node has left: " + node.String())
}

func (ed *eventDelegate) NotifyUpdate(node *memberlist.Node) {
	fmt.Println("A node was updated: " + node.String())
}

func Init(port int, clusters []string) error{
	hostname, _ := os.Hostname()
	config := memberlist.DefaultLocalConfig()
	config.Name = hostname + "-" + strconv.Itoa(port)
	config.BindPort = port
	config.AdvertisePort = port
	config.Delegate = &delegate{}
	config.Events = &eventDelegate{}

	m, err := memberlist.Create(config)
	if err != nil{
		return err
	}
	suc, err := m.Join(clusters)
	if err != nil {
		return err
	}
	log.Printf("Joined successful cluster num: %d \n", suc)

	broadcasts = &memberlist.TransmitLimitedQueue{
		NumNodes: func() int {
			return m.NumMembers()
		},
		RetransmitMult: 3,
	}
	return nil
}

