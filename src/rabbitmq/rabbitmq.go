package rabbitmq

import (
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type RabbitHandler struct {
	RabbitConnection *amqp091.Connection
	RabbitChannel    *amqp091.Channel
}

func NewRabbitHandler(uri string) (*RabbitHandler, error) {
	conn, err := NewRabbitConn(uri)
	if err != nil {
		logrus.Errorln("Failed to open rabbit conncetion")
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		logrus.Errorln("Failed to open rabbit cahnnel")
		return nil, err
	}

	return &RabbitHandler{
		conn, channel,
	}, nil
}

func NewRabbitConn(uri string) (*amqp091.Connection, error) {
	conn, err := amqp091.Dial(uri)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
