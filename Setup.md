
## Step 1: Host Kafka

Assuming that you have Docker and Go installed on your machine, let’s create a directory for the project named kafka-notify. Then, you will pull Bitnami’s Kafka Docker image for the Kafka setup, providing a hassle-free installation:

```
mkdir kafka-notify && cd kafka-notify

curl -sSL \
https://raw.githubusercontent.com/bitnami/containers/main/bitnami/kafka/docker-compose.yml > docker-compose.yml

```

Before starting the Kafka broker, there’s a slight modification required in the docker-compose.yml file. Find the following string:

## Step 2: Modify Kafka to host it on local

Before starting the Kafka broker, there’s a slight modification required in the docker-compose.yml file. Find the following string:

```
KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://:9092
```
Replace it with
```
KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092
```
## Step 3: Running the kafka broker

The above change ensures Kafka advertises its listener on localhost, which allows our local Go application to connect seamlessly. Now, you can start the Kafka broker via the following Docker command:

```
docker-compose up -d
```

## Step 4: Setting up folder structure
Next, you’ll need to create a few directories to organize the project files. The cmd/producer and cmd/consumer directories will contain the main application files, and the pkg/models directory will store the model declarations:
```shell
mkdir -p cmd/producer cmd/consumer pkg/models

# for windows
#mkdir cmd\producer
#mkdir cmd\consumer
#mkdir pkg\models

```
## Step 4: Setting up connection from gin- API to kafka via SARAMA

The last step is to initialize Go modules and install external packages. You’ll be using sarama to establish a connection with the Kafka broker and gin to handle the API endpoints for the notification system:

``` shell
go mod init kafka-notify
go get github.com/IBM/sarama github.com/gin-gonic/gin

```

# How to Create the User and Notification Models

## Step 6: Get the models for Objects created

With the workspace set up, the first step is to create the User and Notification structs. Move to the pkg/models directory, then create a new file named models.go and declare these structs within it:

``` go
package models

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

type Notification struct {
    From    User   `json:"from"`
    To      User   `json:"to"`
    Message string `json:"message"`
}
```

## Step 7: Get the producer and consumer code.

# How to Run

1. Start the producer
Open up a terminal, move into the kafka-notify directory, and run the producer with the following command:

```shell
go run cmd/producer/producer.go
```
2. Start the consumer
Open up a second terminal window, navigate to the kafka-notify directory, and start the consumer by running:
``` shell
go run cmd/consumer/consumer.go
```
3. Sending notifications
With both producer and consumer running, you can simulate sending notifications. Open up a third terminal and use the below curl commands to send notifications:

``` shell
curl -X POST http://localhost:8080/send \
-d "fromID=2&toID=1&message=Bruno started following you."
```
User 2 (Bruno) receives a notification from User 1 (Emma):

``` shell
curl -X POST http://localhost:8080/send \
-d "fromID=1&toID=2&message=Emma mentioned you in a comment: 'Great seeing you yesterday, @Bruno!'"

``` 
User 1 (Emma) receives a notification from User 4 (Lena):

``` shell
curl -X POST http://localhost:8080/send \
-d "fromID=4&toID=1&message=Lena liked your post: 'My weekend getaway!'"
```
4. Retrieving notifications
Finally, you can fetch the notifications of a specific user. You can use the below curl commands to fetch notifications:

Retrieving notifications for User 1 (Emma):

curl http://localhost:8081/notifications/1
Output:

{"notifications": [{"from": {"id": 2, "name": "Bruno"}, "to": {"id": 1, "name": "Emma"}, "message": "Bruno started following you."}]}
{"notifications": [{"from": {"id": 4, "name": "Lena"}, "to": {"id": 1, "name": "Emma"}, "message": "Lena liked your post: 'My weekend getaway!'"}]}
In the above output, you see a JSON response listing all the notifications for Emma. As you send more notifications, they accumulate, and you can fetch them using the consumer’s API.
