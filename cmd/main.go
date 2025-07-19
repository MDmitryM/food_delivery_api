package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	server "github.com/MDmitryM/food_delivery_api/src/pb"
	"github.com/MDmitryM/food_delivery_api/src/pb/api"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		logrus.Fatalf("listening error %v", err)
	}

	_ = godotenv.Load() //убрать в продакшене

	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		logrus.Fatal("AUTH_SERVICE_URL environment variable is required")
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(server.AuthInterceptor(authServiceURL)),
	)
	api.RegisterGatewayServiceServer(s, &server.Server{})

	go func() {
		if err := s.Serve(lis); err != nil {
			logrus.Fatalf("failed to serve: %v", err)
		}
	}()

	logrus.Println("Server started!")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	logrus.Println("Stopping")
	s.GracefulStop()
}
