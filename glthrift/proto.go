package glthrift

import (
	"context"
	"errors"

	"github.com/z-nuo/GLTools/glthrift/gen-go/protoservice"
)

// NewProtoServer 创建当前 ProtoRpcService 的 Thrift 服务端。
func NewProtoServer(addr string, handler protoservice.ProtoRpcService, opts ...ServerOption) (*Server, error) {
	if handler == nil {
		return nil, errors.New("proto rpc service handler is nil")
	}
	processor := protoservice.NewProtoRpcServiceProcessor(handler)
	return NewServer(addr, processor, opts...)
}

// ProtoClient 封装当前 ProtoRpcService 的 Thrift 客户端。
type ProtoClient struct {
	conn   *ClientConn
	client *protoservice.ProtoRpcServiceClient
}

// NewProtoClient 创建并打开当前 ProtoRpcService 的 Thrift 客户端。
func NewProtoClient(addr string, opts ...ClientOption) (*ProtoClient, error) {
	conn, err := NewClientConn(addr, opts...)
	if err != nil {
		return nil, err
	}
	client := protoservice.NewProtoRpcServiceClientFactory(conn.Transport(), conn.ProtocolFactory())
	return &ProtoClient{
		conn:   conn,
		client: client,
	}, nil
}

// DealTwowayMessage 调用 ProtoRpcService 的双向消息接口。
func (c *ProtoClient) DealTwowayMessage(ctx context.Context, req *protoservice.ProtoRequest) (*protoservice.ProtoReply, error) {
	if c == nil || c.client == nil {
		return nil, errors.New("proto rpc service client is nil")
	}
	return c.client.DealTwowayMessage(ctx, req)
}

// DealOnewayMessage 调用 ProtoRpcService 的单向消息接口。
func (c *ProtoClient) DealOnewayMessage(ctx context.Context, req *protoservice.ProtoRequest) error {
	if c == nil || c.client == nil {
		return errors.New("proto rpc service client is nil")
	}
	return c.client.DealOnewayMessage(ctx, req)
}

// Close 关闭 ProtoRpcService 客户端连接。
func (c *ProtoClient) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}
