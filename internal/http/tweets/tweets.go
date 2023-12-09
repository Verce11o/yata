package tweets

import (
	"context"
	pbTweets "github.com/Verce11o/yata-protos/gen/go/tweets"
	"github.com/Verce11o/yata/internal/domain"
	"github.com/Verce11o/yata/internal/lib/files"
	"github.com/Verce11o/yata/internal/lib/response"
	"github.com/Verce11o/yata/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
	"net/http"
)

type Handler struct {
	log       *zap.SugaredLogger
	services  *service.Services
	validator *validator.Validate
}

func NewHandler(log *zap.SugaredLogger, services *service.Services, validator *validator.Validate) *Handler {
	return &Handler{log: log, services: services, validator: validator}
}

func (h *Handler) CreateTweet(c *fiber.Ctx) error {

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

	resp, err := h.services.Tweets.CreateTweet(context.Background(), &pbTweets.CreateTweetRequest{
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
	tweetID := c.Params("id")

	tweet, err := h.services.Tweets.GetTweet(context.Background(), &pbTweets.GetTweetRequest{TweetId: tweetID})

	if err != nil {
		h.log.Errorf("GetTweet:GRPC: %v", err.Error())
		st, _ := status.FromError(err)
		return response.WithGRPCError(c, st.Code())
	}

	return c.Status(http.StatusOK).JSON(domain.TweetResponse{
		TweetID:    tweet.GetTweetId(),
		UserID:     tweet.GetUserId(),
		Text:       tweet.GetText(),
		ImageChunk: tweet.GetImage().GetChunk(),
	})

}

func (h *Handler) UpdateTweet(c *fiber.Ctx) error {
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

	tweet, err := h.services.Tweets.UpdateTweet(context.Background(), &pbTweets.UpdateTweetRequest{
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
		TweetID:    tweet.GetTweetId(),
		UserID:     tweet.GetUserId(),
		Text:       tweet.GetText(),
		ImageChunk: tweet.GetImage().GetChunk(),
	})

}

func (h *Handler) DeleteTweet(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	tweetID := c.Params("id")

	_, err := h.services.Tweets.DeleteTweet(context.Background(), &pbTweets.DeleteTweetRequest{UserId: userID.(string), TweetId: tweetID})

	if err != nil {
		h.log.Errorf("DeleteTweet:GRPC: %v", err.Error())
		st, _ := status.FromError(err)
		return response.WithGRPCError(c, st.Code())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "success",
	})

}
