How to Set Up the Kafka Consumer 📥
After creating the producer, the next step is to set up a consumer that listens to the "notifications" topic and provides an endpoint to list notifications for a specific user.

Let’s move to the cmd/consumer directory and create a new file named consumer.go. Within it, you’ll set up the consumer logic and the Gin-based API:

Code structure:
- Imports
- constants declaration - consumer-group, consumerTopic, ConsumerPort,
						  KafkaServerAddress
- Helper Functions:
	ErrMessagesFound - variable
	getUserIDFromRequest(ctx *gin.Context) (string, error)

- Notificaion storage Model -related content
	UserNotifications map
	NotificationStore struct -have mutex (This struct can be used as class)
	func (ns *NotificationStore) Add(userID string,
	notification models.Notification)

- Kafka related content
  consumer class - having notification store
  func (*Consumer) Setup(sarama.ConsumerGroupSession) error
  func (*Consumer) Cleanup(sarama.ConsumerGroupSession)
  func (consumer *Consumer) ConsumeClaim
  func initializeConsumerGroup()(sarama.ConsumerGroup, error)
  func setupConsumerGroup(ctx context.Context, store *NotificationStore)
  func handleNotifications(ctx *gin.Context, store *NotificationStore)
  

Main

store<=NotificationStore<-{data<-make(UserNotifications)

main-> setupConsumer-> InitializeConsumer:
main-> setupConsumer -> consume Claim -> Add
main -> getNotification


Let’s examine the Kafka-related operations within consumer.go:

Within the initializeConsumerGroup() function:

* config := sarama.NewConfig(): Initializes a new default configuration for Kafka.
* consumerGroup, err := sarama.NewConsumerGroup(…): Creates a new Kafka consumer group that connects to the broker running on localhost:9092. The group name is "notifications-group".

Inside the Consumer struct and its methods:

The Consumer struct has a store field, which is a reference to the NotificationStore to keep track of the received notifications.
*Setup() and Cleanup() methods are required to satisfy the sarama.ConsumerGroupHandler  interface. While they will NOT be used in this tutorial, they can serve  potential roles for initialization and cleanup during message consumption but act as placeholders here.
*In the ConsumeClaim() method: The consumer listens for new messages on the topic. For each message, it fetches the userID (the Key of the message), un-marshals the message into a Notification struct, and adds the notification to the NotificationStore.

In the setupConsumerGroup() function:

* This function sets up the Kafka consumer group, listens for incoming messages, and processes them using the Consumer struct methods.
* It runs a for loop indefinitely, consuming messages from the “notifications” topic and processing any errors that arise.

The handleNotifications() function:

* Initially, it attempts to retrieve the userID from the request. If it doesn’t exist, it returns a 404 Not Found status.
*Then, it fetches the notifications for the provided user ID from the NotificationStore.  Depending on whether the user has notifications, it responds with a 200 OK status and either an empty notifications slice or sends back the current notifications.


Lastly, within the main() function:

*store := &NotificationStore{…}: Creates an instance of NotificationStore to hold the notifications
* ctx, cancel := context.WithCancel(context.Background()): Sets up a cancellable context that can be used to stop the consumer group.
go setupConsumerGroup(ctx, store): Starts the consumer group in a separate Goroutine, allowing it to operate concurrently without blocking the main thread.
The last step is to create a Gin router and define a GET endpoint /notifications/:userID that will fetch the notifications for a specific user via the handleNotifications() function when accessed.

