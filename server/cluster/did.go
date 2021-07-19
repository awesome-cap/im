/**
2 * @Author: Nico
3 * @Date: 2021/6/2 22:03
4 */
package cluster

import (
	"errors"
	"github.com/awesome-cap/im/core/model"
	"github.com/awesome-cap/im/core/util/json"
	"sync"
	"sync/atomic"
	"time"
)

const(
	apply = "incr-apply"
	applyAccess = "incr-apply-access"
	applyRefuse = "incr-apply-refuse"
	increment = "incr-increment"
	incremented = "incr-incremented"
)

type DID struct {
	ID int64 `json:"id"`
	StartTime int64 `json:"startTime"`
	Node string `json:"node"`

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
	tryTimes := 3
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
	if len(d.states) > 0 {
		d.broadcast(apply)
		accessed, err := d.waitForReply()
		if err != nil{
			return 0, err
		}
		if ! accessed {
			return 0, errors.New("refused")
		}
	}
	d.state = 2
	d.increment()
	d.broadcast(increment)
	//if len(d.states) > 0{
	//	d.broadcast(increment)
	//	_, err := d.waitForReply()
	//	if err != nil{
	//		fmt.Printf("waitForReply2 err: %v \n", err)
	//		d.decrement()
	//		return 0, err
	//	}
	//}
	return d.ID, nil
}

func (d *DID) increment(){
	atomic.AddInt64(&d.ID, 1)
}

func (d *DID) decrement(){
	atomic.AddInt64(&d.ID, -1)
}

func (d *DID) start(){
	d.Lock()
	d.state = 1
	d.StartTime = time.Now().UnixNano()
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
	remote := DID{}
	json.Unmarshal([]byte(event.Data), &remote)
	switch event.Type {
	case apply:
		t := applyAccess
		if d.state != 0 && d.StartTime <= remote.StartTime{
			t = applyRefuse
		}
		d.broadcast(t)
	case applyAccess, applyRefuse:
		if _, ok := d.states[remote.Node]; ok && d.state == 1{
			if event.Type == applyAccess {
				d.states[remote.Node] = 1
				if d.ID < remote.ID {
					d.ID = remote.ID
				}
			}else{
				d.states[remote.Node] = -1
			}
			completed := true
			accessed := true
			for _, state := range d.states {
				if state == 0 {
					completed = false
					break
				}
				if state == -1 {
					accessed = false
				}
			}
			if completed {
				d.notify <- accessed
			}
		}
	case increment:
		d.increment()
		//d.broadcast(incremented)
	case incremented:
		if  _, ok := d.states[remote.Node]; ok && d.state == 2 {
			d.states[remote.Node] = 2
			completed := true
			for _, state := range d.states {
				if state == 0 {
					completed = false
					break
				}
			}
			if completed {
				d.notify <- true
			}
		}
	}
}

func (d *DID) broadcast(t string){
	d.Node = ml.LocalNode().Name
	broadcasts.QueueBroadcast(&broadcast{
		msg: json.Marshal(model.Event{
			Type: t,
			Data: string(json.Marshal(d)),
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

