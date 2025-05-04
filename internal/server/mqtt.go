package server

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-transport/transport/mqtt"

	"github.com/tx7do/kratos-transport/broker"
	"github.com/wyuhsin/web-template-go/internal/conf"
	"github.com/wyuhsin/web-template-go/internal/service"
)

// NewMQTTServer create a mqtt server.
func NewMQTTServer(
	c *conf.Server,
	logger log.Logger,
	svc *service.GreeterService,
) *mqtt.Server {
	ctx := context.Background()

	srv := mqtt.NewServer(
		mqtt.WithAddress([]string{c.Mqtt.Addr}),
		mqtt.WithCodec("json"),
	)

	_ = srv.RegisterSubscriber(ctx,
		"/hfp/v2/journey/ongoing/vp/bus/#",
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
