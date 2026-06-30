package glthrift

import (
	"errors"

	thrift "github.com/apache/thrift/lib/go/thrift"
)

// Server 封装 Apache Thrift 服务端实例。
type Server struct {
	addr            string
	server          thrift.TServer
	serverTransport *thrift.TServerSocket
}

// NewServer 创建通用 Thrift 服务端。
func NewServer(addr string, processor thrift.TProcessor, opts ...ServerOption) (*Server, error) {
	if addr == "" {
		return nil, errors.New("thrift server addr is empty")
	}
	if processor == nil {
		return nil, errors.New("thrift processor is nil")
	}

	cfg := ServerConfig{
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
	serverTransport, err := thrift.NewTServerSocketTimeout(cfg.Addr, cfg.SocketTimeout)
	if err != nil {
		return nil, err
	}

	return &Server{
		addr:            cfg.Addr,
		server:          thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory),
		serverTransport: serverTransport,
	}, nil
}

// Serve 启动 Thrift 服务并阻塞直到服务停止或启动失败。
func (s *Server) Serve() error {
	if s == nil || s.server == nil {
		return errors.New("thrift server is nil")
	}
	return s.server.Serve()
}

// Stop 停止 Thrift 服务。
func (s *Server) Stop() error {
	if s == nil || s.server == nil {
		return nil
	}
	return s.server.Stop()
}

// Addr 返回服务端监听地址。
func (s *Server) Addr() string {
	if s == nil {
		return ""
	}
	if s.serverTransport != nil && s.serverTransport.Addr() != nil {
		return s.serverTransport.Addr().String()
	}
	return s.addr
}
