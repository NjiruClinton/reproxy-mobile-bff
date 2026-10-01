package producer

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/IBM/sarama"
	_ "github.com/IBM/sarama"
	"github.com/NjiruClinton/reproxy-mobile-bff/pkg/models"
	"github.com/gin-gonic/gin"
)

const (
	ProducerPort        = ":8080"
	KafkaServiceAddress = "localhost:9092"
	KafkaTopic          = "notifications"
)

var ErrUserNotFoundInProducer = errors.New("user not found")

func findUserByID(id int, users []models.User) (models.User, error) {
	for _, user := range users {
		if user.ID == id {
			return user, nil
		}
	}
	return models.User{}, ErrUserNotFoundInProducer
}

func getIDFromRequest(formValue string, ctx *gin.Context) (int, error) {
	id, err := strconv.Atoi(ctx.PostForm(formValue))
	if err != nil {
		return 0, fmt.Errorf("failed to parse ID from form value %s: %w", formValue, err)
	}
	return id, nil
}

func sendKafkaMessage(producer sarama.SyncProducer, users []models.User, ctx *gin.Context, fromID, toID int) error {
	message := ctx.PostForm("message")

	fromUser, err := findUserByID(fromID, users)
	if err != nil {
		return err
	}

	toUser, err := findUserByID(toID, users)
	if err != nil {
		return err
	}

	notification := models.Notification{
		From: fromUser,
		To:   toUser, Message: message,
	}

	notificationJSON, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: KafkaTopic,
		Key:   sarama.StringEncoder(strconv.Itoa(toUser.ID)),
		Value: sarama.StringEncoder(notificationJSON),
	}

	_, _, err = producer.SendMessage(msg)
	return err
}

func setupProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{KafkaServiceAddress}, config)
	if err != nil {
		return nil, fmt.Errorf("failed to setup producer: %w", err)
	}
	return producer, nil
}

func main() {
	users := []models.User{
		{ID: 1, Name: "Emma"},
		{ID: 2, Name: "Bruno"},
		{ID: 3, Name: "Rick"},
		{ID: 4, Name: "Lena"},
	}

	producer, err := setupProducer()
	if err != nil {
		log.Fatalf("failed to initialize producer: %v", err)
	}
	defer producer.Close()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.POST("/send", sendMessageHandler(producer, users))

	fmt.Printf("Kafka PRODUCER started at http:localhost:%s\n", ProducerPort)

	if err := router.Run(ProducerPort); err != nil {
		log.Printf("failed to run the server: %v", err)
	}
}
