
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