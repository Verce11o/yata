package service

import (
	"context"
	pbTweets "github.com/Verce11o/yata-protos/gen/go/tweets"
	"github.com/Verce11o/yata/internal/domain"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type TweetService struct {
	log    *zap.SugaredLogger
	tracer trace.Tracer
	client pbTweets.TweetsClient
}

func NewTweetService(log *zap.SugaredLogger, tracer trace.Tracer, client pbTweets.TweetsClient) *TweetService {
	return &TweetService{log: log, tracer: tracer, client: client}
}

func (t *TweetService) CreateTweet(ctx context.Context, input domain.CreateTweetRequest) (string, error) {
	ctx, span := t.tracer.Start(ctx, "Service.CreateTweet")
	defer span.End()

	var pbImage *pbTweets.Image

	if input.Image != nil {
		pbImage = &pbTweets.Image{
			Chunk:       input.Image.Chunk,
			ContentType: input.Image.ContentType,
			Name:        input.Image.ImageName,
		}
	}

	resp, err := t.client.CreateTweet(ctx, &pbTweets.CreateTweetRequest{
		UserId: input.UserID,
		Text:   input.Text,
		Image:  pbImage,
	})

	if err != nil {
		t.log.Errorf("cannot create tweet: %v", err)
		return "", err
	}

	return resp.GetTweetId(), nil

}

func (t *TweetService) GetTweet(ctx context.Context, tweetID string) (domain.TweetResponse, error) {
	ctx, span := t.tracer.Start(ctx, "Service.GetTweet")
	defer span.End()

	resp, err := t.client.GetTweet(ctx, &pbTweets.GetTweetRequest{TweetId: tweetID})
	if err != nil {
		t.log.Errorf("cannot get tweet: %v", err)
		return domain.TweetResponse{}, err
	}

	return domain.TweetResponse{
		TweetID:   resp.GetTweetId(),
		UserID:    resp.GetUserId(),
		Text:      resp.GetText(),
		CreatedAt: resp.GetCreatedAt().AsTime(),
	}, nil
}

func (t *TweetService) GetAllTweets(ctx context.Context, cursor string) ([]domain.TweetResponse, string, error) {
	ctx, span := t.tracer.Start(ctx, "Service.GetAllTweets")
	defer span.End()

	resp, err := t.client.GetAllTweets(ctx, &pbTweets.GetAllTweetsRequest{Cursor: cursor})
	if err != nil {
		t.log.Errorf("cannot get tweet: %v", err)
		return nil, "", err
	}

	result := make([]domain.TweetResponse, 0, len(resp.GetTweets()))

	for _, tweet := range resp.GetTweets() {
		item := domain.TweetResponse{
			TweetID:   tweet.GetTweetId(),
			UserID:    tweet.GetUserId(),
			Text:      tweet.GetText(),
			CreatedAt: tweet.GetCreatedAt().AsTime(),
		}
		result = append(result, item)
	}

	return result, resp.GetCursor(), nil

}

func (t *TweetService) UpdateTweet(ctx context.Context, input domain.UpdateTweetRequest) (domain.TweetResponse, error) {
	ctx, span := t.tracer.Start(ctx, "Service.UpdateTweet")
	defer span.End()

	var pbImage *pbTweets.Image

	if input.Image != nil {
		pbImage = &pbTweets.Image{
			Chunk:       input.Image.Chunk,
			ContentType: input.Image.ContentType,
			Name:        input.Image.ImageName,
		}
	}

	resp, err := t.client.UpdateTweet(ctx, &pbTweets.UpdateTweetRequest{
		UserId:  input.UserID,
		TweetId: input.TweetID,
		Text:    input.Text,
		Image:   pbImage,
	})
	if err != nil {
		t.log.Errorf("cannot update tweet: %v", err)
		return domain.TweetResponse{}, err
	}

	return domain.TweetResponse{
		TweetID:   resp.GetTweetId(),
		UserID:    resp.GetUserId(),
		Text:      resp.GetText(),
		CreatedAt: resp.CreatedAt.AsTime(),
	}, nil
}

func (t *TweetService) DeleteTweet(ctx context.Context, userID, tweetID string) error {
	ctx, span := t.tracer.Start(ctx, "Service.DeleteTweet")
	defer span.End()

	_, err := t.client.DeleteTweet(ctx, &pbTweets.DeleteTweetRequest{
		UserId:  userID,
		TweetId: tweetID,
	})
	if err != nil {
		t.log.Errorf("cannot delete tweet: %v", err)
		return err
	}
	return nil
}
