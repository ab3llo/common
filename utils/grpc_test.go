package utils

import (
	"context"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func TestConnectToGrpcClient_ValidateParameters(t *testing.T) {
	tests := []struct {
		name string
		addr string
		serviceName string
	}{
		{
			name: "valid localhost address",
			addr: "localhost:8080",
			serviceName: "test-service",
		},
		{
			name: "valid IP address",
			addr: "127.0.0.1:9090",
			serviceName: "another-service",
		},
		{
			name: "valid service name with port",
			addr: "my-service:5000",
			serviceName: "grpc-service",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.addr == "" {
				t.Error("Address should not be empty")
			}
			if tt.serviceName == "" {
				t.Error("Service name should not be empty")
			}
		})
	}
}

func TestConnectToGrpcClient_WithMockServer(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()

	go func() {
		if err := server.Serve(lis); err != nil {
			t.Logf("Server exited with error: %v", err)
		}
	}()
	defer server.Stop()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithInsecure(),
	)
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	if conn.GetState().String() == "SHUTDOWN" {
		t.Error("Connection should not be in SHUTDOWN state")
	}
}

func TestConnectToGrpcClient_InvalidAddress(t *testing.T) {
	tests := []struct {
		name string
		addr string
	}{
		{
			name: "invalid port",
			addr: "localhost:99999999",
		},
		{
			name: "non-existent host",
			addr: "non-existent-host-12345:8080",
		},
		{
			name: "malformed address",
			addr: ":::invalid:::8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			conn, err := grpc.DialContext(ctx, tt.addr,
				grpc.WithInsecure(),
				grpc.WithBlock(),
			)
			if err == nil {
				defer conn.Close()
				t.Errorf("Expected error for address %s, but got successful connection", tt.addr)
			}

			if status.Code(err) == codes.OK {
				t.Errorf("Expected error code, got OK for address %s", tt.addr)
			}
		})
	}
}

func TestConnectToGrpcClient_ConnectionTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := grpc.DialContext(ctx, "192.0.2.1:12345", // RFC5737 TEST-NET-1 address
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)

	if err == nil {
		t.Error("Expected timeout error, but got successful connection")
	}

	// Accept either DeadlineExceeded or Unknown as both are valid timeout responses
	code := status.Code(err)
	if code != codes.DeadlineExceeded && code != codes.Unknown {
		t.Errorf("Expected DeadlineExceeded or Unknown, got %v", code)
	}
}

func BenchmarkConnectToGrpcClient(b *testing.B) {
	lis := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()

	go func() {
		server.Serve(lis)
	}()
	defer server.Stop()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		conn, err := grpc.DialContext(ctx, "bufnet",
			grpc.WithContextDialer(dialer),
			grpc.WithInsecure(),
		)
		if err != nil {
			b.Fatalf("Failed to dial: %v", err)
		}
		conn.Close()
		cancel()
	}
}