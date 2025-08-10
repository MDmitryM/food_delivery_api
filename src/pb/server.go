package server

import (
	"context"
	"log"

	"github.com/MDmitryM/food_delivery_api/src/pb/api"
	"github.com/MDmitryM/food_delivery_api/src/rabbitmq"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type Server struct {
	api.UnimplementedGatewayServiceServer
	rabbit *rabbitmq.RabbitHandler
}

func (s *Server) CreateOrder(ctx context.Context, in *api.CreateOrderRequest) (*api.CreateOrderResponse, error) {
	logrus.Println("CreteOrder invoked")

	q, err := s.rabbit.RabbitChannel.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		logrus.Errorln("Failed to create queue")
		//return //TODO
	}

	body := "Hello World!"
	err = s.rabbit.RabbitChannel.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		logrus.Errorf("Failed to publish a message, error: %s", err.Error())
	}
	log.Printf(" [x] Sent %s\n", body)

	return &api.CreateOrderResponse{OrderID: "OrderID", Status: "Status"}, nil
}

func NewServer(rabbit *rabbitmq.RabbitHandler) *Server {
	return &Server{rabbit: rabbit}
}
