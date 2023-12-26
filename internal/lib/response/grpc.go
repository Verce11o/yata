package response

import (
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/codes"
	"net/http"
)

func mapGRPCErrCodeToHttpStatus(code codes.Code) (int, string) {
	switch code {
	case codes.Unauthenticated:
		return http.StatusUnauthorized, "invalid credentials"
	case codes.AlreadyExists:
		return http.StatusBadRequest, "already exists"
	case codes.NotFound:
		return http.StatusNotFound, "not found"
	case codes.Internal:
		return http.StatusInternalServerError, "server error"
	case codes.PermissionDenied:
		return http.StatusForbidden, "permission denied"
	case codes.Canceled:
		return http.StatusRequestTimeout, "request canceled"
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout, "deadline exceeded"
	case codes.InvalidArgument:
		return http.StatusBadRequest, "invalid request"
	}
	return http.StatusInternalServerError, "server error"
}

func WithGRPCError(c *fiber.Ctx, code codes.Code) error {
	status, message := mapGRPCErrCodeToHttpStatus(code)
	return c.Status(status).JSON(fiber.Map{
		"message": message,
	})
}
