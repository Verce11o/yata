package service

import (
	pbTweets "github.com/Verce11o/yata-protos/gen/go/tweets"
	"github.com/Verce11o/yata/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func MakeTweetsServiceClient(cfg config.Services) pbTweets.TweetsClient {

	// TODO: add retry
	cc, err := grpc.Dial(cfg.Tweets.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("error while connect to tweets client: %s", err)
	}

	return pbTweets.NewTweetsClient(cc)
}
