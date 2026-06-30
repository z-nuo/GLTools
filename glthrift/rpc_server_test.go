package glthrift

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/z-nuo/GLTools/glthrift/gen-go/protoservice"
)

type testRPCHandler struct {
	mu      sync.Mutex
	oneways []*protoservice.ProtoRequest
}

func (h *testRPCHandler) DealTwowayMessage(ctx context.Context, req *protoservice.ProtoRequest) (*protoservice.ProtoReply, error) {
	return &protoservice.ProtoReply{
		Result_: req.ShardingID + 100,
		Type:    req.Type,
		Content: append([]byte("server:"), req.Content...),
		Trace:   append([]byte(nil), req.Trace...),
	}, nil
}

func (h *testRPCHandler) DealOnewayMessage(ctx context.Context, req *protoservice.ProtoRequest) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.oneways = append(h.oneways, req)
	return nil
}

func (h *testRPCHandler) onewayCount() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.oneways)
}

func freeRPCAddr(t *testing.T) string {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen free addr: %v", err)
	}
	addr := listener.Addr().String()
	if err := listener.Close(); err != nil {
		t.Fatalf("close listener: %v", err)
	}
	return addr
}

func startRPCServer(t *testing.T, server *Server) {
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

func TestRPCServerHandlesProtoTwowayRequest(t *testing.T) {
	handler := &testRPCHandler{}
	server, err := NewProtoServer(freeRPCAddr(t), handler)
	if err != nil {
		t.Fatalf("new proto server: %v", err)
	}
	startRPCServer(t, server)

	conn, err := NewClientConn(server.Addr(), WithClientConnectTimeout(time.Second), WithClientSocketTimeout(time.Second))
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
		Type:       9,
		Content:    []byte("ping"),
		ShardingID: 23,
		Trace:      []byte("trace-001"),
	})
	if err != nil {
		t.Fatalf("deal twoway message: %v", err)
	}
	if reply.Result_ != 123 || reply.Type != 9 || string(reply.Content) != "server:ping" || string(reply.Trace) != "trace-001" {
		t.Fatalf("unexpected reply: %#v", reply)
	}
}
