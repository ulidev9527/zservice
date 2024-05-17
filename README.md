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
        - ~~etcd~~
        - ~~grpc~~
        - ~~nsq~~
        - ~~链路~~
        - token验证
- `zsms` 短信服务
    - ~~三方服务发送短信~~
- `zconfig` 配置服务
    - ~~三方服务文件配置读取~~
    - ~~三方服务结合 `nsq` 进行配置更改监听和同步更新~~

## 主要集成
- [gorm](https://gorm.io)
- [go-redis](https://redis.uptrace.dev)
- [gin](https://gin-gonic.com)
- [zerolog](https://github.com/rs/zerolog)
- [etcd](https://etcd.io)
- [nsq](https://nsq.io)
- [grpc](https://grpc.io)

## Bug
- `debug` 下 `grpc.server` 中的 `lease` 超时会断开失效，需要重连机制

### 预留环境变量
*数组类型使用 `,` 进行分割*
模块|名称|类型|示例|说明
-|-|-|-|-
`zservice`
-|`ZSERVICE_NAME`       |`string` |`service_name`|服务名称
-|`ZSERVICE_VERSION`    |`string` |`0.1.0`| 服务版本,`仅在zservice.Init(name, version)生效，不受其他环境变量影响`
-|`ZSERVICE_REMOTE_ENV_ADDR`| `string` |`"http://127.0.0.1/config"` | 远程环境变量地址
-|`ZSERVICE_REMOTE_ENV_AUTH` | `string` |`"授权字符串"` | 远程环境变量授权, 会添加参数为:`?auth=***`
-|`ZSERVICE_FILES_ENV`|`[]string`|`static/a.env,static/b.env` | 环境变量文件

## 模块
### `zservice`
项目启动模块、框架结构