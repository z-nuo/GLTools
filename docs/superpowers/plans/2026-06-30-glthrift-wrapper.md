# glthrift Wrapper Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build reusable Apache Thrift server/client primitives and a convenience wrapper for the generated `ProtoRpcService`.

**Architecture:** `glthrift` exposes a small generic layer around Apache Thrift protocol, transport, server, and client connection setup. `ProtoClient` and `NewProtoServer` sit on top of the generic layer and delegate request handling to `glthrift/gen-go/protoservice`.

**Tech Stack:** Go 1.25, `github.com/apache/thrift v0.23.0`, generated Go thrift package `github.com/z-nuo/GLTools/glthrift/gen-go/protoservice`.

## Global Constraints

- Do not regenerate thrift code.
- Do not modify `proto.thrift` service definitions.
- Do not add service discovery, load balancing, circuit breaking, retry, or connection pooling.
- Keep APIs framework-neutral and return errors instead of logging or panicking.
- All exported APIs added by this plan must have Chinese GoDoc comments.
- Use TDD: write each behavior test first, run it and verify the expected failure, then implement the minimal code.
- Use `/usr/local/go/bin/go` for local verification because `go` is not currently on this shell PATH.

---

## File Structure

- Create `glthrift/types.go`: protocol/transport enums, config structs, option types, validation helpers, and factory builders.
- Create `glthrift/server.go`: generic `Server`, `NewServer`, `Serve`, `Stop`, and `Addr`.
- Create `glthrift/client.go`: generic `ClientConn`, `NewClientConn`, `Transport`, `ProtocolFactory`, `Addr`, and `Close`.
- Create `glthrift/proto.go`: `NewProtoServer`, `ProtoClient`, `NewProtoClient`, RPC methods, and `Close`.
- Create `glthrift/thrift_test.go`: end-to-end tests using the generated `ProtoRpcService`.
- Modify `docs/usage.md`: add a compact `glthrift` usage section.
- Modify `README.md`: add `glthrift` to package list or usage summary using the repository's existing README style.

---

### Task 1: Generic Thrift Server And Client Connection

**Files:**
- Create: `glthrift/types.go`
- Create: `glthrift/server.go`
- Create: `glthrift/client.go`
- Create: `glthrift/thrift_test.go`

**Interfaces:**
- Produces:
  - `type Protocol string`
  - `const ProtocolBinary Protocol = "binary"`
  - `const ProtocolCompact Protocol = "compact"`
  - `type Transport string`
  - `const TransportBuffered Transport = "buffered"`
  - `const TransportFramed Transport = "framed"`
  - `type ServerOption func(*ServerConfig)`
  - `type ClientOption func(*ClientConfig)`
  - `func WithServerProtocol(protocol Protocol) ServerOption`
  - `func WithServerTransport(transport Transport) ServerOption`
  - `func WithServerSocketTimeout(timeout time.Duration) ServerOption`
  - `func WithClientProtocol(protocol Protocol) ClientOption`
  - `func WithClientTransport(transport Transport) ClientOption`
  - `func WithClientConnectTimeout(timeout time.Duration) ClientOption`
  - `func WithClientSocketTimeout(timeout time.Duration) ClientOption`
  - `func NewServer(addr string, processor thrift.TProcessor, opts ...ServerOption) (*Server, error)`
  - `func (s *Server) Serve() error`
  - `func (s *Server) Stop() error`
  - `func (s *Server) Addr() string`
  - `func NewClientConn(addr string, opts ...ClientOption) (*ClientConn, error)`
  - `func (c *ClientConn) Transport() thrift.TTransport`
  - `func (c *ClientConn) ProtocolFactory() thrift.TProtocolFactory`
  - `func (c *ClientConn) Addr() string`
  - `func (c *ClientConn) Close() error`

- [ ] **Step 1: Write failing generic integration tests**

Add this initial test scaffolding to `glthrift/thrift_test.go`:

```go
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
	mu       sync.Mutex
	oneways []*protoservice.ProtoRequest
}

func (h *testProtoHandler) DealTwowayMessage(ctx context.Context, msg *protoservice.ProtoRequest) (*protoservice.ProtoReply, error) {
	return &protoservice.ProtoReply{
		Result:  msg.ShardingID + 1,
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
	if reply.Result != 42 || reply.Type != 7 || string(reply.Content) != "reply:body" || string(reply.Trace) != "trace" {
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
```

- [ ] **Step 2: Run tests to verify they fail because APIs are missing**

Run: `/usr/local/go/bin/go test ./glthrift`

Expected: FAIL with compile errors mentioning undefined `Server`, `NewServer`, `NewClientConn`, and option functions.

- [ ] **Step 3: Implement generic types and factories**

Create `glthrift/types.go` with exported Chinese GoDoc comments and helper functions for protocol/transport factory creation. Use `thrift.NewTBinaryProtocolFactoryDefault()`, `thrift.NewTCompactProtocolFactory()`, `thrift.NewTBufferedTransportFactory(8192)`, and `thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())`.

- [ ] **Step 4: Implement generic server**

Create `glthrift/server.go`. Validate `addr` and `processor`; build `thrift.NewTServerSocketTimeout(addr, cfg.SocketTimeout)`; build factories; use `thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)`. Store the socket as `serverTransport` so `Addr()` can return the actual listener address when available, otherwise the configured address.

- [ ] **Step 5: Implement generic client connection**

Create `glthrift/client.go`. Validate `addr`; create socket with `thrift.NewTSocketTimeout(addr, cfg.ConnectTimeout, cfg.SocketTimeout)`; wrap it using the selected transport factory; call `Open()` on the wrapped transport; return `ClientConn`. `Close()` should return nil when called on nil or already closed connection.

- [ ] **Step 6: Run tests to verify generic layer passes**

Run: `/usr/local/go/bin/go test ./glthrift`

Expected: PASS for all tests currently in `glthrift`.

---

### Task 2: ProtoRpcService Convenience Layer

**Files:**
- Modify: `glthrift/thrift_test.go`
- Create: `glthrift/proto.go`

**Interfaces:**
- Consumes:
  - `func NewServer(addr string, processor thrift.TProcessor, opts ...ServerOption) (*Server, error)`
  - `func NewClientConn(addr string, opts ...ClientOption) (*ClientConn, error)`
- Produces:
  - `func NewProtoServer(addr string, handler protoservice.ProtoRpcService, opts ...ServerOption) (*Server, error)`
  - `type ProtoClient struct`
  - `func NewProtoClient(addr string, opts ...ClientOption) (*ProtoClient, error)`
  - `func (c *ProtoClient) DealTwowayMessage(ctx context.Context, req *protoservice.ProtoRequest) (*protoservice.ProtoReply, error)`
  - `func (c *ProtoClient) DealOnewayMessage(ctx context.Context, req *protoservice.ProtoRequest) error`
  - `func (c *ProtoClient) Close() error`

- [ ] **Step 1: Write failing proto convenience tests**

Append these tests to `glthrift/thrift_test.go`:

```go
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
	if reply.Result != 10 || reply.Type != 3 || string(reply.Content) != "reply:hello" {
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
```

- [ ] **Step 2: Run tests to verify they fail because proto APIs are missing**

Run: `/usr/local/go/bin/go test ./glthrift`

Expected: FAIL with compile errors mentioning undefined `NewProtoServer`, `NewProtoClient`, and `ProtoClient`.

- [ ] **Step 3: Implement proto convenience layer**

Create `glthrift/proto.go`. `NewProtoServer` validates nil handler and calls `protoservice.NewProtoRpcServiceProcessor(handler)`. `NewProtoClient` opens `ClientConn` and creates `protoservice.NewProtoRpcServiceClientFactory(conn.Transport(), conn.ProtocolFactory())`. RPC methods delegate directly to generated client. `Close()` delegates to `ClientConn.Close()`.

- [ ] **Step 4: Run tests to verify proto layer passes**

Run: `/usr/local/go/bin/go test ./glthrift`

Expected: PASS.

---

### Task 3: Usage Documentation And Package Verification

**Files:**
- Modify: `docs/usage.md`
- Modify: `README.md`

**Interfaces:**
- Consumes:
  - All public APIs produced by Tasks 1 and 2.

- [ ] **Step 1: Add usage docs**

Append a `## glthrift` section to `docs/usage.md` with imports for `glthrift` and `protoservice`, a minimal handler, server startup, client call, and `Close()` call.

- [ ] **Step 2: Add README summary**

Update `README.md` using the existing package summary style to mention `glthrift` as an Apache Thrift server/client wrapper with `ProtoRpcService` convenience APIs.

- [ ] **Step 3: Run full verification**

Run: `/usr/local/go/bin/go test ./...`

Expected: PASS for all packages.

- [ ] **Step 4: Inspect final diff**

Run: `git diff -- glthrift docs/usage.md README.md`

Expected: Diff contains only the generic thrift wrapper, proto convenience wrapper, tests, and docs updates.
