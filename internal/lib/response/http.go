package response

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

var (
	ErrInvalidRequest   = errors.New("invalid request")
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidImage     = errors.New("invalid image")
	ErrInvalidCode      = errors.New("invalid code or expired")
	ErrPasswordMismatch = errors.New("password mismatch")
)

func mapErrorWithCode(err error) int {
	switch {
	case errors.Is(err, ErrInvalidRequest):
		return http.StatusBadRequest
	case errors.Is(err, ErrInvalidImage):
		return http.StatusBadRequest
	case errors.Is(err, ErrPasswordMismatch):
		return http.StatusBadRequest
	case errors.Is(err, ErrInvalidCode):
		return http.StatusBadRequest
	case errors.Is(err, fiber.ErrUpgradeRequired):
		return http.StatusUpgradeRequired
	case errors.Is(err, ErrUserNotFound):
		return http.StatusNotFound
	}

	return http.StatusInternalServerError
}

// WithError responds to request with provided error
func WithError(c *fiber.Ctx, err error) error {

	return c.Status(mapErrorWithCode(err)).JSON(fiber.Map{
		"message": err.Error(),
	})
}

// ReadRequest parses and validates request
func ReadRequest(c *fiber.Ctx, v *validator.Validate, request any) error {
	if err := c.BodyParser(&request); err != nil {
		return ErrInvalidRequest
	}

	if err := v.Struct(request); err != nil {
		validateErr := err.(validator.ValidationErrors)

		return ValidationError{validateErr}
	}

	return nil

}
