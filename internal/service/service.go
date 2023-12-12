package service

import (
	pbSSO "github.com/Verce11o/yata-protos/gen/go/sso"
	pbTweets "github.com/Verce11o/yata-protos/gen/go/tweets"
	"github.com/Verce11o/yata/internal/config"
	trace "github.com/Verce11o/yata/internal/lib/metrics/tracer"
)

type Services struct {
	Auth   pbSSO.AuthClient
	Tweets pbTweets.TweetsClient
}

func NewServices(cfg *config.Config, tracer *trace.JaegerTracing) *Services {
	return &Services{Auth: MakeAuthServiceClient(cfg.Services, tracer), Tweets: MakeTweetsServiceClient(cfg.Services, tracer)}
}
