package protocol

import (
	"encoding/binary"
	"errors"
	"strings"
)

var (
	prefix = []byte{11, 24, 21, 126, 127}
	lenSize = 4
	idSize = 8
	LengthError = errors.New("Data length error. ")
)

type Msg struct {
	ID int64 `json:"id"`
	Data []byte `json:"data"`
}

func Encode(msg Msg) []byte{
	idBytes, lenBytes := make([]byte, idSize), make([]byte, lenSize)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(msg.Data)))
	binary.BigEndian.PutUint64(idBytes, uint64(msg.ID))
	data := make([]byte, 0)
	data = append(data, prefix...)
	data = append(data, lenBytes...)
	data = append(data, idBytes...)
	data = append(data, )
	return append(data, msg.Data...)
}

func Decode(data []byte) (*Msg, int, error){
	start := strings.Index(string(data), string(prefix))
	if start == -1 {
		return nil, 0, LengthError
	}
	start = start + len(prefix)
	if len(data) < start + lenSize + idSize{
		return nil, 0, LengthError
	}
	size := int(binary.BigEndian.Uint32(data[start:start + lenSize]))
	start = start + lenSize
	id := int64(binary.BigEndian.Uint64(data[start:start + idSize]))
	start = start + idSize
	if len(data) < start + size {
		return nil, 0, LengthError
	}
	return &Msg{
		ID: id,
		Data: data[start: start + size],
	}, start + size, nil
}