package glthrift

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	thrift "github.com/apache/thrift/lib/go/thrift"
	"github.com/z-nuo/GLTools/glthrift/gen-go/protoservice"
)

type testProtoHandler struct {
	mu      sync.Mutex
	oneways []*protoservice.ProtoRequest
}

func (h *testProtoHandler) DealTwowayMessage(ctx context.Context, msg *protoservice.ProtoRequest) (*protoservice.ProtoReply, error) {
	return &protoservice.ProtoReply{
		Result_: msg.ShardingID + 1,
		Type:    msg.Type,
		Content: append([]byte("reply:"), msg.Content...),
		Trace:   append([]byte(nil), msg.Trace...),
	}, nil
}

func (h *testProtoHandler) DealOnewayMessage(ctx context.Context, msg *protoservice.ProtoRequest) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.oneways = append(h.oneways, msg)
	return nil
}

func (h *testProtoHandler) onewayCount() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.oneways)
}

func freeAddr(t *testing.T) string {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen free addr: %v", err)
	}
	addr := listener.Addr().String()
	if err := listener.Close(); err != nil {
		t.Fatalf("close free addr listener: %v", err)
	}
	return addr
}

func startTestServer(t *testing.T, server *Server) {
	t.Helper()
	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Serve()
	}()
	t.Cleanup(func() {
		_ = server.Stop()
		select {
		case <-errCh:
		case <-time.After(time.Second):
			t.Fatalf("server did not stop")
		}
	})
	time.Sleep(50 * time.Millisecond)
}

func TestGenericServerAndClientConnCompleteTwowayCall(t *testing.T) {
	handler := &testProtoHandler{}
	processor := protoservice.NewProtoRpcServiceProcessor(handler)
	addr := freeAddr(t)

	server, err := NewServer(addr, processor)
	if err != nil {
		t.Fatalf("new server: %v", err)
	}
	startTestServer(t, server)

	conn, err := NewClientConn(addr, WithClientConnectTimeout(time.Second), WithClientSocketTimeout(time.Second))
	if err != nil {
		t.Fatalf("new client conn: %v", err)
	}
	t.Cleanup(func() {
		if err := conn.Close(); err != nil {
			t.Fatalf("close conn: %v", err)
		}
	})

	client := protoservice.NewProtoRpcServiceClientFactory(conn.Transport(), conn.ProtocolFactory())
	reply, err := client.DealTwowayMessage(context.Background(), &protoservice.ProtoRequest{
		Type:       7,
		Content:    []byte("body"),
		ShardingID: 41,
		Trace:      []byte("trace"),
	})
	if err != nil {
		t.Fatalf("twoway call: %v", err)
	}
	if reply.Result_ != 42 || reply.Type != 7 || string(reply.Content) != "reply:body" || string(reply.Trace) != "trace" {
		t.Fatalf("unexpected reply: %#v", reply)
	}
}

func TestNewClientConnFailsForUnavailableAddress(t *testing.T) {
	addr := freeAddr(t)

	conn, err := NewClientConn(addr, WithClientConnectTimeout(50*time.Millisecond))
	if err == nil {
		_ = conn.Close()
		t.Fatalf("expected connection error")
	}
}

func TestClientConnCloseIsIdempotent(t *testing.T) {
	handler := &testProtoHandler{}
	server, err := NewServer(freeAddr(t), protoservice.NewProtoRpcServiceProcessor(handler))
	if err != nil {
		t.Fatalf("new server: %v", err)
	}
	startTestServer(t, server)

	conn, err := NewClientConn(server.Addr())
	if err != nil {
		t.Fatalf("new client conn: %v", err)
	}
	if err := conn.Close(); err != nil {
		t.Fatalf("first close: %v", err)
	}
	if err := conn.Close(); err != nil {
		t.Fatalf("second close: %v", err)
	}
}

func TestInvalidGenericInputsReturnErrors(t *testing.T) {
	if _, err := NewServer("", thrift.NewTMultiplexedProcessor()); err == nil {
		t.Fatalf("expected empty server addr error")
	}
	if _, err := NewServer("127.0.0.1:0", nil); err == nil {
		t.Fatalf("expected nil processor error")
	}
	if _, err := NewClientConn(""); err == nil {
		t.Fatalf("expected empty client addr error")
	}
	if _, err := NewServer("127.0.0.1:0", thrift.NewTMultiplexedProcessor(), WithServerProtocol(Protocol("unknown"))); err == nil {
		t.Fatalf("expected unsupported protocol error")
	}
	if _, err := NewServer("127.0.0.1:0", thrift.NewTMultiplexedProcessor(), WithServerTransport(Transport("unknown"))); err == nil {
		t.Fatalf("expected unsupported transport error")
	}
}

func TestProtoClientDealTwowayMessage(t *testing.T) {
	handler := &testProtoHandler{}
	server, err := NewProtoServer(freeAddr(t), handler)
	if err != nil {
		t.Fatalf("new proto server: %v", err)
	}
	startTestServer(t, server)

	client, err := NewProtoClient(server.Addr(), WithClientConnectTimeout(time.Second), WithClientSocketTimeout(time.Second))
	if err != nil {
		t.Fatalf("new proto client: %v", err)
	}
	t.Cleanup(func() {
		if err := client.Close(); err != nil {
			t.Fatalf("close proto client: %v", err)
		}
	})

	reply, err := client.DealTwowayMessage(context.Background(), &protoservice.ProtoRequest{
		Type:       3,
		Content:    []byte("hello"),
		ShardingID: 9,
	})
	if err != nil {
		t.Fatalf("proto twoway: %v", err)
	}
	if reply.Result_ != 10 || reply.Type != 3 || string(reply.Content) != "reply:hello" {
		t.Fatalf("unexpected proto reply: %#v", reply)
	}
}

func TestProtoClientDealOnewayMessage(t *testing.T) {
	handler := &testProtoHandler{}
	server, err := NewProtoServer(freeAddr(t), handler)
	if err != nil {
		t.Fatalf("new proto server: %v", err)
	}
	startTestServer(t, server)

	client, err := NewProtoClient(server.Addr())
	if err != nil {
		t.Fatalf("new proto client: %v", err)
	}
	t.Cleanup(func() {
		if err := client.Close(); err != nil {
			t.Fatalf("close proto client: %v", err)
		}
	})

	if err := client.DealOnewayMessage(context.Background(), &protoservice.ProtoRequest{
		Type:       5,
		Content:    []byte("event"),
		ShardingID: 11,
	}); err != nil {
		t.Fatalf("proto oneway: %v", err)
	}

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if handler.onewayCount() == 1 {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("oneway handler was not called")
}

func TestNewProtoServerRejectsNilHandler(t *testing.T) {
	if _, err := NewProtoServer("127.0.0.1:0", nil); err == nil {
		t.Fatalf("expected nil handler error")
	}
}
