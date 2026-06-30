package glthrift

import (
	"errors"

	"github.com/apache/thrift/lib/go/thrift"
)

// ClientConn 封装 Apache Thrift 客户端连接。
type ClientConn struct {
	addr            string
	transport       thrift.TTransport
	protocolFactory thrift.TProtocolFactory
}

// NewClientConn 创建并打开通用 Thrift 客户端连接。
func NewClientConn(addr string, opts ...ClientOption) (*ClientConn, error) {
	if addr == "" {
		return nil, errors.New("thrift client addr is empty")
	}

	cfg := ClientConfig{
		Addr:      addr,
		Protocol:  ProtocolBinary,
		Transport: TransportBuffered,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}

	protocolFactory, err := protocolFactory(cfg.Protocol)
	if err != nil {
		return nil, err
	}
	transportFactory, err := transportFactory(cfg.Transport)
	if err != nil {
		return nil, err
	}
	socket := thrift.NewTSocketConf(cfg.Addr, &thrift.TConfiguration{
		ConnectTimeout: cfg.ConnectTimeout,
		SocketTimeout:  cfg.SocketTimeout,
	})

	transport, err := transportFactory.GetTransport(socket)
	if err != nil {
		return nil, err
	}
	if err := transport.Open(); err != nil {
		return nil, err
	}

	return &ClientConn{
		addr:            cfg.Addr,
		transport:       transport,
		protocolFactory: protocolFactory,
	}, nil
}

// Transport 返回已打开的 Thrift transport。
func (c *ClientConn) Transport() thrift.TTransport {
	if c == nil {
		return nil
	}
	return c.transport
}

// ProtocolFactory 返回当前连接使用的 Thrift protocol factory。
func (c *ClientConn) ProtocolFactory() thrift.TProtocolFactory {
	if c == nil {
		return nil
	}
	return c.protocolFactory
}

// Addr 返回客户端连接地址。
func (c *ClientConn) Addr() string {
	if c == nil {
		return ""
	}
	return c.addr
}

// Close 关闭客户端连接，允许重复调用。
func (c *ClientConn) Close() error {
	if c == nil || c.transport == nil || !c.transport.IsOpen() {
		return nil
	}
	return c.transport.Close()
}
