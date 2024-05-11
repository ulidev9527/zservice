# zconfig
配置服务  
主要负责其他微服务的配置管理，提供给其它服务使用  
其他服务配合 `ZServiceConfig` 里面的 `RemoteEnvAddr` 参数使用  

pb: `protoc zconfig.proto --go_out=. --go-grpc_out=.`