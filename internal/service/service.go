package service

import (
	pbComments "github.com/Verce11o/yata-protos/gen/go/comments"
	pbNotifications "github.com/Verce11o/yata-protos/gen/go/notifications"
	pbSSO "github.com/Verce11o/yata-protos/gen/go/sso"
	pbTweets "github.com/Verce11o/yata-protos/gen/go/tweets"
	"github.com/Verce11o/yata/internal/config"
	trace "github.com/Verce11o/yata/internal/lib/metrics/tracer"
	"time"
)

type Services struct {
	Auth          pbSSO.AuthClient
	Tweets        pbTweets.TweetsClient
	Comments      pbComments.CommentsClient
	Notifications pbNotifications.NotificationsClient
}

const (
	grpcRetriesCount = 5
	grpcTimeout      = 5 * time.Second
)

// TODO add ping on start

func NewServices(cfg *config.Config, tracer *trace.JaegerTracing) *Services {
	return &Services{
		Auth:          MakeAuthServiceClient(cfg.Services, tracer, grpcRetriesCount, grpcTimeout),
		Tweets:        MakeTweetsServiceClient(cfg.Services, tracer, grpcRetriesCount, grpcTimeout),
		Comments:      MakeCommentsServiceClient(cfg.Services, tracer, grpcRetriesCount, grpcTimeout),
		Notifications: MakeNotificationsServiceClient(cfg.Services, tracer, grpcRetriesCount, grpcTimeout),
	}
}
