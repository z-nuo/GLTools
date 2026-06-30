package glthrift

import (
	"fmt"
	"time"

	thrift "github.com/apache/thrift/lib/go/thrift"
)

const defaultBufferSize = 8192

// Protocol 表示 Apache Thrift 使用的协议类型。
type Protocol string

const (
	// ProtocolBinary 表示使用 Thrift binary protocol。
	ProtocolBinary Protocol = "binary"
	// ProtocolCompact 表示使用 Thrift compact protocol。
	ProtocolCompact Protocol = "compact"
)

// Transport 表示 Apache Thrift 使用的传输封装类型。
type Transport string

const (
	// TransportBuffered 表示使用 buffered transport。
	TransportBuffered Transport = "buffered"
	// TransportFramed 表示使用 framed transport。
	TransportFramed Transport = "framed"
)

// ServerConfig 表示 Thrift 服务端配置。
type ServerConfig struct {
	// Addr 是服务端监听地址，例如 "127.0.0.1:9090"。
	Addr string
	// Protocol 是服务端协议类型，空值默认使用 ProtocolBinary。
	Protocol Protocol
	// Transport 是服务端传输类型，空值默认使用 TransportBuffered。
	Transport Transport
	// SocketTimeout 是服务端接受的客户端连接读写超时时间。
	SocketTimeout time.Duration
}

// ClientConfig 表示 Thrift 客户端连接配置。
type ClientConfig struct {
	// Addr 是服务端连接地址，例如 "127.0.0.1:9090"。
	Addr string
	// Protocol 是客户端协议类型，空值默认使用 ProtocolBinary。
	Protocol Protocol
	// Transport 是客户端传输类型，空值默认使用 TransportBuffered。
	Transport Transport
	// ConnectTimeout 是客户端建立 TCP 连接的超时时间。
	ConnectTimeout time.Duration
	// SocketTimeout 是客户端连接读写超时时间。
	SocketTimeout time.Duration
}

// ServerOption 用于调整 Thrift 服务端配置。
type ServerOption func(*ServerConfig)

// ClientOption 用于调整 Thrift 客户端配置。
type ClientOption func(*ClientConfig)

// WithServerProtocol 设置服务端 Thrift 协议类型。
func WithServerProtocol(protocol Protocol) ServerOption {
	return func(cfg *ServerConfig) {
		cfg.Protocol = protocol
	}
}

// WithServerTransport 设置服务端 Thrift 传输类型。
func WithServerTransport(transport Transport) ServerOption {
	return func(cfg *ServerConfig) {
		cfg.Transport = transport
	}
}

// WithServerSocketTimeout 设置服务端接受连接后的读写超时时间。
func WithServerSocketTimeout(timeout time.Duration) ServerOption {
	return func(cfg *ServerConfig) {
		cfg.SocketTimeout = timeout
	}
}

// WithClientProtocol 设置客户端 Thrift 协议类型。
func WithClientProtocol(protocol Protocol) ClientOption {
	return func(cfg *ClientConfig) {
		cfg.Protocol = protocol
	}
}

// WithClientTransport 设置客户端 Thrift 传输类型。
func WithClientTransport(transport Transport) ClientOption {
	return func(cfg *ClientConfig) {
		cfg.Transport = transport
	}
}

// WithClientConnectTimeout 设置客户端建立连接的超时时间。
func WithClientConnectTimeout(timeout time.Duration) ClientOption {
	return func(cfg *ClientConfig) {
		cfg.ConnectTimeout = timeout
	}
}

// WithClientSocketTimeout 设置客户端连接读写超时时间。
func WithClientSocketTimeout(timeout time.Duration) ClientOption {
	return func(cfg *ClientConfig) {
		cfg.SocketTimeout = timeout
	}
}

func protocolFactory(protocol Protocol) (thrift.TProtocolFactory, error) {
	if protocol == "" {
		protocol = ProtocolBinary
	}
	switch protocol {
	case ProtocolBinary:
		return thrift.NewTBinaryProtocolFactoryDefault(), nil
	case ProtocolCompact:
		return thrift.NewTCompactProtocolFactory(), nil
	default:
		return nil, fmt.Errorf("unsupported thrift protocol: %s", protocol)
	}
}

func transportFactory(transport Transport) (thrift.TTransportFactory, error) {
	if transport == "" {
		transport = TransportBuffered
	}
	switch transport {
	case TransportBuffered:
		return thrift.NewTBufferedTransportFactory(defaultBufferSize), nil
	case TransportFramed:
		return thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory()), nil
	default:
		return nil, fmt.Errorf("unsupported thrift transport: %s", transport)
	}
}
