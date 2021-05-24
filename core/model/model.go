package model

import "time"

type Msg struct {
	Data []byte `json:"data"`
}

type Event struct {
	Type string      `json:"type"`
	Data string      `json:"data"`
	From *Client 	 `json:"from"`
}

func (e Event) Resp(code int, data []byte, msg string) *Resp{
	return &Resp{
		Code: code,
		Data: data,
		Msg: msg,
		Type: e.Type,
		From: e.From,
		Time: time.Now(),
	}
}

type Client struct {
	ID int64 `json:"id"`
	Name string `json:"name"`
	ChatID int64 `json:"chatId"`
}

type Chat struct {
	ID int64 `json:"id"`
	Name string `json:"name"`
	Creator string `json:"creator"`
	CreateID int64 `json:"createId"`
}

type Resp struct {
	Code int `json:"code"`
	Type string `json:"type"`
	Data []byte `json:"data"`
	Msg string `json:"msg"`
	From *Client 	 `json:"from"`
	Time time.Time `json:"time"`
}