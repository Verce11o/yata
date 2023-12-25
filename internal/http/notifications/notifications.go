package notifications

import (
	pb "github.com/Verce11o/yata-protos/gen/go/notifications"
	"github.com/Verce11o/yata/internal/lib/response"
	"github.com/Verce11o/yata/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
	"net/http"
)

type Handler struct {
	log       *zap.SugaredLogger
	tracer    trace.Tracer
	services  *service.Services
	validator *validator.Validate
}

func NewHandler(log *zap.SugaredLogger, tracer trace.Tracer, services *service.Services, validator *validator.Validate) *Handler {
	return &Handler{log: log, tracer: tracer, services: services, validator: validator}
}

func (h *Handler) SubscribeToUser(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "Gateway.SubscribeToUser")
	defer span.End()

	userID := c.Locals("userID")
	toUserID := c.Params("id")

	_, err := h.services.Notifications.SubscribeToUser(ctx, &pb.SubscribeToUserRequest{
		UserId:   userID.(string),
		ToUserId: toUserID,
	})

	if err != nil {
		h.log.Errorf("SubscribeToUser:GRPC: %v", err.Error())
		st, _ := status.FromError(err)
		return response.WithGRPCError(c, st.Code())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}

func (h *Handler) UnSubscribeFromUser(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "Gateway.UnSubscribeFromUser")
	defer span.End()

	userID := c.Locals("userID")
	toUserID := c.Params("id")

	_, err := h.services.Notifications.UnSubscribeFromUser(ctx, &pb.UnSubscribeFromUserRequest{
		UserId:   userID.(string),
		ToUserId: toUserID,
	})

	if err != nil {
		h.log.Errorf("UnSubscribeFromUser:GRPC: %v", err.Error())
		st, _ := status.FromError(err)
		return response.WithGRPCError(c, st.Code())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}
