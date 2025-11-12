package proto

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// RegisterAllHandlersFromEndpoint registers all service handlers from a gRPC endpoint
func RegisterAllHandlersFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := RegisterHouseholdServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		return err
	}
	if err := RegisterMemberServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		return err
	}
	if err := RegisterMealServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		return err
	}
	if err := RegisterEventServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		return err
	}
	return nil
}

// RegisterAllHandlers registers all service handlers from a gRPC client connection
func RegisterAllHandlers(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	if err := RegisterHouseholdServiceHandler(ctx, mux, conn); err != nil {
		return err
	}
	if err := RegisterMemberServiceHandler(ctx, mux, conn); err != nil {
		return err
	}
	if err := RegisterMealServiceHandler(ctx, mux, conn); err != nil {
		return err
	}
	if err := RegisterEventServiceHandler(ctx, mux, conn); err != nil {
		return err
	}
	return nil
}

// RegisterAllHandlersServer registers all service handlers for direct server calls
func RegisterAllHandlersServer(ctx context.Context, mux *runtime.ServeMux,
	householdServer HouseholdServiceServer,
	memberServer MemberServiceServer,
	mealServer MealServiceServer,
	eventServer EventServiceServer) error {

	if err := RegisterHouseholdServiceHandlerServer(ctx, mux, householdServer); err != nil {
		return err
	}
	if err := RegisterMemberServiceHandlerServer(ctx, mux, memberServer); err != nil {
		return err
	}
	if err := RegisterMealServiceHandlerServer(ctx, mux, mealServer); err != nil {
		return err
	}
	if err := RegisterEventServiceHandlerServer(ctx, mux, eventServer); err != nil {
		return err
	}
	return nil
}
