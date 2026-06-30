package glthrift

import (
	"context"
	"testing"
	"time"

	"github.com/z-nuo/GLTools/glthrift/gen-go/protoservice"
)

func TestRPCClientCallsProtoTwowayAndOneway(t *testing.T) {
	handler := &testRPCHandler{}
	server, err := NewProtoServer(freeRPCAddr(t), handler)
	if err != nil {
		t.Fatalf("new proto server: %v", err)
	}
	startRPCServer(t, server)

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
		Type:       5,
		Content:    []byte("hello"),
		ShardingID: 7,
	})
	if err != nil {
		t.Fatalf("deal twoway message: %v", err)
	}
	if reply.Result_ != 107 || reply.Type != 5 || string(reply.Content) != "server:hello" {
		t.Fatalf("unexpected reply: %#v", reply)
	}

	if err := client.DealOnewayMessage(context.Background(), &protoservice.ProtoRequest{
		Type:       6,
		Content:    []byte("event"),
		ShardingID: 8,
	}); err != nil {
		t.Fatalf("deal oneway message: %v", err)
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
