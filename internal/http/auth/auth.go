package auth

import (
	"github.com/Verce11o/yata/internal/domain"
	"github.com/Verce11o/yata/internal/lib/response"
	"github.com/Verce11o/yata/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type Handler struct {
	log       *zap.SugaredLogger
	tracer    trace.Tracer
	service   service.Auth
	validator *validator.Validate
}

func NewHandler(log *zap.SugaredLogger, tracer trace.Tracer, service service.Auth, validator *validator.Validate) *Handler {
	return &Handler{log: log, tracer: tracer, service: service, validator: validator}
}

func (h *Handler) SignUp(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "Gateway.Signup")
	defer span.End()

	var input domain.SignUpInput

	if err := response.ReadRequest(c, h.validator, &input); err != nil {
		h.log.Errorf("Signup:HTTP: %s", err.Error())
		return response.WithError(c, err)
	}

	userID, err := h.service.Register(ctx, input)

	if err != nil {
		h.log.Errorf("Signup:GRPC: %s", err.Error())
		st, _ := status.FromError(err)
		return response.WithGRPCError(c, st.Code())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"id": userID,
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

	token, err := h.service.Login(ctx, input)

	if err != nil {
		h.log.Errorf("Login:GRPC: %s", err.Error())
		st, _ := status.FromError(err)
		return response.WithGRPCError(c, st.Code())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"token": token,
	})

}

func (h *Handler) Verify(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "Gateway.Verify")
	defer span.End()

	userID := c.Locals("userID")
	h.log.Debug(userID)

	err := h.service.VerifyUser(ctx, userID.(string))

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

	code := c.Query("code")

	if code == "" {
		h.log.Errorf("invalid code")
		return response.WithError(c, response.ErrInvalidCode)
	}

	err := h.service.CheckVerify(ctx, code)

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

	user, err := h.service.GetUserByID(ctx, userID.(string))

	if err != nil {
		h.log.Errorf("GetUserByID:GRPC: %v", err.Error())
		st, _ := status.FromError(err)
		return response.WithGRPCError(c, st.Code())
	}

	return c.Status(http.StatusOK).JSON(user)
}

func (h *Handler) ForgotPassword(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "Gateway.ForgotPassword")
	defer span.End()

	userID := c.Locals("userID")

	err := h.service.ForgotPassword(ctx, userID.(string))

	if err != nil {
		h.log.Errorf("ForgotPassword:GRPC: %v", err.Error())
		st, _ := status.FromError(err)
		return response.WithGRPCError(c, st.Code())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "success",
	})

}

func (h *Handler) VerifyPassword(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "Gateway.VerifyPassword")
	defer span.End()

	code := c.Query("code")

	if code == "" {
		h.log.Errorf("invalid code")
		return response.WithError(c, response.ErrInvalidCode)
	}

	err := h.service.VerifyPassword(ctx, code)

	if err != nil {
		h.log.Errorf("VerifyPassword:GRPC: %v", err.Error())
		st, _ := status.FromError(err)
		return response.WithGRPCError(c, st.Code())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "success",
	})

}

func (h *Handler) ResetPassword(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "Gateway.ResetPassword")
	defer span.End()

	var input domain.ResetPasswordRequest

	if err := response.ReadRequest(c, h.validator, &input); err != nil {
		h.log.Errorf("ResetPassword:HTTP: %s", err.Error())
		return response.WithError(c, err)
	}

	code := c.Locals("code")

	if code == nil {
		h.log.Errorf("ResetPassword:HTTP: %s", response.ErrInvalidCode)
		return response.WithError(c, response.ErrInvalidCode)
	}

	userID := c.Locals("userID")

	err := h.service.ResetPassword(ctx, code.(string), userID.(string), input)

	if err != nil {
		h.log.Errorf("ResetPassword:GRPC: %v", err.Error())
		st, _ := status.FromError(err)

		if st.Code() == codes.InvalidArgument {
			return response.WithError(c, response.ErrPasswordMismatch)
		}

		return response.WithGRPCError(c, st.Code())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "success",
	})

}
