# caolog - Go 日志库

一个高性能的 Go 日志库，支持传统 API 和新的 Event-based API。

## 包结构

- `github.com/CaoStudio/caolog` - 主包（传统 API）
- `github.com/CaoStudio/caolog/core` - 新日志核心（Event-based API）
- `github.com/CaoStudio/caolog/plugin` - 插件（recovery、trace）
- `github.com/CaoStudio/caolog/kratos_log` - Kratos 集成

## 快速开始

### 传统 API

```go
import "github.com/CaoStudio/caolog"

caolog.Info("Hello, World!")
caolog.Error("Something went wrong")
```

### Event-based API (推荐)

```go
import "github.com/CaoStudio/caolog/core"

// 创建 Logger
log := core.NewLogger(nil, core.DebugLevel)

// 链式日志
log.Info().
    Str("user", "john").
    Int("age", 30).
    Msg("User logged in")

// 包级函数
core.Warn().Str("component", "api").Msg("Warning message")
```

## 功能特性

- **链式 API**：类似 zerolog 的 `Info().Str(...).Int(...).Msg(...)` 风格
- **字段支持**：Str、Int、Int32、Int64、Uint、Uint32、Uint64、Bool、Float32、Float64、Any、Time、Dur、Err
- **级别过滤**：Debug、Info、Warn、Error、Fatal、Panic
- **性能优化**：使用 sync.Pool 重用 strings.Builder
- **调用位置自动获取**：正确显示文件名和行号

## 配置

```go
// 设置日志级别
core.WithLevel(core.InfoLevel)

// 设置自定义 Writer
core.WithWriter(os.Stderr)

// 使用自定义 Logger
customLogger := core.NewLogger(os.Stdout, core.DebugLevel)
core.WithLogger(customLogger)
```

## 字段方法示例

```go
log.Info().
    Str("string", "value").
    Int("int", 42).
    Int32("int32", 123).
    Int64("int64", 456).
    Uint("uint", 789).
    Uint32("uint32", 101112).
    Uint64("uint64", 131415).
    Bool("bool", true).
    Float32("float32", 3.14).
    Float64("float64", 2.718).
    Time("time", time.Now()).
    Dur("duration", 5*time.Second).
    Err(errors.New("error")).
    Any("any", map[string]interface{}{"key": "value"}).
    Msg("Complex log message")
```

## 输出格式

```
[2026-04-16 - 01:48:40] [INFO] main.go:12 User logged in user=john age=30
```

## 性能优化

- 使用 `sync.Pool` 重用 `strings.Builder`，减少内存分配
- 预分配缓冲区大小，减少扩容次数
- 使用 `strconv` 直接格式化数字，避免反射

## 插件

### Recovery 插件

自动捕获 panic 并记录日志：

```go
import "github.com/CaoStudio/caolog/plugin"

recovery := plugin.NewRecovery(logger)
defer recovery.Recovery()
```

### Trace 插件

添加追踪信息到日志：

```go
import "github.com/CaoStudio/caolog/plugin"

traceLogger := plugin.NewTraceLogger(logger)
```

## Kratos 集成

```go
import "github.com/CaoStudio/caolog/kratos_log"

kratosLogger := kratoslog.NewKratosLogger(logger)
```

## 许可证

MIT