package utils

import (
	"os"

	"github.com/hmlylab/common/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	log = logger.NewLogger()
)

func ConnectToGrpcClient(name, addr string) *grpc.ClientConn {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(
		insecure.NewCredentials(),
	))
	if err != nil {
		log.Error("Failed to connect to gRPC server: " + name + " " + err.Error())
		os.Exit(1) // Exit if we can't connect to the gRPC server
	}
	return conn
}
