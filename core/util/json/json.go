package json

import "encoding/json"

func Marshal(v interface{}) []byte{
	data, _ := json.Marshal(v)
	return data
}

func Unmarshal(data []byte, v interface{}){
	_ = json.Unmarshal(data, v)
}