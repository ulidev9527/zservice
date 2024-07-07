# zservice
一个基于 `golang` 的服务开发框架(套件?)

## 目录结构
- `service` 服务目录, 所有的服务都在这里面进行.
    - `zauth` 一个授权服务，用于支撑其它服务的授权管理/登陆管理/资源上传等，使其他服务尽量保持无状态运行。
- `zservice` 框架目录
    - `const.***` 一些常量
    - `fn.***` 一些方法
    - `type.***` 一些对象
- `zserviceex` 基于 `zservice` 扩展的服务/插件，用于其它业务服务调用
    - `dbservice` 数据相关处理的服务,集成 [gorm](https://gorm.io) 和 [go-redis](https://redis.uptrace.dev),用于数据存储。
    - `etcdservice` [etcd](https://etcd.io)集成使用
    - `grpcservice` [grpc](https://grpc.io)集成使用
    - `ginservice` [gin](https://gin-gonic.com)集成使用
    - `nsqservice` [nsq](https://nsq.io) 集成使用

## 主要集成
- [gorm](https://gorm.io)
- [go-redis](https://redis.uptrace.dev)
- [gin](https://gin-gonic.com)
- [zerolog](https://github.com/rs/zerolog)
- [etcd](https://etcd.io)
- [nsq](https://nsq.io)
- [grpc](https://grpc.io)