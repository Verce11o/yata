package http

import (
	"github.com/gofiber/fiber/v2"
)

type AuthHandler interface {
	SignUp(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	Verify(ctx *fiber.Ctx) error
	Activate(ctx *fiber.Ctx) error
}

type TweetsHandler interface {
	CreateTweet(ctx *fiber.Ctx) error
	GetTweet(ctx *fiber.Ctx) error
	GetAllTweets(ctx *fiber.Ctx) error
	UpdateTweet(ctx *fiber.Ctx) error
	DeleteTweet(ctx *fiber.Ctx) error
}

type CommentsHandler interface {
	CreateComment(ctx *fiber.Ctx) error
	GetComment(ctx *fiber.Ctx) error
	GetAllTweetComments(ctx *fiber.Ctx) error
	UpdateComment(ctx *fiber.Ctx) error
	DeleteComment(ctx *fiber.Ctx) error
}

type AuthMiddlewareHandler interface {
	AuthMiddleware(ctx *fiber.Ctx) error
}

type Handlers struct {
	AuthHandler
	TweetsHandler
	CommentsHandler
	AuthMiddlewareHandler
}

func NewHandlers(authHandler AuthHandler, tweetsHandler TweetsHandler, commentsHandler CommentsHandler, middlewareHandler AuthMiddlewareHandler) *Handlers {
	return &Handlers{AuthHandler: authHandler, TweetsHandler: tweetsHandler, CommentsHandler: commentsHandler, AuthMiddlewareHandler: middlewareHandler}
}

func (h *Handlers) InitRoutes(app *fiber.App) {
	api := app.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.Post("/signup", h.SignUp)
			auth.Post("/login", h.Login)

		}

		user := api.Group("/user")
		{
			user.Post("/verify", h.AuthMiddleware, h.Verify)
			user.Get("/activate", h.Activate)
		}

		tweets := api.Group("/tweets", h.AuthMiddleware)
		{
			tweets.Post("/", h.CreateTweet)
			tweets.Get("/", h.GetAllTweets)
			tweets.Get("/:id", h.GetTweet)
			tweets.Put("/:id", h.UpdateTweet)
			tweets.Delete("/:id", h.DeleteTweet)

			comments := tweets.Group("/:id/comments")
			{
				comments.Get("/", h.GetAllTweetComments)
				comments.Post("/", h.CreateComment)
				comments.Put("/:comment_id", h.UpdateComment)
				comments.Delete("/:comment_id", h.DeleteComment)
			}

		}

	}
}
