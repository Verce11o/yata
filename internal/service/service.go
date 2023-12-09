package service

import (
	pbSSO "github.com/Verce11o/yata-protos/gen/go/sso"
	pbTweets "github.com/Verce11o/yata-protos/gen/go/tweets"
	"github.com/Verce11o/yata/internal/config"
)

type Services struct {
	Auth   pbSSO.AuthClient
	Tweets pbTweets.TweetsClient
}

func NewServices(cfg *config.Config) *Services {
	return &Services{Auth: MakeAuthServiceClient(cfg.Services), Tweets: MakeTweetsServiceClient(cfg.Services)}
}
