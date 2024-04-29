package main

import (
	"zservice/zservice"
)

var DBService = zservice.NewService(&zservice.ZServiceConfig{
	Name: "DBService",
	OnBeforeStart: func(s *zservice.ZService) {
		zservice.LogDebug("DBService")
	},
	OnStart: func(s *zservice.ZService) {
		s.StartDone()
	},
})

func main() {

	service := zservice.NewService(&zservice.ZServiceConfig{Name: "TestService"})
	restService := zservice.NewRestService(&zservice.ZServiceRESTConfig{
		Name: "RestService",
		Addr: "127.0.0.1:8080",
		OnBeforeStart: func(rs *zservice.RestService) {
			DBService.WaitingDone()
		},
	})

	service.AddService(DBService)
	service.AddService(restService.ZService)

	service.Start()
	service.WaitingDone()

}
