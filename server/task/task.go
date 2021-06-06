/**
2 * @Author: Nico
3 * @Date: 2021/6/6 18:31
4 */
package task

import (
	"errors"
	"fmt"
	"github.com/awesome-cmd/chat/core/util/async"
	"github.com/awesome-cmd/chat/server/chats"
	"strconv"
	"strings"
	"time"
)

// cron, only supported ,- * ? /
type cron struct {
	jobs []job
}

type job struct {
	expression string
	fn func()
}

func Start(){
	c := &cron{}
	c.Add("0 */1 * * * ?", func() {
		chatList := chats.Chats()
		for _, chat := range chatList {
			if chat.Deleted || chat.LastActiveTime.Add(2 * time.Hour).Before(time.Now()) {
				chats.DeleteChatPhysically(chat.ID)
			}
		}
	})
	async.Async(c.Start)
}

func (j *job) formats() ([]string, error){
	strs := strings.Split(j.expression, " ")
	if len(strs) < 6 {
		return nil, errors.New(fmt.Sprintf("invalid cron expression %s \n", j.expression))
	}
	return strs, nil
}

func (c *cron) Add(expression string, fn func()){
	c.jobs = append(c.jobs, job{
		expression: expression,
		fn: fn,
	})
}

func (c *cron) Start(){
	for{
		fs := formats()
		for _, job := range c.jobs {
			jfs, err := job.formats()
			if err == nil {
				mat := true
				for i, v := range jfs{
					if ! matched(v, fs[i], i){
						mat = false
						break
					}
				}
				if mat {
					job.fn()
				}
			}
		}
		time.Sleep(time.Millisecond * 998)
	}
}

func formats() []int{
	t := time.Now()
	return []int{t.Second(), t.Minute(), t.Hour(), t.Day(), int(t.Month()), int(t.Weekday()), t.Year()}
}

func matched(exp string, t int, position int) bool{
	if exp == "*" {
		return true
	}
	if strings.Contains(exp, "-") {
		strs := strings.Split(exp, "-")
		l, err := strconv.ParseInt(strs[0], 10, 64)
		if err != nil {
			return false
		}
		r, err := strconv.ParseInt(strs[1], 10, 64)
		if err != nil {
			return false
		}
		return int64(t) >= l && int64(t) <= r
	}
	if exp == "?" && (position == 3 || position == 5) {
		return true
	}
	if strings.Contains(exp, ",") {
		strs := strings.Split(exp, ",")
		for _, v := range strs {
			if matched(v, t, position) {
				return true
			}
		}
		return false
	}
	if strings.Contains(exp, "/") {
		strs := strings.Split(exp, "/")
		r, err := strconv.ParseInt(strs[1], 10, 64)
		if err != nil {
			return false
		}
		s := 0
		if strs[0] != "*" {
			l, err := strconv.ParseInt(strs[0], 10, 64)
			if err != nil {
				return false
			}
			s = int(l)
		}
		return t >= s && (t - s) % int(r) == 0
	}
	return exp == strconv.Itoa(t)
}

