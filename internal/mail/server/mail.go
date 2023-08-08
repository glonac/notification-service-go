package mail

import (
	"fmt"
	"net/smtp"
	"notification-service/internal/config"
)

type MailService interface {
	SendMail(to []string, msg string) error
}

type mailService struct {
	cnf config.ConfigSmtp
}

func NewMailService(cnf *config.ConfigSmtp) MailService {
	return &mailService{cnf: *cnf}
}

func (s *mailService) SendMail(to []string, msg string) error {
	message := []byte(msg)
	auth := smtp.PlainAuth("", s.cnf.From, s.cnf.Password, s.cnf.Host)
	err := smtp.SendMail(s.cnf.Host+":"+s.cnf.Port, auth, s.cnf.From, to, message)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
