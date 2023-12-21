
Package name

Constants - Port, Kafka Address, Topic Name

Helper functions:-
User Not found Error,
FindUserby ID (id int, users []models.User) (models.User, error)
Get ID from Request (formValue string, ctx *gin.Context) (int, error)

Other functions:
sendKafkaMessage(producer sarama.SyncProducer,
	users []models.User, ctx *gin.Context, fromID, toID int)

sendMessageHandler(producer sarama.SyncProducer,
	users []models.User)

setupProducer() (sarama.SyncProducer, error)

Call Flow:
Main -> Setup Producers:
Main -> SendMessageHandler ->sendkafkaMessage

Read it from main.


Within the setupProducer() function:

* config := sarama.NewConfig(): Initializes a new default configuration for Kafka. Think of it as setting up the parameters before connecting to the broker.

* config.Producer.Return.Successes = true: Ensures that the producer receives an acknowledgment 
once the message is successfully stored in the "notifications" topic.

* producer, err := sarama.NewSyncProducer(…): Initializes a synchronous Kafka producer that connects to the Kafka broker running at localhost:9092.

Inside the sendKafkaMessage() function:

* This function starts by retrieving the message from the context and then attempts to find both the sender and the recipient using their IDs.
*notification := models.Notification{…}: Initializes a Notification struct that encapsulates information about the sender, the recipient, and the actual message.
*msg := &sarama.ProducerMessage{…}: Constructs a ProducerMessage for the "notifications" topic, setting the recipient’s ID as the Key and the message content, which is the serialized form of the Notification as the Value.
producer.SendMessage(msg): Sends the constructed message to the "notifications" topic.

In the sendMessageHandler() function:

This function serves as an endpoint handler for the /send POST request. It processes the incoming request to ensure valid sender and recipient IDs are provided.
After fetching the IDs, it invokes the sendKafkaMessage() function to send the Kafka message. Depending on the result, it dispatches appropriate HTTP responses: a 404 Not Found for nonexistent users, a 400 Bad Request for invalid IDs, and a 500 Internal Server Error for other failures, along with a specific error message.

Finally, within the main() function:

You initialize a Kafka producer via the setupProducer() function.
Then, you create a Gin router via gin.Default(), setting up a web server. Next, you define a POST endpoint /send to handle notifications. This endpoint expects the sender and recipient’s IDs and the message content.
The notification is processed upon receiving a POST request via the sendMessageHandler() function, and an appropriate HTTP response is dispatched.
The above setup provides a simple way to simulate users sending notifications to each other and showcases how these notifications are produced to the "notifications" topic.