# 通讯协议

- `HTTP + JSON`
    - 对内部使用 `HTTP` 服务，简化服务复杂度
    - 对外使用 `HTTPS` 服务
- `gRPC + Protobuf`

# Gin

- 支持路由功能
- 一进程多服务
- 中间件
- `HTTP` 参数的解析

### 路由分组

```go
v1 := router.Group("/v1")
{
...
}
```

### `HTTP` & `HTTPS`

```go
var eg errgroup.Group
insecureServer := &http.Server{}
secureServer := &http.Server{}

eg.Go(func () error {
err := insecureServer.ListenAndServe()
})

eg.Go(func () error {
err := secureServer.ListenAndServeTLS("cert.pem", "key.pem")
})
```

### 参数解析

- 路径参数 (`path` `uri`)
    - `gin.Default().GET("/user/:name")`
    - `ShouldBindUri` `BindUri`
- 查询字符串参数 (`query` `form`)
    - `/a?name=xxx`
    - `ShouldBindQuery` `BindQuery`
- 表单参数 (`form` `form`)
    - `curl -X POST -F 'username=colin' -F 'password=colin1234' /login`
    - `ShouldBind`
- `HTTP` 头参数 (`header` `header`)
    - `Content-Type: application/json`
    - `ShouldBindHeader` `BindHeader`
- 消息体参数 (`body`)
    - `-d '{"username":"colin","password":"colin1234"}'`
    - `ShouldBindJSON` `BindJSON`

> `ShouldBindWith(obj interface{}, b binding.Binding) error` 参数解析底层函数，绑定失败，只返回错误，不会终止请求
> `MustBindWith(obj interface{}, b binding.Binding) error` 绑定失败，返回错误并终止请求

```go
gin.Default().GET("/:name/:id")

type Person struct {
ID string `uri:"id" binding:"required"`
Name string `uri:"name" binding:"required"`
}

if err := c.ShouldBindUri(&Person); err != nil {
}
```

### 中间件

请求在到达实际处理函数之前，会被一系列加载的中间件进行处理，中间件可以进行一系列的处理

- 中间做成可加载的，通过配置文件指定程序启动加载那些中间件
- 只将一些通用的，必要的功能做成中间件
- 本身支持的中间件
    - `gin.Logger()`
        - 将日志写到 `gin.DefaultWriter` 默认为 `os.Stdout`
    - `gin.Recovery()`
        - 从 `panic` 恢复，并返回 `500`
    - `gin.CustomRecovery(handle gin.RecoveryFunc)`
        - 恢复，同时调用传入的方法
    - `gin.BasicAuth()`
        - 用户名密码的基本认证

```go
router := gin.New()
router.Use(gin.Logger(), gin.Recovery())
```

### 认证和授权

认证确认用户是否有使用系统的权限，授权是确认用户是否有访问某个资源的权限

##### 四种基本认证方式

- `Basic`
    - 将用户名密码进行 `base64` 编码放在 `HTTP Authorization Header`
    - 需要将它和 `SSL` 配合使用
- `Digest`
    - 摘要验证，和基本验证兼容
    - 不使用明文发送密码
    - 防止重放攻击
    - 防止报文篡改
    - 步骤
        - 客户端请求资源
        - 服务认证失败，返回 `WWW-Authenticate` 头，里面包含需要认证的信息
        - 根据头中信息，选择加密算法，并使用密码随机数 `nonce` 计算密码摘要，再次请求服务端
        - 服务端验证摘要，并返回结果
- `OAuth`
    - 允许用户让第三方应用访问该用户资源
    - 密码
        - 用户名密码直接给第三方，第三方换取令牌
    - 隐藏 (前端应用)
        - 提供一个跳转到其他网站的 `URL` 并请求授权
        - 登录成功之后重定向返回并携带令牌
    - 凭借 (无前端)
        - 在命令行向第三方应用请求授权, 携带第三方应用提前颁发的 `secretID` 和 `secretKey`
        - 验证身份，返回结果
    - 授权码
        - 提前申请一个授权码，再使用授权码获取令牌
- `Bearer`

### 策略模式
- `auto`
  - 根据 `http` 头自动选择认证方式
  - `Authorization: Basic XX`
  - `Authorization: Bearer XX`
- `basic`
  - 实现 `Basic` 认证
- `jwt`
  - 实现 `Bearer` 认证
- `cache`
  - `Bearer` 认证的一种实现，采用 `JWT` 其中密钥是从内存中获取

### 验证步骤
- JWT
  - 从请求头中提取 `token`
  - 调用 `ParseWithClaims`
    - `ParseUnverified`
      - 解析 `Token`
        - `Header`
        - `Payload`
        - 加密函数
    - `keyFunc`
      - 获取密钥
  - `KeyExpired`
    - 验证是否过期
  - 设置 `Header` 中的 `username`