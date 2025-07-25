package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	server "github.com/MDmitryM/food_delivery_api/src/pb"
	"github.com/MDmitryM/food_delivery_api/src/pb/api"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":"+os.Getenv("API_PORT"))
	if err != nil {
		logrus.Fatalf("listening error %v", err)
	}

	authServiceURL := "http://" + os.Getenv("AUTH_HOST") + ":" + os.Getenv("AUTH_PORT")

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
