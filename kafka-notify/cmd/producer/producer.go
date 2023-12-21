package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"kafka-notify/pkg/models"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

const (
	ProducerPort       = ":8080"
	KafkaServerAddress = "localhost:9092"
	KafkaTopic         = "notifications"
)

// ============== HELPER FUNCTIONS ==============
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
		return 0, fmt.Errorf(
			"failed to parse ID from form value %s: %w", formValue, err)
	}
	return id, nil
}

// ============== KAFKA RELATED FUNCTIONS ==============
func sendKafkaMessage(producer sarama.SyncProducer,
	users []models.User, ctx *gin.Context, fromID, toID int) error {
	message := ctx.PostForm("message") //message to send

	fromUser, err := findUserByID(fromID, users) //find sender in out own db
	if err != nil {
		return err
	}

	toUser, err := findUserByID(toID, users) //find receiver in our own db
	if err != nil {
		return err
	}

	notification := models.Notification{ //create notification
		From: fromUser,
		To:   toUser, Message: message,
	}

	notificationJSON, err := json.Marshal(notification) //convert to json
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	msg := &sarama.ProducerMessage{ //send message to kafka
		Topic: KafkaTopic,
		Key:   sarama.StringEncoder(strconv.Itoa(toUser.ID)),
		Value: sarama.StringEncoder(notificationJSON),
	}

	_, _, err = producer.SendMessage(msg)
	return err
}

func sendMessageHandler(producer sarama.SyncProducer,
	users []models.User) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fromID, err := getIDFromRequest("fromID", ctx) //sender ID
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		toID, err := getIDFromRequest("toID", ctx) //receiver ID
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		err = sendKafkaMessage(producer, users, ctx, fromID, toID) //send message
		if errors.Is(err, ErrUserNotFoundInProducer) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Notification sent successfully!",
		})
	}
}

func setupProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()            //default config
	config.Producer.Return.Successes = true // enable message acknowledgements
	producer, err := sarama.NewSyncProducer([]string{KafkaServerAddress},
		config) // setup producer
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

	fmt.Printf("Kafka PRODUCER 📨 started at http://localhost%s\n",
		ProducerPort)

	if err := router.Run(ProducerPort); err != nil {
		log.Printf("failed to run the server: %v", err)
	}
}
