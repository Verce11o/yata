package service

import (
	"context"
	"github.com/Verce11o/yata/internal/clients"
	"github.com/Verce11o/yata/internal/config"
	"github.com/Verce11o/yata/internal/domain"
	trace "github.com/Verce11o/yata/internal/lib/metrics/tracer"
	"go.uber.org/zap"
	"time"
)

type Auth interface {
	Register(ctx context.Context, input domain.SignUpInput) (string, error)
	VerifyUser(ctx context.Context, userID string) error
	CheckVerify(ctx context.Context, code string) error
	Login(ctx context.Context, input domain.SignInInput) (string, error)
	GetUserByID(ctx context.Context, userID string) (domain.GetUserResponse, error)
	ForgotPassword(ctx context.Context, userID string) error
	VerifyPassword(ctx context.Context, code string) error
	ResetPassword(ctx context.Context, code string, userID string, input domain.ResetPasswordRequest) error
}

type Tweet interface {
	CreateTweet(ctx context.Context, input domain.CreateTweetRequest) (string, error)
	GetTweet(ctx context.Context, tweetID string) (domain.TweetResponse, error)
	GetAllTweets(ctx context.Context, cursor string) ([]domain.TweetResponse, string, error)
	UpdateTweet(ctx context.Context, input domain.UpdateTweetRequest) (domain.TweetResponse, error)
	DeleteTweet(ctx context.Context, userID, tweetID string) error
}

type Comment interface {
	CreateComment(ctx context.Context, input domain.CreateCommentRequest) (string, error)
	GetComment(ctx context.Context, commentID string) (domain.CommentResponse, error)
	GetAllTweetComments(ctx context.Context, cursor string, tweetID string) ([]domain.CommentResponse, string, error)
	UpdateComment(ctx context.Context, input domain.UpdateCommentRequest) (domain.CommentResponse, error)
	DeleteComment(ctx context.Context, commentID, userID string) error
}

type Notification interface {
	SubscribeToUser(ctx context.Context, userID, toUserID string) error
	UnSubscribeFromUser(ctx context.Context, userID, toUserID string) error
	GetNotifications(ctx context.Context, userID string) ([]domain.Notification, error)
	MarkNotificationAsRead(ctx context.Context, userID, notificationID string) error
	ReadAllNotifications(ctx context.Context, userID string) error
}

type Services struct {
	Auth          Auth
	Tweets        Tweet
	Comments      Comment
	Notifications Notification
}

const (
	grpcRetriesCount = 5
	grpcTimeout      = 5 * time.Second
)

// TODO add ping on start

func NewServices(cfg config.Services, log *zap.SugaredLogger, tracer *trace.JaegerTracing) *Services {
	return &Services{
		Auth:          NewAuthService(log, tracer.Tracer, clients.MakeAuthServiceClient(cfg, tracer, grpcRetriesCount, grpcTimeout)),
		Tweets:        NewTweetService(log, tracer.Tracer, clients.MakeTweetsServiceClient(cfg, tracer, grpcRetriesCount, grpcTimeout)),
		Comments:      NewCommentService(log, tracer.Tracer, clients.MakeCommentsServiceClient(cfg, tracer, grpcRetriesCount, grpcTimeout)),
		Notifications: NewNotificationService(log, tracer.Tracer, clients.MakeNotificationsServiceClient(cfg, tracer, grpcRetriesCount, grpcTimeout)),
	}
}
