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

	// Get the absolute path to the .env file in the root directory
	envPath := filepath.Join("..", ".env") // Go up one level from current directory

	// Alternative: if you know the exact path relative to your executable
	// envPath := "/path/to/your/project/.env"

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

	client := messaging.NewClient(conn, ch)

	err = client.ConsumeTasks(queueName)

	if err != nil {
		log.Fatalln("Failed to start consuming :", err)
	}

}
