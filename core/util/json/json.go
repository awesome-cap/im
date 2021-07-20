package json

import (
	"github.com/json-iterator/go"
	"unsafe"
)

//var json = jsoniter.ConfigFastest

//
//type bytesDecoder struct {}
//type bytesEncoder struct {}
//
//func (bd *bytesDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
//	str := iter.ReadArray()
//
//	mayBlank, _ := time.Parse(consts.AppTimeFormat, str)
//	now, err := time.ParseInLocation(consts.AppTimeFormat, str, loc)
//
//	if err != nil {
//		*((*time.Time)(ptr)) = time.Unix(0, 0)
//	} else if mayBlank.IsZero() {
//		*((*time.Time)(ptr)) = mayBlank
//	} else {
//		*((*time.Time)(ptr)) = now
//	}
//}
//
//func (be *bytesEncoder) IsEmpty(ptr unsafe.Pointer) bool {
//	ts := *((*time.Time)(ptr))
//	return ts.IsZero()
//}
//
//func (be *bytesEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
//	ts := *((*time.Time)(ptr))
//	if !ts.IsZero() {
//		timestamp := ts.Unix()
//		tm := time.Unix(timestamp, 0)
//		format := tm.Format(consts.AppTimeFormat)
//		stream.WriteString(format)
//	} else {
//		mayBlank, _ := time.Parse(consts.AppTimeFormat, consts.BlankString)
//		stream.WriteString(mayBlank.Format(consts.AppTimeFormat))
//	}

//}

func init(){
	jsoniter.RegisterTypeEncoderFunc("[]uint8", func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
		t := *((*[]byte)(ptr))
		stream.WriteString(string(t))
	}, nil)
	jsoniter.RegisterTypeDecoderFunc("[]uint8", func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		str := iter.ReadString()
		*((*[]byte)(ptr)) = []byte(str)
	})
}

func Marshal(v interface{}) []byte{
	data, _ := jsoniter.Marshal(v)
	return data
}

func Unmarshal(data []byte, v interface{}){
	_ = jsoniter.Unmarshal(data, v)
}