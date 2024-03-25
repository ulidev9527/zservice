package main

import "zservice/zservice"

var DBService = zservice.NewService("DBService", func(s *zservice.ZService) {
}, func(s *zservice.ZService) {
	close(s.ChanLock)
})

func main() {

	launcher := zservice.NewService("test", nil, nil)

	launcher.AddService(DBService)

	launcher.Start()

	<-launcher.ChanLock

}
