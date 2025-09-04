package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"github.com/suhas-developer07/webScraperContentAnalysis/scraper-worker/internal/messaging"
)

// Things to do
// 1. consume messages from rabbitmq with this "urls" queue name
// 2. take the consumed msgs and give it one by one to worker (Goscolly) it extracts text from it
// 3. put the extracted text to kafka with message like this ("task_id": 123, "url": "...", "row_text": "...")

func main() {
	fmt.Print("hello world")

	envPath := filepath.Join("..", ".env") // Go up one level from current directory

	// Load the .env file from the specified path
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading .env file from %s: %v", envPath, err)
	}

	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		log.Fatal("RABBITMQ_URL environment variable is not set")
	}

	queueName := os.Getenv("QUEUE_NAME")
	if queueName == "" {
		log.Fatal("queueName is not present in the env file ")
	}

	conn, err := amqp.Dial(rabbitmqURL)

	if err != nil {
		log.Fatalln("Failed to Connect rabbitMQ", err)
	}

	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		log.Fatalln("Failed to open Channel", err)
	}

	defer ch.Close()

	bootstrapServers := os.Getenv("KAFKA_BOOTSTRAP_SERVICE")
	kafkaTopic := os.Getenv("KAFKA_RAW_CONTENT_TOPIC")

	kafkaProducer, err := messaging.NewKafkaProducer(bootstrapServers, kafkaTopic)

	if err != nil {
		log.Fatalf("Failed to create kafka producer : %v", err)
	}

	defer kafkaProducer.Close()

	client := messaging.NewClient(conn, ch, kafkaProducer)

	err = client.ConsumeTasks(queueName)

	if err != nil {
		log.Fatalln("Failed to start consuming :", err)
	}

}
