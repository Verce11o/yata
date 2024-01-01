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
	//TODO implement me
	panic("implement me")
}

func (t *TweetService) GetAllTweets(ctx context.Context, cursor string) ([]domain.TweetResponse, string, error) {
	//TODO implement me
	panic("implement me")
}

func (t *TweetService) UpdateTweet(ctx context.Context, input domain.UpdateTweetRequest) (domain.TweetResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (t *TweetService) DeleteTweet(ctx context.Context, userID, tweetID string) error {
	//TODO implement me
	panic("implement me")
}
