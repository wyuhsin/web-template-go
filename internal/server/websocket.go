package server

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-transport/transport/websocket"

	"github.com/wyuhsin/web-template-go/internal/conf"
	"github.com/wyuhsin/web-template-go/internal/service"
)

// NewWebsocketServer create a websocket server.
func NewWebsocketServer(c *conf.Server, _ log.Logger, svc *service.GreeterService) *websocket.Server {
	srv := websocket.NewServer(
		websocket.WithAddress(c.Ws.Addr),
		websocket.WithPath(c.Ws.Path),
		websocket.WithConnectHandle(svc.OnWebsocketConnect),
		websocket.WithCodec("json"),
	)

	svc.SetWebsocketServer(srv)

	srv.RegisterMessageHandler(1,
		func(sessionId websocket.SessionID, payload websocket.MessagePayload) error {
			switch t := payload.(type) {
			case any:
				return svc.OnChatMessage(sessionId, t)
			default:
				return fmt.Errorf("unsupported type: %T", t)
			}
		},
		func() websocket.Any { return struct{}{} },
	)

	return srv
}
