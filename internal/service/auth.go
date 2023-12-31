package service

import (
	"context"
	pbSSO "github.com/Verce11o/yata-protos/gen/go/sso"
	"github.com/Verce11o/yata/internal/domain"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type AuthService struct {
	log    *zap.SugaredLogger
	tracer trace.Tracer
	client pbSSO.AuthClient
}

func NewAuthService(log *zap.SugaredLogger, tracer trace.Tracer, client pbSSO.AuthClient) *AuthService {
	return &AuthService{log: log, tracer: tracer, client: client}
}

func (s *AuthService) Register(ctx context.Context, input domain.SignUpInput) (string, error) {
	ctx, span := s.tracer.Start(ctx, "Service.Register")
	defer span.End()

	resp, err := s.client.Register(ctx, &pbSSO.RegisterRequest{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	})

	if err != nil {
		s.log.Errorf("cannot register user: %v", err)
		return "", err
	}

	return resp.GetUserId(), nil
}

func (s *AuthService) VerifyUser(ctx context.Context, userID string) error {
	ctx, span := s.tracer.Start(ctx, "Service.VerifyUser")
	defer span.End()

	_, err := s.client.VerifyUser(ctx, &pbSSO.VerifyRequest{
		UserId: userID,
	})

	if err != nil {
		s.log.Errorf("cannot verify user: %v", err)
		return err
	}

	return nil
}

func (s *AuthService) CheckVerify(ctx context.Context, code string) error {
	ctx, span := s.tracer.Start(ctx, "Service.CheckVerify")
	defer span.End()

	_, err := s.client.CheckVerify(ctx, &pbSSO.CheckVerifyRequest{
		Code: code,
	})

	if err != nil {
		s.log.Errorf("cannot check verify: %v", err)
		return err
	}
	return nil
}

func (s *AuthService) Login(ctx context.Context, input domain.SignInInput) (string, error) {
	ctx, span := s.tracer.Start(ctx, "Service.Login")
	defer span.End()

	resp, err := s.client.Login(ctx, &pbSSO.LoginRequest{
		Email:    input.Email,
		Password: input.Password,
	})

	if err != nil {
		s.log.Errorf("cannot login user: %v", err)
		return "", err
	}

	return resp.GetToken(), nil
}

func (s *AuthService) GetUserByID(ctx context.Context, userID string) (domain.GetUserResponse, error) {
	ctx, span := s.tracer.Start(ctx, "Service.GetUserByID")
	defer span.End()

	user, err := s.client.GetUserByID(ctx, &pbSSO.GetUserRequest{UserId: userID})

	if err != nil {
		s.log.Errorf("cannot get user by id: %v", err)
		return domain.GetUserResponse{}, err
	}

	return domain.GetUserResponse{
		UserID:     uuid.MustParse(user.GetUserId()),
		Username:   user.GetUsername(),
		Email:      user.GetEmail(),
		IsVerified: user.GetIsVerified(),
		CreatedAt:  user.GetCreatedAt().AsTime(),
	}, nil

}

func (s *AuthService) ForgotPassword(ctx context.Context, userID string) error {
	ctx, span := s.tracer.Start(ctx, "Service.ForgotPassword")
	defer span.End()

	_, err := s.client.ForgotPassword(ctx, &pbSSO.ForgotPasswordRequest{UserId: userID})

	if err != nil {
		s.log.Errorf("cannot send forgot password request: %v", err)
		return err
	}

	return nil
}

func (s *AuthService) VerifyPassword(ctx context.Context, code string) error {
	ctx, span := s.tracer.Start(ctx, "Service.VerifyPassword")
	defer span.End()

	_, err := s.client.VerifyPassword(ctx, &pbSSO.VerifyPasswordRequest{Code: code})

	if err != nil {
		s.log.Errorf("cannot verify password: %v", err)
		return err
	}

	return nil
}

func (s *AuthService) ResetPassword(ctx context.Context, code string, userID string, input domain.ResetPasswordRequest) error {
	ctx, span := s.tracer.Start(ctx, "Service.ResetPassword")
	defer span.End()

	_, err := s.client.ResetPassword(ctx, &pbSSO.ResetPasswordRequest{
		Code:       code,
		UserId:     userID,
		Password:   input.Password,
		PasswordRe: input.PasswordRe,
	})

	if err != nil {
		s.log.Errorf("cannot reset password: %v", err)
		return err
	}

	return nil
}
