package fb

import (
	"errors"
	"fmt"
	"zserviceapps/packages/zservice"

	flatbuffers "github.com/google/flatbuffers/go"
)

var pool_builder = zservice.NewPool(func() *flatbuffers.Builder {
	return flatbuffers.NewBuilder(0)
}, func(obj *flatbuffers.Builder) {
	obj.Reset()
})

func GetBuilder() *flatbuffers.Builder     { return pool_builder.Get() }
func PutBuilder(item *flatbuffers.Builder) { pool_builder.Put(item) }

func VerifyTable(buf []byte, offset flatbuffers.UOffsetT, elemSize flatbuffers.VOffsetT) error {
	if int(offset) >= len(buf) {
		return errors.New("offset out of bounds")
	}

	vtableOffset := offset - flatbuffers.UOffsetT(flatbuffers.GetUOffsetT(buf[offset:]))
	if int(vtableOffset) >= len(buf) {
		return errors.New("vtable offset out of bounds")
	}

	vtableLen := flatbuffers.GetVOffsetT(buf[vtableOffset:])
	if vtableLen < elemSize {
		return fmt.Errorf("vtable too short: expected %d, got %d", elemSize, vtableLen)
	}

	return nil
}
