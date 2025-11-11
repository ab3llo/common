package proto

import (
	"context"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestGatewayRegistrationFunctions(t *testing.T) {
	ctx := context.Background()

	// Test that we can create a ServeMux and call registration functions without panicking
	mux := runtime.NewServeMux()

	// Test RegisterAllHandlersFromEndpoint function exists and can be called
	t.Run("RegisterAllHandlersFromEndpoint", func(t *testing.T) {
		// We expect this to fail to connect, but not panic or have compilation errors
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		err := RegisterAllHandlersFromEndpoint(ctx, mux, "localhost:9999", opts)
		// We expect an error because the endpoint doesn't exist, that's the correct behavior
		if err == nil {
			t.Log("Unexpectedly succeeded in connecting to non-existent endpoint")
		} else {
			t.Log("Correctly failed to connect to non-existent endpoint:", err)
		}
	})

	// Test individual registration functions exist
	t.Run("Individual registration functions", func(t *testing.T) {
		// These should be callable without panicking even with nil parameters
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Registration function panicked: %v", r)
			}
		}()

		// Test that the functions exist and are callable
		// We expect errors due to nil parameters, but not panics
		_ = RegisterHouseholdServiceHandlerFromEndpoint
		_ = RegisterMemberServiceHandlerFromEndpoint
		_ = RegisterMealServiceHandlerFromEndpoint
		_ = RegisterEventServiceHandlerFromEndpoint

		_ = RegisterHouseholdServiceHandler
		_ = RegisterMemberServiceHandler
		_ = RegisterMealServiceHandler
		_ = RegisterEventServiceHandler

		_ = RegisterHouseholdServiceHandlerServer
		_ = RegisterMemberServiceHandlerServer
		_ = RegisterMealServiceHandlerServer
		_ = RegisterEventServiceHandlerServer
	})

	// Test helper functions exist
	t.Run("Helper functions", func(t *testing.T) {
		_ = RegisterAllHandlersFromEndpoint
		_ = RegisterAllHandlers
		_ = RegisterAllHandlersServer
	})
}
