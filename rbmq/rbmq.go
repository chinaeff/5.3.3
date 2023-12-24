package rbmq

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/ptflp/gopubsub/queue"
	"go.uber.org/ratelimit"
	_ "go.uber.org/ratelimit"
	"log"
	"net/http"
)

var (
	rateLimiter = ratelimit.New(5)
	RabbitMQ    queue.MessageQueuer
)

func RbHandler(c *gin.Context) {
	if !rateLimiter.Take().IsZero() {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too Many Requests"})
		return
	}

	err := RabbitMQ.Publish("user_rate_limit_exceeded", []byte("User rate limit exceeded"))
	if err != nil {
		log.Fatal("error public rabbit", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Request processed successfully"})
}

func StartNotificationService() {
	messages, err := RabbitMQ.Subscribe("user_rate_limit_exceeded")
	if err != nil {
		log.Fatal("Ошибка подписки на канал сообщений RabbitMQ:", err)
	}

	processMessages(messages)
}

func processMessages(messages <-chan queue.Message) {
	for msg := range messages {
		msgString := msg

		sendEmailNotification(msgString)
		sendSMSNotification(msgString)
	}
}

func sendEmailNotification(msg queue.Message) {
	log.Println("Email:", msg)
}

func sendSMSNotification(msg queue.Message) {
	log.Println("SMS:", msg)
}
