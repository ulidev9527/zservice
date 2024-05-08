# zservice
一个基于 `golang` 的服务开发框架(套件?)

## 目录结构
- `internal` 内置模块/服务
    - `dbservice` 数据相关处理的服务
- `service` 服务目录, 所有的服务都在这里面进行.
- `zservice` 框架目录

## 开发计划
- ~~建立仓库~~
- `zservice`
    - 初步,测试和探索阶段
        - ~~框架结构搭建~~
        - ~~gorm~~
        - ~~go-redis~~
        - ~~zerolog~~
        - ~~gin~~
        - etcd
        - nsq
- `zconfig`
    - 配置服务

## 主要集成
- [gorm](https://gorm.io)
- [go-redis](https://redis.uptrace.dev)
- [gin](https://gin-gonic.com)
- [zerolog](https://github.com/rs/zerolog)
- [etcd](https://etcd.io)
- [nsq](https://nsq.io)
- [grpc](https://grpc.io)