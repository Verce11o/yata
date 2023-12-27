package websocket

import (
	"github.com/Verce11o/yata/internal/service"
	"github.com/gofiber/contrib/websocket"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"sync"
)

type WsClients map[string]*websocket.Conn

type Handler struct {
	log      *zap.SugaredLogger
	tracer   trace.Tracer
	services *service.Services
	Clients  WsClients
	mu       sync.Mutex
}

func NewHandler(log *zap.SugaredLogger, tracer trace.Tracer, services *service.Services) *Handler {
	return &Handler{log: log, tracer: tracer, services: services, Clients: make(WsClients)}
}

func (h *Handler) EstablishConnection(c *websocket.Conn) {
	userID := c.Locals("userID")

	h.mu.Lock()
	h.Clients[userID.(string)] = c
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.Clients, userID.(string))
		h.mu.Unlock()

		err := c.Close()
		if err != nil {
			h.log.Errorf("error closing ws connection: %v", err)
		}

	}()

	select {}
}
