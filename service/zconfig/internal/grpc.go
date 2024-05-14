package internal

import (
	"zservice/service/zconfig/zconfig_pb"

	"google.golang.org/grpc"
)

var Grpc *grpc.Server

func InitGrpc() {
	zconfig_pb.RegisterZconfigServer(Grpc, &ZconfigServer{})
}
