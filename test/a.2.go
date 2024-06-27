package main

import "zservice/zservice"

func getE() *zservice.Error {
	return zservice.NewError("a.2 EEE").SetCode(200)
}
