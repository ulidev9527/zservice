package internal

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice/ex/grpcservice"

	"google.golang.org/grpc"
)

var GrpcService *grpcservice.GrpcService
var Grpc *grpc.Server

func InitGrpc() {
	zauth_pb.RegisterZauthServer(Grpc, &ZauthServer{})
}
