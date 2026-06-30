# glthrift Wrapper Design

## 背景

`glthrift` 当前包含 `proto.thrift` 及 Apache Thrift 生成的 Go 代码，生成服务为 `protoservice.ProtoRpcService`。生成代码已经提供底层 client、processor、请求结构和响应结构，但使用方仍需要手动拼装 transport、protocol、server、连接关闭和超时配置。

本次目标是在 `glthrift` 中补齐两层封装：

- 通用 Apache Thrift 服务端和客户端底座，适配任意 thrift 生成服务。
- 当前 `ProtoRpcService` 的便捷服务端和客户端，直接服务已有 `proto.thrift`。

## 目标

- 提供可复用的 thrift server 封装，接收任意 `thrift.TProcessor`。
- 提供可复用的 thrift client 连接封装，负责打开和关闭 transport。
- 提供 `ProtoRpcService` 专用便捷 API，减少业务代码中的 thrift 拼装代码。
- 支持常见配置：监听地址、连接地址、协议类型、传输类型、连接超时、socket 超时。
- 所有导出 API 提供中文 GoDoc 注释。
- 通过单元测试覆盖 twoway、oneway、连接失败和关闭行为。

## 非目标

- 不重新生成 thrift 代码。
- 不修改 `proto.thrift` 的服务定义。
- 不内置服务发现、负载均衡、熔断、重试或连接池。
- 不绑定业务日志、配置中心或特定框架。

## API 设计

### 通用服务端

新增 `Server` 类型和 `ServerConfig`：

```go
type ServerConfig struct {
    Addr          string
    Protocol      Protocol
    Transport     Transport
    SocketTimeout time.Duration
}

type Server struct {
    // 内部持有 thrift.TServer
}
```

核心构造函数：

```go
func NewServer(addr string, processor thrift.TProcessor, opts ...ServerOption) (*Server, error)
```

服务端行为：

- `Serve()` 阻塞运行 thrift server。
- `Stop()` 停止 thrift server。
- `Addr()` 返回配置的监听地址。
- 监听地址必须显式传入，空地址返回错误。
- 默认协议为 binary protocol。
- 默认传输为 buffered transport。

### 通用客户端连接

新增 `ClientConn` 类型和 `ClientConfig`：

```go
type ClientConfig struct {
    Addr          string
    Protocol      Protocol
    Transport     Transport
    ConnectTimeout time.Duration
    SocketTimeout  time.Duration
}

type ClientConn struct {
    // 内部持有 thrift.TTransport 和 thrift.TProtocolFactory
}
```

核心构造函数：

```go
func NewClientConn(addr string, opts ...ClientOption) (*ClientConn, error)
```

客户端连接行为：

- 构造时打开 thrift transport；失败返回 error。
- `Transport()` 返回已打开的 `thrift.TTransport`。
- `ProtocolFactory()` 返回协议工厂。
- `Close()` 关闭连接，允许重复调用。
- `Addr()` 返回连接地址。

### 协议与传输选项

新增枚举类型：

```go
type Protocol string

const (
    ProtocolBinary  Protocol = "binary"
    ProtocolCompact Protocol = "compact"
)

type Transport string

const (
    TransportBuffered Transport = "buffered"
    TransportFramed   Transport = "framed"
)
```

新增 option：

- `WithServerProtocol(protocol Protocol)`
- `WithServerTransport(transport Transport)`
- `WithServerSocketTimeout(timeout time.Duration)`
- `WithClientProtocol(protocol Protocol)`
- `WithClientTransport(transport Transport)`
- `WithClientConnectTimeout(timeout time.Duration)`
- `WithClientSocketTimeout(timeout time.Duration)`

非法协议或传输类型在构造阶段返回明确错误。

## ProtoRpcService 便捷层

新增 `ProtoServer`：

```go
func NewProtoServer(addr string, handler protoservice.ProtoRpcService, opts ...ServerOption) (*Server, error)
```

该函数内部创建 `protoservice.NewProtoRpcServiceProcessor(handler)`，再交给通用 `NewServer`。

新增 `ProtoClient`：

```go
type ProtoClient struct {
    // 内部持有 *ClientConn 和 *protoservice.ProtoRpcServiceClient
}

func NewProtoClient(addr string, opts ...ClientOption) (*ProtoClient, error)
func (c *ProtoClient) DealTwowayMessage(ctx context.Context, req *protoservice.ProtoRequest) (*protoservice.ProtoReply, error)
func (c *ProtoClient) DealOnewayMessage(ctx context.Context, req *protoservice.ProtoRequest) error
func (c *ProtoClient) Close() error
```

便捷层不吞掉 Apache Thrift 返回的错误，调用失败时将原始错误返回给使用方。

## 错误处理

- processor 或 handler 为 nil 时，构造函数返回 error。
- 地址为空时，构造函数返回 error。
- transport 打开失败时，客户端构造函数返回 error。
- 不支持的 protocol 或 transport 返回 error。
- `Close()` 可重复调用；重复关闭不应 panic。

## 测试计划

- 写失败测试验证通用 server/client 能完成一次 `DealTwowayMessage` 调用。
- 写失败测试验证 `ProtoClient.DealTwowayMessage` 能拿到 handler 返回内容。
- 写失败测试验证 `ProtoClient.DealOnewayMessage` 能触达 handler。
- 写失败测试验证客户端连接不存在的地址时返回 error。
- 写失败测试验证 `Close()` 重复调用安全。

## 文档

在 `docs/usage.md` 和 `README.md` 中补充 `glthrift` 简短用法，展示服务端启动和客户端调用的最小示例。
