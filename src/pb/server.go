package server

import (
	"context"

	"github.com/MDmitryM/food_delivery_api/src/pb/api"
	"github.com/sirupsen/logrus"
)

type Server struct {
	api.UnimplementedGatewayServiceServer
}

func (s *Server) CreateOrder(ctx context.Context, in *api.CreateOrderRequest) (*api.CreateOrderResponse, error) {
	logrus.Println("CreteOrder invoked")
	return &api.CreateOrderResponse{OrderID: "OrderID", Status: "Status"}, nil
}

func NewServer() *Server {
	return &Server{}
}
