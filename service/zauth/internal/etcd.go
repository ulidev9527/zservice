package internal

import (
	"zservice/zservice/ex/etcdservice"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var EtcdService *etcdservice.EtcdService
var Etcd *clientv3.Client

func InitEtcd() {

}
