package internal

import (
	"zservice/service/zauth/zauth_pb"

	"google.golang.org/grpc"
)

var Grpc *grpc.Server

func InitGrpc() {
	zauth_pb.RegisterZauthServer(Grpc, &ZauthServer{})
}
