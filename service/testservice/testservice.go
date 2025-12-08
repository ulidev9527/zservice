package main

import "zserviceapps/packages/zservice"

func main() {

	zservice.NewService(zservice.ServiceOptions{

		Name: "testservice",
	})

	zservice.Start()

	zservice.WaitStart()

	zservice.WaitStop()

}
