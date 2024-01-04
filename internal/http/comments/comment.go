package comments

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

func NewHandler(log *zap.SugaredLogger, tracer trace.Tracer, services *service.Services, validator *validator.Validate) *Handler {
	return &Handler{log: log, tracer: tracer, services: services, validator: validator}
}

func (h *Handler) CreateComment(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "Gateway.CreateComment")
	defer span.End()

	userID := c.Locals("userID")

	tweetID := c.Params("id")

	text := c.FormValue("text")

	imageInput, err := c.FormFile("image")

	var image *domain.Image

	if err == nil {
		contentType, bytes, imageName, err := files.PrepareImage(imageInput)

		if err != nil {
			h.log.Debugf("CreateComment:HTTP: %v", err.Error())
			return response.WithError(c, err)
		}

		image = &domain.Image{
			Chunk:       bytes,
			ContentType: contentType,
			ImageName:   imageName,
		}

	}

	commentID, err := h.services.Comments.CreateComment(ctx, domain.CreateCommentRequest{
		UserID:  userID.(string),
		TweetID: tweetID,
		Text:    text,
		Image:   image,
	})

	if err != nil {
		h.log.Errorf("CreateComment:GRPC: %v", err.Error())
		st, _ := status.FromError(err)
		return response.WithGRPCError(c, st.Code())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"id": commentID,
	})
}

func (h *Handler) GetComment(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "Gateway.GetComment")
	defer span.End()

	commentID := c.Params("id")

	comment, err := h.services.Comments.GetComment(ctx, commentID)

	if err != nil {
		h.log.Errorf("GetComment:GRPC: %v", err.Error())
		st, _ := status.FromError(err)
		return response.WithGRPCError(c, st.Code())
	}

	return c.Status(http.StatusOK).JSON(comment)
}

func (h *Handler) GetAllTweetComments(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "Gateway.GetComment")
	defer span.End()

	tweetID := c.Params("id")

	h.log.Debugf("tweetID: %v", tweetID)

	cursor := c.Query("cursor")

	comments, cursor, err := h.services.Comments.GetAllTweetComments(ctx, cursor, tweetID)

	if err != nil {
		h.log.Errorf("GetAllTweetComments:GRPC: %v", err.Error())
		return response.WithError(c, err)
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"data":   comments,
		"cursor": cursor,
	})
}

func (h *Handler) UpdateComment(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "Gateway.UpdateComment")
	defer span.End()

	userID := c.Locals("userID")
	tweetID := c.Params("id")
	commentID := c.Params("comment_id")

	text := c.FormValue("text")

	imageInput, err := c.FormFile("image")

	var image *domain.Image

	if err == nil {
		contentType, bytes, imageName, err := files.PrepareImage(imageInput)

		if err != nil {
			h.log.Debugf("UpdateComment:HTTP: %v", err.Error())
			return response.WithError(c, err)
		}

		image = &domain.Image{
			Chunk:       bytes,
			ContentType: contentType,
			ImageName:   imageName,
		}

	}

	comment, err := h.services.Comments.UpdateComment(ctx, domain.UpdateCommentRequest{
		TweetID:   tweetID,
		UserID:    userID.(string),
		Text:      text,
		Image:     image,
		CommentID: commentID,
	})

	if err != nil {
		h.log.Errorf("UpdateComment:GRPC: %v", err.Error())
		st, _ := status.FromError(err)
		return response.WithGRPCError(c, st.Code())
	}

	return c.Status(http.StatusOK).JSON(comment)
}

func (h *Handler) DeleteComment(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "Gateway.DeleteComment")
	defer span.End()

	userID := c.Locals("userID")
	commentID := c.Params("comment_id")

	err := h.services.Comments.DeleteComment(ctx, userID.(string), commentID)

	if err != nil {
		h.log.Errorf("DeleteComment:GRPC: %v", err.Error())
		st, _ := status.FromError(err)
		return response.WithGRPCError(c, st.Code())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}
