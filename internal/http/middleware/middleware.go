package middleware

import (
	"context"
	"errors"
	"github.com/Verce11o/yata-protos/gen/go/sso"
	"github.com/Verce11o/yata/internal/config"
	"github.com/Verce11o/yata/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"strings"
)

type Handler struct {
	log      *zap.SugaredLogger
	services *service.Services
	cfg      *config.Config
}

func NewMiddlewareHandler(log *zap.SugaredLogger, services *service.Services, cfg *config.Config) *Handler {
	return &Handler{log: log, services: services, cfg: cfg}
}

func (h *Handler) AuthMiddleware(ctx *fiber.Ctx) error {
	header := ctx.Get("Authorization")

	if header == "" {
		h.log.Infof("AuthMiddleware: empty header")
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "empty authorization header",
		})
	}

	headerParts := strings.Fields(header)
	if len(headerParts) != 2 {
		h.log.Infof("AuthMiddleware: invalid header")
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid authorization header",
		})
	}

	userID, err := ParseToken(headerParts[1], h.cfg.App.JWT.Secret)

	h.log.Debugf("AuthMiddleware: PassedUserID: %v", userID)

	if errors.Is(err, jwt.ErrTokenExpired) {
		h.log.Errorf("AuthMiddleware: %v", err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "token expired",
		})
	}

	if err != nil {
		h.log.Errorf("AuthMiddleware: %v", err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "server error",
		})
	}

	user, err := h.services.Auth.GetUserByID(context.Background(), &sso.GetUserRequest{UserId: userID})

	h.log.Debugf("AuthMiddleware: GetUser: %v", user.GetUserId())

	if err != nil {
		h.log.Errorf("AuthMiddleware: %v", err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "server error",
		})
	}

	ctx.Locals("userID", userID)

	return ctx.Next()
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
