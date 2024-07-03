package main

import (
	"time"
	"zservice/zservice"
)

type TT struct {
}

func (t TT) Now() TT {
	return TT{}
}
func (t *TT) Add() {

}

func init() {
	zservice.Init("test", "1.0.0")
}

func main() {

	data := zservice.ZtimeNow()

	zservice.LogError(
		data.UTC().Truncate(0).Unix(),
		data.UnixMilli(),
		data.UnixMicro(),
		data.UnixNano(),
		data.IsZero(),
		"---",
		time.UnixMilli(1720000993974).String(),
		time.UnixMilli(0).String(),
		// data.Unix(),
		// data.UnixMilli(),
		// data.UnixMicro(),
		// data.UnixNano(),
		// data.IsZero(),

		zservice.MaxInt(1823, 2312, 1, 31, 32, 13, 123, 1, 31, 3),

		zservice.MinInt(1823, 2312, 1, 31, 32, 13, 123, 1, 31, 3),

		zservice.MaxInt64(1823, 2312, 1, 31, 32, 13, 123, 1, 31, 3),

		zservice.MinInt64(1823, 2312, 1, 31, 32, 13, 123, 1, 31, 3),
	)
}
