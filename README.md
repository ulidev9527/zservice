# zservice
一个基于 `golang` 的服务开发框架(套件?)

## 目录结构
- `examples` 用例
- `zservice` 框架目录
    - `const.***` 一些常量
    - `fn.***` 一些方法
    - `type.***` 一些对象
- `zserviceex` 基于 `zservice` 扩展的服务/插件，用于其它业务服务调用
    - `dbservice` [gorm](https://gorm.io) + [go-redis](https://redis.uptrace.dev) 数据管理
    - `etcdservice` [etcd](https://etcd.io) 服务发现
    - `grpcservice` [grpc](https://grpc.io) 远程调用过程
    - `ginservice` [gin](https://gin-gonic.com) web服务
    - `nsqservice` [nsq](https://nsq.io) 消息队列
    - `nbioservice` [nbio](https://github.com/lesismal/nbio) tcp/udp/websocket通信