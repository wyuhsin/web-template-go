package server

import (
	"context"
	"fmt"

	"github.com/tx7do/kratos-transport/broker"
	"github.com/tx7do/kratos-transport/transport/rabbitmq"
	"github.com/wyuhsin/web-template-go/internal/conf"
	"github.com/wyuhsin/web-template-go/internal/service"

	"github.com/go-kratos/kratos/v2/log"
)

const (
	EXCHANGE    = "server.greeter.exchange"
	ROUTING_KEY = "server.greeter.routingkey"
)

func NewRabbitMQServer(
	c *conf.Server,
	logger log.Logger,
	greeter *service.GreeterService,
) *rabbitmq.Server {
	srv := rabbitmq.NewServer(
		rabbitmq.WithAddress([]string{c.Rabbitmq.Addr}),
		rabbitmq.WithExchange(EXCHANGE, true),
		rabbitmq.WithCodec("json"),
	)

	srv.RegisterSubscriber(
		context.Background(),
		ROUTING_KEY,
		func(ctx context.Context, evt broker.Event) error {
			switch t := evt.Message().Body.(type) {
			case any:
				return nil
			default:
				return fmt.Errorf("unsupported type: %T", t)
			}
		},
		func() broker.Any { return struct{}{} },
	)

	return srv
}
