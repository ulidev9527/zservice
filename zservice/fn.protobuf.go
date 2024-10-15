package zservice

import "google.golang.org/protobuf/proto"

func ProtobufMustMarshal(pb proto.Message) []byte {
	bt, e := proto.Marshal(pb)
	if e != nil {
		LogError(e)
	}
	return bt
}

func ProtobufMustUnmarshal(bt []byte, pb proto.Message) {
	e := proto.Unmarshal(bt, pb)
	if e != nil {
		LogError(e)
	}
}
