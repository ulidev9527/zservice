package main

import (
	"errors"
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
	e := zservice.NewError("222")
	if errors.Is(e, &zservice.Error{}) {
		zservice.LogError(1)
	} else {
		zservice.LogError(2)
	}

}
