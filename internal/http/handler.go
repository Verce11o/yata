package http

import "github.com/gofiber/fiber/v2"

type AuthHandler interface {
	SignUp(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
}

type Handlers struct {
	AuthHandler
}

func NewHandlers(authHandler AuthHandler) *Handlers {
	return &Handlers{AuthHandler: authHandler}
}

func (h *Handlers) InitRoutes(app *fiber.App) {
	api := app.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.Post("/signup", h.SignUp)
			auth.Post("/login", h.Login)
		}
	}
}
