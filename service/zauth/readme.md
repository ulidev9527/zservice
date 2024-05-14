# zauth
权限系统/认证系统  
`protoc zauth.proto --go_out=. --go-grpc_out=.`  


## DB

### PermissionTable
用于存储具体权限   

字段|类型|说明
-|-|-
Name | `string` | 权限显示名称
Action | `string` | 权限动作，详情查看[PermissionTable.Action](#permissiontableaction)

#### `PermissionTable.Action`
存储所有权限动作  
权限动作由: [服务名]/[协议]/功能 组成，使用 `/` 进行分割，其中第一个和第二个为固定格式.  

**示例**
- `zauth/rpc/CheckAuth`, 表示 `zauth` 服务的 `rpc` 下的 `checkauth` 操作。
- `zauth/http_get/xxx`, 表示 `zauth` 服务的 `http_get` 协议的 `xxx` 操作。
- `auth/http_post/xxx/xxx`, 可以多层操作