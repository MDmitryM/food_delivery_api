package main

import (
	"net"

	server "github.com/MDmitryM/food_delivery_api/src/pb"
	"github.com/MDmitryM/food_delivery_api/src/pb/api"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		logrus.Fatalf("listening error %v", err)
	}

	s := grpc.NewServer()
	api.RegisterGatewayServiceServer(s, &server.Server{})
	if err := s.Serve(lis); err != nil {
		logrus.Fatalf("failed to serve: %v", err)
	}

	logrus.Println("Server started!")
}
