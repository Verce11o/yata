package tweets

import (
	"github.com/Verce11o/yata/internal/domain"
	"github.com/Verce11o/yata/internal/lib/files"
	"github.com/Verce11o/yata/internal/lib/response"
	"github.com/Verce11o/yata/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
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

func NewHandler(log *zap.SugaredLogger, trace trace.Tracer, services *service.Services, validator *validator.Validate) *Handler {
	return &Handler{log: log, tracer: trace, services: services, validator: validator}
}

func (h *Handler) CreateTweet(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "Gateway.CreateTweet")
	defer span.End()

	userID := c.Locals("userID")

	text := c.FormValue("text")

	imageInput, err := c.FormFile("image")

	var image *domain.Image

	if err == nil {
		contentType, bytes, imageName, err := files.PrepareImage(imageInput)

		if err != nil {
			h.log.Debugf("CreateTweet:HTTP: %v", err.Error())
			return response.WithError(c, err)
		}

		image = &domain.Image{
			Chunk:       bytes,
			ContentType: contentType,
			ImageName:   imageName,
		}

	}

	tweetID, err := h.services.Tweets.CreateTweet(ctx, domain.CreateTweetRequest{
		UserID: userID.(string),
		Text:   text,
		Image:  image,
	})

	if err != nil {
		h.log.Errorf("CreateTweet:GRPC: %v", err.Error())
		st, _ := status.FromError(err)
		return response.WithGRPCError(c, st.Code())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"id": tweetID,
	})

}

func (h *Handler) GetTweet(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "Gateway.GetTweet")
	defer span.End()

	tweetID := c.Params("id")

	tweet, err := h.services.Tweets.GetTweet(ctx, tweetID)

	if err != nil {
		h.log.Errorf("GetTweet:GRPC: %v", err.Error())
		st, _ := status.FromError(err)
		return response.WithGRPCError(c, st.Code())
	}

	return c.Status(http.StatusOK).JSON(tweet)

}

func (h *Handler) GetAllTweets(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "Gateway.GetAllTweets")
	defer span.End()

	cursor := c.Query("cursor")

	tweets, cursor, err := h.services.Tweets.GetAllTweets(ctx, cursor)

	if err != nil {
		h.log.Errorf("GetAllTweets:GRPC: %v", err.Error())
		return response.WithError(c, err)
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"data":   tweets,
		"cursor": cursor,
	})

}

func (h *Handler) UpdateTweet(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "Gateway.UpdateTweet")
	defer span.End()

	userID := c.Locals("userID")
	tweetID := c.Params("id")

	text := c.FormValue("text")

	imageInput, err := c.FormFile("image")

	var image *domain.Image

	if err == nil {
		contentType, bytes, imageName, err := files.PrepareImage(imageInput)

		if err != nil {
			h.log.Debugf("UpdateTweet:HTTP: %v", err.Error())
			return response.WithError(c, err)
		}

		image = &domain.Image{
			Chunk:       bytes,
			ContentType: contentType,
			ImageName:   imageName,
		}

	}

	tweet, err := h.services.Tweets.UpdateTweet(ctx, domain.UpdateTweetRequest{
		UserID:  userID.(string),
		Text:    text,
		Image:   image,
		TweetID: tweetID,
	})

	if err != nil {
		h.log.Errorf("UpdateTweet:GRPC: %v", err.Error())
		st, _ := status.FromError(err)
		return response.WithGRPCError(c, st.Code())
	}

	return c.Status(http.StatusOK).JSON(domain.TweetResponse{
		TweetID: tweet.TweetID,
		Text:    tweet.Text,
	})

}

func (h *Handler) DeleteTweet(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "Gateway.DeleteTweet")
	defer span.End()

	userID := c.Locals("userID")
	tweetID := c.Params("id")

	err := h.services.Tweets.DeleteTweet(ctx, userID.(string), tweetID)

	if err != nil {
		h.log.Errorf("DeleteTweet:GRPC: %v", err.Error())
		st, _ := status.FromError(err)
		return response.WithGRPCError(c, st.Code())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "success",
	})

}
