package service

import (
	"context"
	"errors"
	"pet-project/internal/realtime"
	"pet-project/internal/repository"
	"pet-project/pkg/model"
	"time"
)

type NotificationService struct {
	Repository    repository.NotificationRepository
	ClientManager *realtime.ClientManager
}

func (s *NotificationService) Create(ctx context.Context, notif *model.Notification) error {
	if notif.Message == "" {
		return errors.New("message cannot be empty")
	}

	if notif.UserID <= 0 {
		return errors.New("invalid user ID")
	}

	notif.CreatedAt = time.Now()
	notif.IsRead = false

	err := s.Repository.Create(ctx, notif)
	if err != nil {
		return err
	}

	// Отправляем уведомление через WebSocket если клиент подключен
	if s.ClientManager != nil {
		s.ClientManager.Send(notif.UserID, *notif)
	}

	return nil
}

func (s *NotificationService) GetByUserID(ctx context.Context, userID int, limit, offset int) ([]model.Notification, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}

	if limit <= 0 {
		limit = 10 // Дефолтный лимит
	}
	if offset < 0 {
		offset = 0
	}

	notifications, err := s.Repository.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (s *NotificationService) MarkAsRead(ctx context.Context, userID int, notifIDs []int) error {
	if userID <= 0 {
		return errors.New("invalid user ID")
	}

	if len(notifIDs) == 0 {
		return errors.New("no notification IDs provided")
	}

	err := s.Repository.MarkAsRead(ctx, userID, notifIDs)
	if err != nil {
		return err
	}

	return nil
}

func (s *NotificationService) CountUnread(ctx context.Context, userID int) (int, error) {
	if userID <= 0 {
		return 0, errors.New("invalid user ID")
	}

	count, err := s.Repository.CountUnread(ctx, userID)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CreateNotificationForUser создает уведомление для конкретного пользователя
func (s *NotificationService) CreateNotificationForUser(ctx context.Context, userID int, notifType, message string) error {
	notif := &model.Notification{
		UserID:  userID,
		Type:    notifType,
		Message: message,
	}
	return s.Create(ctx, notif)
}

// CreateNotificationForMultipleUsers создает уведомления для нескольких пользователей
func (s *NotificationService) CreateNotificationForMultipleUsers(ctx context.Context, userIDs []int, notifType, message string) error {
	for _, userID := range userIDs {
		notif := &model.Notification{
			UserID:  userID,
			Type:    notifType,
			Message: message,
		}
		if err := s.Create(ctx, notif); err != nil {
			return err
		}
	}
	return nil
}
