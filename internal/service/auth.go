package service

import (
	pbSSO "github.com/Verce11o/yata-protos/gen/go/sso"
	"github.com/Verce11o/yata/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func MakeAuthServiceClient(cfg config.Services) pbSSO.AuthClient {

	// TODO: add retry
	cc, err := grpc.Dial(cfg.Auth.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("error while connect to auth client: %s", err)
	}

	return pbSSO.NewAuthClient(cc)
}
