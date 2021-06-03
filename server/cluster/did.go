/**
2 * @Author: Nico
3 * @Date: 2021/6/2 22:03
4 */
package cluster

import (
	"errors"
	"github.com/awesome-cmd/chat/core/model"
	"github.com/awesome-cmd/chat/core/util/json"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const(
	apply = "incr-apply"
	applyAccess = "incr-apply-access"
	applyRefuse = "incr-apply-refuse"
	increment = "incr-increment"
)

type DID struct {
	id int64

	// 0 finished
	// 1 starting
	state int

	// 0 wait for ack
	// 1 access
	// -1 refuse
	states map[string]int

	notify chan bool
	sync.RWMutex
}

func NextID() (int64, error){
	tryTimes := 1000
	for i := 0; i < tryTimes; i ++{
		id, err := did.next()
		if err == nil {
			return id, nil
		}
	}
	return 0, errors.New("did err")
}

func (d *DID) next() (int64, error){
	d.start()
	defer d.finished()
	d.broadcast(apply, "")
	accessed, err := d.waitForReply()
	if err != nil{
		return 0, err
	}
	if ! accessed {
		return 0, errors.New("refused")
	}
	d.increment()
	d.broadcast(increment, strconv.FormatInt(d.id, 10))
	return d.id, nil
}

func (d *DID) increment(){
	atomic.AddInt64(&d.id, 1)
}

func (d *DID) start(){
	d.Lock()
	d.state = 1
	d.notify = make(chan bool)
	d.states = map[string]int{}
	for _, v := range ml.Members() {
		if v.Name == ml.LocalNode().Name {
			continue
		}
		d.states[v.Name] = 0
	}
}

func (d *DID) finished(){
	d.Unlock()
	d.state = 0
}

func (d *DID) process(event model.Event){
	form := event.From.Name
	switch event.Type {
	case apply:
		t := applyAccess
		if d.state != 0 {
			t = applyRefuse
		}
		d.broadcast(t, strconv.FormatInt(d.id, 10))
	case applyAccess, applyRefuse:
		if _, ok := d.states[form]; ok && d.state == 1{
			if event.Type == applyAccess {
				d.states[form] = 1
			}else{
				d.states[form] = -1
			}
			remoteId, _ := strconv.ParseInt(event.Data, 10, 64)
			if d.id < remoteId {
				d.id = remoteId
			}
			completed := true
			accessed := true
			for _, state := range d.states {
				if state == 0 {
					completed = false
					break
				}
				if state == 2 {
					accessed = false
				}
			}
			if completed {
				d.notify <- accessed
			}
		}
	case increment:
		d.id, _ = strconv.ParseInt(event.Data, 10, 64)
	}
}

func (d *DID) broadcast(t string, data string){
	broadcasts.QueueBroadcast(&broadcast{
		msg: json.Marshal(model.Event{
			Type: t,
			Data: data,
			From: &model.Client{
				Name: ml.LocalNode().Name,
			},
		}),
	})
}

func (d *DID) waitForReply() (bool, error){
	select {
	case accessed := <- d.notify:
		return accessed, nil
	case <- time.After(time.Second * 1):
		return false, errors.New("timeout")
	}
}

