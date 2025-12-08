# zservice
一个基于 `golang` 的服务开发模版

## 目录结构
- `packages` 用例
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


pb 编译:  
protoc --go_out=. --go-grpc_out=. --csharp_out=. xxx.proto


### 扩
- pb 生成
    - `protoc_ex.sh xxxx`
    - `xxxx`表示 `.protoc` 文件或者文件目录
    - 会自动编译目录下所有的 `.protoc` 文件
    - protoc --go_out=. --go-grpc_out=. --csharp_out=. xxx.proto

- luban 配置生成
    - `luban_ex.sh xxxx`
        - `xxxx`表示 `luban` 配置根目录下的文件或者 `luban` 配置根目录

- fb 生成        