# 添加上下文

返回一个新的错误，记录堆栈，同时在原有基础上添加上下文

```go
import (
stderrors "errors"
)

err := stderrors.New("original error")
err = errors.Wrap(err, "context")
```

# 检索错误原因

```go
err := errors.Cause(err).(type)

```

# 格式化打印错误

- `%s`
    - 打印错误，如果有 Cause 递归
- `%v`
    - 同上
- `%-v`
    - 打印错误链中最后一个错误
- `%+v`
    - 打印所有错误
- `%#-v`
    - `json` 格式打印最后一个错误
- `%#+v`

# 检索错误

```go
for err, ok := err.(stackTracer); ok {
for _, f := range err.StackTrace() {
}
}
```

# 业务错误代码

```go
type withCode struct {
    err   error // error 错误
    code  int // 业务错误码
    cause error // cause error
    *stack // 错误堆栈
}
errors.WrapCode(err, "get user failed.")
```