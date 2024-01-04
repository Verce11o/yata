package service

import (
	"context"
	pbComments "github.com/Verce11o/yata-protos/gen/go/comments"
	"github.com/Verce11o/yata/internal/domain"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type CommentService struct {
	log    *zap.SugaredLogger
	tracer trace.Tracer
	client pbComments.CommentsClient
}

func NewCommentService(log *zap.SugaredLogger, tracer trace.Tracer, client pbComments.CommentsClient) *CommentService {
	return &CommentService{log: log, tracer: tracer, client: client}
}

func (c *CommentService) CreateComment(ctx context.Context, input domain.CreateCommentRequest) (string, error) {
	ctx, span := c.tracer.Start(ctx, "Service.CreateComment")
	defer span.End()

	var pbImage *pbComments.Image

	if input.Image != nil {
		pbImage = &pbComments.Image{
			Chunk:       input.Image.Chunk,
			ContentType: input.Image.ContentType,
			Name:        input.Image.ImageName,
		}
	}

	resp, err := c.client.CreateComment(ctx, &pbComments.CreateCommentRequest{
		UserId:  input.UserID,
		TweetId: input.TweetID,
		Text:    input.Text,
		Image:   pbImage,
	})

	if err != nil {
		c.log.Errorf("cannot create comment: %v", err)
		return "", err
	}

	return resp.GetCommentId(), nil
}

func (c *CommentService) GetComment(ctx context.Context, commentID string) (domain.CommentResponse, error) {
	ctx, span := c.tracer.Start(ctx, "Service.GetComment")
	defer span.End()

	resp, err := c.client.GetComment(ctx, &pbComments.GetCommentRequest{CommentId: commentID})
	if err != nil {
		c.log.Errorf("cannot get comment: %v", err)
		return domain.CommentResponse{}, err
	}
	return domain.CommentResponse{
		CommentID: resp.GetCommentId(),
		UserID:    resp.GetUserId(),
		TweetID:   resp.GetTweetId(),
		Text:      resp.GetText(),
		CreatedAt: resp.GetCreatedAt().AsTime(),
	}, nil
}

func (c *CommentService) GetAllTweetComments(ctx context.Context, cursor string, tweetID string) ([]domain.CommentResponse, string, error) {
	ctx, span := c.tracer.Start(ctx, "Service.GetAllTweetComments")
	defer span.End()

	resp, err := c.client.GetAllTweetComments(ctx, &pbComments.GetAllTweetCommentsRequest{
		Cursor:  cursor,
		TweetId: tweetID,
	})

	if err != nil {
		c.log.Errorf("cannot get all tweet comments: %v", err)
		return nil, "", err
	}

	result := make([]domain.CommentResponse, 0, len(resp.GetComments()))

	for _, comment := range resp.GetComments() {
		item := domain.CommentResponse{
			CommentID: comment.GetCommentId(),
			TweetID:   comment.GetTweetId(),
			UserID:    comment.GetUserId(),
			Text:      comment.GetText(),
			CreatedAt: comment.GetCreatedAt().AsTime(),
		}
		result = append(result, item)
	}

	return result, resp.GetCursor(), nil
}

func (c *CommentService) UpdateComment(ctx context.Context, input domain.UpdateCommentRequest) (domain.CommentResponse, error) {
	ctx, span := c.tracer.Start(ctx, "Service.UpdateComment")
	defer span.End()

	var pbImage *pbComments.Image

	if input.Image != nil {
		pbImage = &pbComments.Image{
			Chunk:       input.Image.Chunk,
			ContentType: input.Image.ContentType,
			Name:        input.Image.ImageName,
		}
	}

	resp, err := c.client.UpdateComment(ctx, &pbComments.UpdateCommentRequest{
		CommentId: input.CommentID,
		UserId:    input.UserID,
		TweetId:   input.TweetID,
		Text:      input.Text,
		Image:     pbImage,
	})

	if err != nil {
		c.log.Errorf("cannot update comment: %v", err)
		return domain.CommentResponse{}, err
	}

	return domain.CommentResponse{
		CommentID: resp.GetCommentId(),
		TweetID:   resp.GetTweetId(),
		UserID:    resp.GetUserId(),
		Text:      resp.GetText(),
		CreatedAt: resp.CreatedAt.AsTime(),
	}, nil
}

func (c *CommentService) DeleteComment(ctx context.Context, commentID, userID string) error {
	ctx, span := c.tracer.Start(ctx, "Service.DeleteComment")
	defer span.End()

	_, err := c.client.DeleteComment(ctx, &pbComments.DeleteCommentRequest{
		CommentId: commentID,
		UserId:    userID,
	})

	if err != nil {
		c.log.Errorf("cannot delete comment: %v", err)
		return err
	}

	return nil
}
