package middleware

import (
	"errors"
	"fmt"
	"github.com/Verce11o/yata-protos/gen/go/sso"
	"github.com/Verce11o/yata/internal/config"
	"github.com/Verce11o/yata/internal/domain"
	"github.com/Verce11o/yata/internal/lib/response"
	"github.com/Verce11o/yata/internal/lib/token"
	"github.com/Verce11o/yata/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

type Handler struct {
	log       *zap.SugaredLogger
	tracer    trace.Tracer
	services  *service.Services
	cfg       *config.Config
	validator *validator.Validate
}

func NewMiddlewareHandler(log *zap.SugaredLogger, trace trace.Tracer, services *service.Services, cfg *config.Config, validator *validator.Validate) *Handler {
	return &Handler{log: log, tracer: trace, services: services, cfg: cfg, validator: validator}
}

func (h *Handler) AuthMiddleware(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), fmt.Sprintf("%s: %s", c.Method(), c.OriginalURL()))
	c.SetUserContext(ctx)
	defer span.End()

	header := c.Get("Authorization")

	if header == "" {
		h.log.Infof("AuthMiddleware: empty header")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "empty authorization header",
		})
	}

	span.AddEvent("strings.fields")

	headerParts := strings.Fields(header)
	if len(headerParts) != 2 {
		h.log.Infof("AuthMiddleware: invalid header")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid authorization header",
		})
	}

	span.AddEvent("parseToken")

	userID, err := token.ParseToken(headerParts[1], h.cfg.App.JWT.Secret)

	if errors.Is(err, jwt.ErrTokenExpired) {
		h.log.Errorf("AuthMiddleware: %v", err.Error())
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "token expired",
		})
	}

	if err != nil {
		h.log.Errorf("AuthMiddleware: %v", err.Error())
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "server error",
		})
	}

	// I guess it's useless
	//span.AddEvent("call Auth Service")
	//
	//_, err = h.services.Auth.GetUserByID(ctx, &sso.GetUserRequest{UserId: userID})
	//
	//if err != nil {
	//	h.log.Errorf("AuthMiddleware: %v", err.Error())
	//	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
	//		"message": "server error",
	//	})
	//}

	c.Locals("userID", userID)

	span.AddEvent("next request")
	return c.Next()
}

// TODO find better solution (this middleware calls everytime when sending request to password reset. must call only once)

func (h *Handler) PasswordResetMiddleware(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "PasswordResetMiddleware")
	c.SetUserContext(ctx)
	defer span.End()

	var input domain.ResetPasswordRequestMiddleware

	if err := response.ReadRequest(c, h.validator, &input); err != nil {
		h.log.Errorf("ResetPassword:HTTP: %s", err.Error())
		return response.WithError(c, err)
	}

	_, err := h.services.Auth.VerifyPassword(ctx, &sso.VerifyPasswordRequest{Code: input.Code})

	if err != nil {
		h.log.Errorf("PasswordResetMiddleware:GRPC: %v", err.Error())
		st, _ := status.FromError(err)
		if st.Code() == codes.InvalidArgument {
			return response.WithError(c, response.ErrInvalidCode)
		}
		return response.WithGRPCError(c, st.Code())
	}

	c.Locals("code", input.Code)

	return c.Next()
}
