package zservice

import (
	"encoding/binary"

	"google.golang.org/protobuf/proto"
)

// pb 数据转二进制
func ProtobufMustMarshal(pb proto.Message) []byte {
	bt, e := proto.Marshal(pb)
	if e != nil {
		LogError(e)
	}
	return bt
}

// pb 数据转二进制
func ProtobufMustMarshal_String(pb proto.Message) string {
	return string(ProtobufMustMarshal(pb))
}

// 二进制转 pb 数据
func ProtobufMustUnmarshal(bt []byte, pb proto.Message) {
	e := proto.Unmarshal(bt, pb)
	if e != nil {
		LogError(e)
	}
}

// pb 数据转二进制，携带消息ID
func ProtobufMustMarshal_MsgID(msgID uint16, pb proto.Message) []byte {

	id := make([]byte, 2)
	binary.BigEndian.PutUint16(id, msgID)

	return append(id, ProtobufMustMarshal(pb)...)
}

// 二进制转 pb 数据，携带消息ID
func ProtobufMustUnmarshal_MsgID(bt []byte, pb proto.Message) uint16 {

	// 解析消息ID
	idbt := bt[:2]
	msgID := binary.BigEndian.Uint16(idbt)
	bt = bt[2:]

	ProtobufMustUnmarshal(bt[2:], pb)
	return msgID
}

func ProtobufMustUnmarshal_MsgID_BT(bt []byte) (uint16, []byte) {
	// 解析消息ID
	idbt := bt[:2]
	msgID := binary.BigEndian.Uint16(idbt)
	bt = bt[2:]

	return msgID, bt
}
