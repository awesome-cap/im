package protocol

import (
	"encoding/binary"
	"errors"
	"io"
)

var (
	lenSize = 4
	idSize = 8
	LengthError = errors.New("Data length error. ")
)

type Msg struct {
	ID int64 `json:"id"`
	Data []byte `json:"data"`
}

func ReadUint32(reader io.Reader) (uint32, error) {
	data := make([]byte, 4)
	_, err := io.ReadFull(reader, data)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(data), nil
}

func ReadUint64(reader io.Reader) (uint64, error) {
	data := make([]byte, 8)
	_, err := io.ReadFull(reader, data)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint64(data), nil
}

func Encode(msg Msg) []byte{
	idBytes, lenBytes := make([]byte, idSize), make([]byte, lenSize)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(msg.Data)))
	binary.BigEndian.PutUint64(idBytes, uint64(msg.ID))
	data := make([]byte, 0)
	data = append(data, idBytes...)
	data = append(data, lenBytes...)
	return append(data, msg.Data...)
}

func Decode(r io.Reader) (*Msg, error){
	id, err := ReadUint64(r)
	if err != nil{
		return nil, err
	}
	l, err := ReadUint32(r)
	if err != nil{
		return nil, err
	}
	dataBytes := make([]byte, l)
	_, err = io.ReadFull(r, dataBytes)
	if err != nil{
		return nil, err
	}
	return &Msg{
		ID: int64(id),
		Data: dataBytes,
	}, nil
}