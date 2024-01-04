package service

import (
	"context"
	pbNotifications "github.com/Verce11o/yata-protos/gen/go/notifications"
	"github.com/Verce11o/yata/internal/domain"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type NotificationService struct {
	log    *zap.SugaredLogger
	tracer trace.Tracer
	client pbNotifications.NotificationsClient
}

func NewNotificationService(log *zap.SugaredLogger, tracer trace.Tracer, client pbNotifications.NotificationsClient) *NotificationService {
	return &NotificationService{log: log, tracer: tracer, client: client}
}

func (n *NotificationService) SubscribeToUser(ctx context.Context, userID, toUserID string) error {
	ctx, span := n.tracer.Start(ctx, "Service.SubscribeUser")
	defer span.End()

	_, err := n.client.SubscribeToUser(ctx, &pbNotifications.SubscribeToUserRequest{
		UserId:   userID,
		ToUserId: toUserID,
	})

	if err != nil {
		n.log.Errorf("cannot subscribe to user: %v", err)
		return err
	}

	return nil
}

func (n *NotificationService) UnSubscribeFromUser(ctx context.Context, userID, toUserID string) error {
	ctx, span := n.tracer.Start(ctx, "Service.UnSubscribeFromUser")
	defer span.End()

	_, err := n.client.UnSubscribeFromUser(ctx, &pbNotifications.UnSubscribeFromUserRequest{
		UserId:   userID,
		ToUserId: toUserID,
	})

	if err != nil {
		n.log.Errorf("cannot unsubscribe from user: %v", err)
		return err
	}

	return nil
}

func (n *NotificationService) GetNotifications(ctx context.Context, userID string) ([]domain.Notification, error) {
	ctx, span := n.tracer.Start(ctx, "Service.GetNotifications")
	defer span.End()

	resp, err := n.client.GetNotifications(ctx, &pbNotifications.GetNotificationsRequest{UserId: userID})

	if err != nil {
		n.log.Errorf("cannot get user notifications: %v", err)
		return nil, err
	}

	result := make([]domain.Notification, 0, len(resp.GetNotifications()))

	for _, notif := range resp.GetNotifications() {
		item := domain.Notification{
			NotificationID: notif.GetNotificationId(),
			UserID:         notif.GetUserId(),
			SenderID:       notif.GetSenderId(),
			Read:           notif.GetRead(),
			CreatedAt:      notif.GetCreatedAt().AsTime(),
			Type:           notif.GetType(),
		}
		result = append(result, item)
	}

	return result, nil
}

func (n *NotificationService) MarkNotificationAsRead(ctx context.Context, userID, notificationID string) error {
	ctx, span := n.tracer.Start(ctx, "Service.MarkNotificationAsRead")
	defer span.End()

	_, err := n.client.MarkNotificationAsRead(ctx, &pbNotifications.MarkNotificationAsReadRequest{
		UserId:         userID,
		NotificationId: notificationID,
	})

	if err != nil {
		n.log.Errorf("cannot mark notification as read: %v", err)
		return err
	}

	return nil
}

func (n *NotificationService) ReadAllNotifications(ctx context.Context, userID string) error {
	ctx, span := n.tracer.Start(ctx, "Service.ReadAllNotifications")
	defer span.End()

	_, err := n.client.ReadAllNotifications(ctx, &pbNotifications.ReadAllNotificationsRequest{UserId: userID})
	if err != nil {
		n.log.Errorf("cannot read all notifications: %v", err)
		return err
	}

	return nil
}
