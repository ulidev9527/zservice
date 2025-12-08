protobuf 扩展  


不要有 option 行  
```proto
option go_package = "./pb";
option csharp_namespace = "pb";
```

会处理文件夹下所有的 .proto 文件

### 执行
```
go run tools/protoc_ex/*.go xxxx.proto
go run tools/protoc_ex/*.go  proto文件夹
```