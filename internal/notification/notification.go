package notification

import (
	"errors"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
	"log"
	"time"
)

var (
	ErrorNoAuth      = errors.New("no such notification")
	ErrorWhileFetch  = errors.New("error while fetch")
	ErrorWhileCreate = errors.New("error while create")
)

type DbRepository struct {
	connect *gorm.DB
}

func (d *DbRepository) Create(notification Notification) (createNotification Notification, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			err = ErrorWhileCreate
		}
	}()
	d.connect.Create(&notification)
	return notification, nil
}

type NotificationRepository interface {
	Create(notification Notification) (createNotification Notification, err error)
}

type Notification struct {
	gorm.Model
	ID        uint      `json:"id" gorm:"primary_key"`
	Email     string    `json:"email"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func NewNotification(email string) (Notification, error) {
	validate := validator.New()
	notification := Notification{
		Email: email,
	}
	errs := validate.Struct(notification)
	if errs != nil {
		return Notification{}, errs
	}
	return notification, nil
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &DbRepository{db}
}
