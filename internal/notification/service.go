package notification

import (
	"fmt"
	mail "notification-service/internal/mail/server"
)

type notificationService struct {
	repo NotificationRepository
	smtp mail.MailService
}

type NotificationService interface {
	Create(notification Notification) error
}

// TODO add test for this and test this
func (s *notificationService) Create(notification Notification) error {
	newNotification, err := s.repo.Create(notification)
	if err != nil {
		return err
	}

	err = s.smtp.SendMail([]string{newNotification.Email}, newNotification.Text)

	if err != nil {
		return fmt.Errorf("Error while send", err)
	}
	return nil
}

func NewNotificationService(repo NotificationRepository, smtp mail.MailService) NotificationService {
	return &notificationService{repo: repo, smtp: smtp}
}
