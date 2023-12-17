package tweets

import (
	pbTweets "github.com/Verce11o/yata-protos/gen/go/tweets"
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

	var image *pbTweets.Image

	if err == nil {
		contentType, bytes, imageName, err := files.PrepareImage(imageInput)

		if err != nil {
			h.log.Debugf("CreateTweet:HTTP: %v", err.Error())
			return response.WithError(c, err)
		}

		image = &pbTweets.Image{
			Chunk:       bytes,
			ContentType: contentType,
			Name:        imageName,
		}

	}

	resp, err := h.services.Tweets.CreateTweet(ctx, &pbTweets.CreateTweetRequest{
		UserId: userID.(string),
		Text:   text,
		Image:  image,
	})

	if err != nil {
		h.log.Errorf("CreateTweet:GRPC: %v", err.Error())
		st, _ := status.FromError(err)
		return response.WithGRPCError(c, st.Code())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"id": resp.GetTweetId(),
	})

}

func (h *Handler) GetTweet(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "Gateway.GetTweet")
	defer span.End()

	tweetID := c.Params("id")

	tweet, err := h.services.Tweets.GetTweet(ctx, &pbTweets.GetTweetRequest{TweetId: tweetID})

	if err != nil {
		h.log.Errorf("GetTweet:GRPC: %v", err.Error())
		st, _ := status.FromError(err)
		return response.WithGRPCError(c, st.Code())
	}

	return c.Status(http.StatusOK).JSON(domain.TweetResponse{
		TweetID: tweet.GetTweetId(),
		UserID:  tweet.GetUserId(),
		Text:    tweet.GetText(),
	})

}

func (h *Handler) GetAllTweets(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "Gateway.GetTweet")
	defer span.End()

	cursor := c.Query("cursor")

	resp, err := h.services.Tweets.GetAllTweets(ctx, &pbTweets.GetAllTweetsRequest{Cursor: cursor})

	if err != nil {
		h.log.Errorf("GetAllTweets:GRPC: %v", err.Error())
		return response.WithError(c, err)
	}

	tweets := resp.GetTweets()

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"data":   tweets,
		"cursor": resp.GetCursor(),
	})

}

func (h *Handler) UpdateTweet(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "Gateway.UpdateTweet")
	defer span.End()

	userID := c.Locals("userID")
	tweetID := c.Params("id")

	text := c.FormValue("text")

	imageInput, err := c.FormFile("image")

	var image *pbTweets.Image

	if err == nil {
		contentType, bytes, imageName, err := files.PrepareImage(imageInput)

		if err != nil {
			h.log.Debugf("UpdateTweet:HTTP: %v", err.Error())
			return response.WithError(c, err)
		}

		image = &pbTweets.Image{
			Chunk:       bytes,
			ContentType: contentType,
			Name:        imageName,
		}

	}

	tweet, err := h.services.Tweets.UpdateTweet(ctx, &pbTweets.UpdateTweetRequest{
		UserId:  userID.(string),
		Text:    text,
		Image:   image,
		TweetId: tweetID,
	})

	if err != nil {
		h.log.Errorf("UpdateTweet:GRPC: %v", err.Error())
		st, _ := status.FromError(err)
		return response.WithGRPCError(c, st.Code())
	}

	return c.Status(http.StatusOK).JSON(domain.TweetResponse{
		TweetID: tweet.GetTweetId(),
		Text:    tweet.GetText(),
	})

}

func (h *Handler) DeleteTweet(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "Gateway.DeleteTweet")
	defer span.End()

	userID := c.Locals("userID")
	tweetID := c.Params("id")

	_, err := h.services.Tweets.DeleteTweet(ctx, &pbTweets.DeleteTweetRequest{UserId: userID.(string), TweetId: tweetID})

	if err != nil {
		h.log.Errorf("DeleteTweet:GRPC: %v", err.Error())
		st, _ := status.FromError(err)
		return response.WithGRPCError(c, st.Code())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "success",
	})

}
