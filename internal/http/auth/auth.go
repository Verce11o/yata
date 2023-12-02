package auth

import (
	"context"
	pbSSO "github.com/Verce11o/yata-protos/gen/go/sso"
	"github.com/Verce11o/yata/internal/domain"
	"github.com/Verce11o/yata/internal/lib/response"
	"github.com/Verce11o/yata/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
	"net/http"
)

type Handler struct {
	log       *zap.SugaredLogger
	services  *service.Services
	validator *validator.Validate
}

func NewHandler(log *zap.SugaredLogger, services *service.Services, validator *validator.Validate) *Handler {
	return &Handler{log: log, services: services, validator: validator}
}

func (h *Handler) SignUp(c *fiber.Ctx) error {

	var input domain.SignUpInput

	if err := response.ReadRequest(c, h.validator, &input); err != nil {
		h.log.Errorf("Signup:HTTP: %s", err.Error())
		return response.WithError(c, err)
	}

	resp, err := h.services.Auth.Register(context.Background(), &pbSSO.RegisterRequest{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	})

	if err != nil {
		h.log.Errorf("Signup:GRPC: %s", err.Error())

		st, _ := status.FromError(err)

		return response.WithGRPCError(c, st.Code())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"id": resp.GetUserId(),
	})
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var input domain.SignInInput

	if err := response.ReadRequest(c, h.validator, &input); err != nil {
		h.log.Errorf("Login:HTTP: %s", err.Error())
		return response.WithError(c, err)
	}

	resp, err := h.services.Auth.Login(context.Background(), &pbSSO.LoginRequest{
		Email:    input.Email,
		Password: input.Password,
	})

	if err != nil {
		h.log.Errorf("Login:GRPC: %s", err.Error())

		st, _ := status.FromError(err)

		return response.WithGRPCError(c, st.Code())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"token": resp.GetToken(),
	})

}
