package protocol

import (
	"encoding/binary"
	"errors"
	"strings"
)

var (
	prefix = []byte{11, 24, 21, 126, 127}
	length = 4
	LengthError = errors.New("Data length error. ")
)

func Encode(data []byte) []byte{
	buf := make([]byte, length)
	binary.BigEndian.PutUint32(buf, uint32(len(data)))
	msg := make([]byte, 0)
	msg = append(msg, prefix...)
	msg = append(msg, buf...)
	return append(msg, data...)
}

func Decode(data []byte) ([]byte, int, error){
	start := strings.Index(string(data), string(prefix))
	if start == -1 {
		return nil, 0, LengthError
	}
	start = start + len(prefix)
	if len(data) <= start {
		return nil, 0, LengthError
	}
	size := int(binary.BigEndian.Uint32(data[start:start + length]))
	start = start + length
	if len(data) < start + size {
		return nil, 0, LengthError
	}
	return data[start: start + size], start + size, nil
}