package http

import (
	authHandler "github.com/Verce11o/yata/internal/http/auth"
	commentsHandler "github.com/Verce11o/yata/internal/http/comments"
	middlewareHandler "github.com/Verce11o/yata/internal/http/middleware"
	notificationHandler "github.com/Verce11o/yata/internal/http/notifications"
	tweetHandler "github.com/Verce11o/yata/internal/http/tweets"
	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	auth          *authHandler.Handler
	tweets        *tweetHandler.Handler
	comments      *commentsHandler.Handler
	notifications *notificationHandler.Handler
	middleware    *middlewareHandler.Handler
}

func NewHandlers(auth *authHandler.Handler, tweets *tweetHandler.Handler, comments *commentsHandler.Handler, notifications *notificationHandler.Handler, middleware *middlewareHandler.Handler) *Handlers {
	return &Handlers{auth: auth, tweets: tweets, comments: comments, notifications: notifications, middleware: middleware}
}

func (h *Handlers) InitRoutes(app *fiber.App) {
	api := app.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.Post("/signup", h.auth.SignUp)
			auth.Post("/login", h.auth.Login)

		}

		user := api.Group("/user", h.middleware.AuthMiddleware)
		{
			user.Get("/getme", h.auth.GetUserByID)
			user.Post("/verify", h.auth.Verify)
			user.Get("/activate", h.auth.Activate)

			user.Post("/forgot-password", h.auth.ForgotPassword)
			user.Get("/verify-password", h.auth.VerifyPassword)
			user.Put("/reset-password", h.middleware.PasswordResetMiddleware, h.auth.ResetPassword)

			subscribe := user.Group("/:id")
			{
				subscribe.Post("/subscribe", h.notifications.SubscribeToUser)
				subscribe.Post("/unsubscribe", h.notifications.UnSubscribeFromUser)
			}

		}

		tweets := api.Group("/tweets", h.middleware.AuthMiddleware)
		{
			tweets.Post("/", h.tweets.CreateTweet)
			tweets.Get("/", h.tweets.GetAllTweets)
			tweets.Get("/:id", h.tweets.GetTweet)
			tweets.Put("/:id", h.tweets.UpdateTweet)
			tweets.Delete("/:id", h.tweets.DeleteTweet)

			comments := tweets.Group("/:id/comments")
			{
				comments.Get("/", h.comments.GetAllTweetComments)
				comments.Post("/", h.comments.CreateComment)
				comments.Put("/:comment_id", h.comments.UpdateComment)
				comments.Delete("/:comment_id", h.comments.DeleteComment)
			}

		}

		notifications := api.Group("/notifications", h.middleware.AuthMiddleware)
		{
			notifications.Get("/", h.notifications.GetNotifications)
			notifications.Post("/read-notification", h.notifications.MarkNotificationAsRead)
			notifications.Post("/read-all-notifications", h.notifications.ReadAllNotifications)
		}

	}
}
