package middleware

import (
	"errors"
	"fmt"
	"github.com/Verce11o/yata-protos/gen/go/sso"
	"github.com/Verce11o/yata/internal/config"
	"github.com/Verce11o/yata/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"strings"
)

type Handler struct {
	log      *zap.SugaredLogger
	tracer   trace.Tracer
	services *service.Services
	cfg      *config.Config
}

func NewMiddlewareHandler(log *zap.SugaredLogger, trace trace.Tracer, services *service.Services, cfg *config.Config) *Handler {
	return &Handler{log: log, tracer: trace, services: services, cfg: cfg}
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

	userID, err := ParseToken(headerParts[1], h.cfg.App.JWT.Secret)

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

	span.AddEvent("call Auth Service")

	_, err = h.services.Auth.GetUserByID(ctx, &sso.GetUserRequest{UserId: userID})

	if err != nil {
		h.log.Errorf("AuthMiddleware: %v", err.Error())
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "server error",
		})
	}

	c.Locals("userID", userID)

	span.AddEvent("next request")
	return c.Next()
}

type tokenClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}

func ParseToken(token string, secret string) (string, error) {

	parsedToken, err := jwt.ParseWithClaims(token, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := parsedToken.Claims.(*tokenClaims)
	if !ok {
		return "", err
	}

	return claims.UserID, nil
}
