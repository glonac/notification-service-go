package main

import (
	"context"
	"fmt"
	"github.com/streadway/amqp"
	"go.uber.org/dig"
	"gorm.io/gorm"
	"log"
	"net"
	"notification-service/internal/config"
	mail "notification-service/internal/mail/server"
	"notification-service/internal/notification"
	"notification-service/internal/storage"
	"notification-service/internal/tracer"
)

func main() {
	cfg := config.MustLoad()
	connect := storage.Initialize(cfg.DataBaseConf)

	err := tracer.NewTracer("http://jaeger:14268/api/traces", "server")
	if err != nil {
		log.Fatal(err)
	}

	defer tracer.Tracer.Shutdown(context.Background())
	lis, err := net.Listen("tcp", cfg.GrpcConf.Host+":"+cfg.GrpcConf.Port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	containerNotification := BuildContainerNotificationModule(connect, cfg.SmtpConf)

	// Define RabbitMQ server URL.
	amqpServerURL := "amqp://admin:password@rabbitmq:5672/"

	// Create a new RabbitMQ connection.
	connectRabbitMQ, err := amqp.Dial(amqpServerURL)
	if err != nil {
		panic(err)
	}
	defer connectRabbitMQ.Close()

	// Opening a channel to our RabbitMQ instance over
	// the connection we have already established.
	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		panic(err)
	}
	defer channelRabbitMQ.Close()

	// Subscribing to QueueService1 for getting messages.
	messages, err := channelRabbitMQ.Consume(
		"event", // queue name
		"",      // consumer
		true,    // auto-ack
		false,   // exclusive
		false,   // no local
		false,   // no wait
		nil,     // arguments
	)
	if err != nil {
		log.Println(err)
	}

	// Build a welcome message.
	log.Println("Successfully connected to RabbitMQ")
	log.Println("Waiting for messages")

	// Make a channel to receive messages into infinite loop.
	forever := make(chan bool)

	err = containerNotification.Invoke(func(s notification.NotificationService) {
		go func() {
			for message := range messages {
				// For example, show received message in a console.

				s.Create(notification.Notification{Email: "check@mail.ru", Text: "msg"})
				log.Printf(" > Received message: %s\n", message.Body)
			}
		}()

	})
	if err != nil {
		fmt.Println(err)
	}

	log.Printf("server listening at %v\n", lis.Addr())
	<-forever
}

func BuildContainerNotificationModule(connectDb *gorm.DB, cnf *config.ConfigSmtp) *dig.Container {
	container := dig.New()
	err := container.Provide(func() notification.NotificationRepository {
		return notification.NewNotificationRepository(connectDb)
	})

	if err != nil {
		fmt.Println(err)
	}

	err = container.Provide(notification.NewNotificationService)

	if err != nil {
		fmt.Println(err)
	}

	err = container.Provide(func() mail.MailService {
		return mail.NewMailService(cnf)
	})

	if err != nil {
		fmt.Println(err)
	}

	return container

}
