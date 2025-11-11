package proto

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Example usage functions - these are examples of how to use the registration functions

// SetupGatewayFromEndpoint sets up the gRPC gateway to proxy to a gRPC server endpoint
func SetupGatewayFromEndpoint(grpcEndpoint string, httpPort string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Create a new gRPC gateway mux
	mux := runtime.NewServeMux()

	// Setup gRPC dial options
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Register all service handlers
	if err := RegisterAllHandlersFromEndpoint(ctx, mux, grpcEndpoint, opts); err != nil {
		return fmt.Errorf("failed to register gateway handlers: %w", err)
	}

	// Start HTTP server
	log.Printf("Starting HTTP server on port %s, proxying to gRPC server at %s", httpPort, grpcEndpoint)
	return http.ListenAndServe(":"+httpPort, mux)
}

// SetupGatewayWithConnection sets up the gRPC gateway using an existing gRPC connection
func SetupGatewayWithConnection(conn *grpc.ClientConn, httpPort string) error {
	ctx := context.Background()

	// Create a new gRPC gateway mux
	mux := runtime.NewServeMux()

	// Register all service handlers
	if err := RegisterAllHandlers(ctx, mux, conn); err != nil {
		return fmt.Errorf("failed to register gateway handlers: %w", err)
	}

	// Start HTTP server
	log.Printf("Starting HTTP server on port %s with existing gRPC connection", httpPort)
	return http.ListenAndServe(":"+httpPort, mux)
}

// SetupGatewayWithServers sets up the gRPC gateway for direct server calls (in-process)
func SetupGatewayWithServers(
	householdServer HouseholdServiceServer,
	memberServer MemberServiceServer,
	mealServer MealServiceServer,
	eventServer EventServiceServer,
	httpPort string) error {

	ctx := context.Background()

	// Create a new gRPC gateway mux
	mux := runtime.NewServeMux()

	// Register all service handlers for direct server calls
	if err := RegisterAllHandlersServer(ctx, mux, householdServer, memberServer, mealServer, eventServer); err != nil {
		return fmt.Errorf("failed to register server handlers: %w", err)
	}

	// Start HTTP server
	log.Printf("Starting HTTP server on port %s with direct server handlers", httpPort)
	return http.ListenAndServe(":"+httpPort, mux)
}
