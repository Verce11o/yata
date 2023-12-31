package notifications

import (
	pb "github.com/Verce11o/yata-protos/gen/go/notifications"
	"github.com/Verce11o/yata/internal/lib/response"
	"github.com/Verce11o/yata/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

	_, err := uuid.Parse(toUserID)

	if err != nil {
		h.log.Errorf("SubscribeToUser:HTTP: %v", err.Error())
		return response.WithError(c, response.ErrInvalidRequest)
	}

	_, err = h.services.Auth.GetUserByID(ctx, toUserID)

	if err != nil {
		h.log.Errorf("SubscribeToUser:HTTP: %v", err.Error())
		return response.WithError(c, response.ErrUserNotFound)
	}

	_, err = h.services.Notifications.SubscribeToUser(ctx, &pb.SubscribeToUserRequest{
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

	_, err := uuid.Parse(toUserID)

	if err != nil {
		h.log.Errorf("UnSubscribeFromUser:HTTP: %v", err.Error())
		return response.WithError(c, response.ErrInvalidRequest)
	}

	_, err = h.services.Notifications.UnSubscribeFromUser(ctx, &pb.UnSubscribeFromUserRequest{
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

func (h *Handler) GetNotifications(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "Gateway.GetNotifications")
	defer span.End()

	userID := c.Locals("userID")

	resp, err := h.services.Notifications.GetNotifications(ctx, &pb.GetNotificationsRequest{UserId: userID.(string)})

	if err != nil {
		h.log.Errorf("GetNotifications:GRPC: %v", err.Error())
		st, _ := status.FromError(err)
		return response.WithGRPCError(c, st.Code())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"data": resp.Notifications,
	})

}

func (h *Handler) MarkNotificationAsRead(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "Gateway.MarkNotificationAsRead")
	defer span.End()

	userID := c.Locals("userID")
	notificationID := c.Query("id")

	_, err := uuid.Parse(notificationID)

	if err != nil {
		h.log.Errorf("MarkNotificationAsRead:HTTP: %v", err.Error())
		return response.WithError(c, response.ErrInvalidRequest)
	}

	_, err = h.services.Notifications.MarkNotificationAsRead(ctx, &pb.MarkNotificationAsReadRequest{
		UserId:         userID.(string),
		NotificationId: notificationID,
	})

	if err != nil {
		h.log.Errorf("MarkNotificationAsRead:GRPC: %v", err.Error())
		st, _ := status.FromError(err)
		return response.WithGRPCError(c, st.Code())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "success",
	})

}

func (h *Handler) ReadAllNotifications(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "Gateway.ReadAllNotifications")
	defer span.End()

	userID := c.Locals("userID")

	_, err := h.services.Notifications.ReadAllNotifications(ctx, &pb.ReadAllNotificationsRequest{UserId: userID.(string)})

	if err != nil {
		h.log.Errorf("ReadAllNotifications:GRPC: %v", err.Error())
		st, _ := status.FromError(err)
		return response.WithGRPCError(c, st.Code())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "success",
	})

}
