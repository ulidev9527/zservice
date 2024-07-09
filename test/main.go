package main

import (
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

}
