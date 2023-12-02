package service

import (
	pbSSO "github.com/Verce11o/yata-protos/gen/go/sso"
	"github.com/Verce11o/yata/internal/config"
)

type Services struct {
	Auth pbSSO.AuthClient
}

func NewServices(cfg *config.Config) *Services {
	return &Services{Auth: MakeAuthServiceClient(cfg.Services)}
}
