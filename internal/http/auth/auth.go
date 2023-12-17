package auth

import (
	pbSSO "github.com/Verce11o/yata-protos/gen/go/sso"
	"github.com/Verce11o/yata/internal/domain"
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

func (h *Handler) SignUp(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "Gateway.Signup")
	defer span.End()

	var input domain.SignUpInput

	if err := response.ReadRequest(c, h.validator, &input); err != nil {
		h.log.Errorf("Signup:HTTP: %s", err.Error())
		return response.WithError(c, err)
	}

	resp, err := h.services.Auth.Register(ctx, &pbSSO.RegisterRequest{
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
	ctx, span := h.tracer.Start(c.UserContext(), "Gateway.Login")
	defer span.End()

	var input domain.SignInInput

	if err := response.ReadRequest(c, h.validator, &input); err != nil {
		h.log.Errorf("Login:HTTP: %s", err.Error())
		return response.WithError(c, err)
	}

	resp, err := h.services.Auth.Login(ctx, &pbSSO.LoginRequest{
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

func (h *Handler) Verify(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "Gateway.Verify")
	defer span.End()

	userID := c.Locals("userID")
	h.log.Debug(userID)

	_, err := h.services.Auth.VerifyUser(ctx, &pbSSO.VerifyRequest{
		UserId: userID.(string),
	})

	if err != nil {
		h.log.Errorf("Verify:GRPC: %v", err.Error())
		st, _ := status.FromError(err)
		return response.WithGRPCError(c, st.Code())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "success",
	})

}

func (h *Handler) Activate(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "Gateway.Activate")
	defer span.End()

	//userID := c.Locals("userID")

	code := c.Query("code")

	if code == "" {
		h.log.Errorf("invalid code")
		return response.WithError(c, response.ErrInvalidCode)
	}

	_, err := h.services.Auth.CheckVerify(ctx, &pbSSO.CheckVerifyRequest{
		Code: code,
	})

	if err != nil {
		h.log.Errorf("Activate:GRPC: %v", err.Error())
		st, _ := status.FromError(err)
		return response.WithGRPCError(c, st.Code())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "success",
	})

}

func (h *Handler) GetUserByID(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "Gateway.GetUserByID")
	defer span.End()

	userID := c.Locals("userID")

	user, err := h.services.Auth.GetUserByID(ctx, &pbSSO.GetUserRequest{UserId: userID.(string)})

	if err != nil {
		h.log.Errorf("GetUserByID:GRPC: %v", err.Error())
		st, _ := status.FromError(err)
		return response.WithGRPCError(c, st.Code())
	}

	return c.Status(http.StatusOK).JSON(domain.GetUserResponse{
		UserID:     uuid.MustParse(user.GetUserId()),
		Username:   user.GetUsername(),
		Email:      user.GetEmail(),
		IsVerified: user.GetIsVerified(),
		CreatedAt:  user.GetCreatedAt().AsTime(),
	})
}
