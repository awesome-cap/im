package model

type Msg struct {
	Data []byte `json:"data"`
}

type Event struct {
	Type string      `json:"type"`
	Data string      `json:"data"`
	From Client `json:"from"`
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
	Data interface{} `json:"data"`
	Msg string `json:"msg"`
}