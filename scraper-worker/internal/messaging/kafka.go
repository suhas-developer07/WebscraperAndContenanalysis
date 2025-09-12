package messaging

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/suhas-developer07/webScraperContentAnalysis/scraper-worker/internal/models.go"
)

type KafkaProducer struct {
	producer *kafka.Producer
	topic    string
}

func NewKafkaProducer(bootstrapServers, topic string) (*KafkaProducer, error) {

	configMap := &kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"acks":              "all",
		"retries":           5,
		"retry.backoff.ms":  1000,
		"linger.ms":         10,
		"compression.type":  "snappy",
	}

	producer, err := kafka.NewProducer(configMap)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer : %w", err)
	}

	log.Println("Connected to Kafka...")

	/*
		    starts goroutine to handle delivery reports (success/errors).
			this function is critical for understanding fate of our message
	*/
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("FAILED to deliver message to [%s]: %v\n", ev.TopicPartition, ev.TopicPartition.Error)
				} else {
					log.Printf("Successfully delivered message to [%s] (offset %d)\n", *ev.TopicPartition.Topic, ev.TopicPartition.Offset)
				}
			case kafka.Error:
				log.Printf("kafka producer error: %v\n", ev)

			default:
				log.Printf("Ignored event: %v\n", ev)
			}
		}
	}()
	return &KafkaProducer{
		producer: producer,
		topic:    topic,
	}, nil
}

func (kp *KafkaProducer) SendResult(result models.ScrapedResult) error {
	jsonData, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to Marshal result to JSON for taskID %d:%w", result.TaskID, err)
	}

	if kp.producer == nil {
		return fmt.Errorf("kafka producer is not initialized")
	}

	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &kp.topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(fmt.Sprintf("%d", result.JobID)),
		Value: jsonData,
	}

	/*
		     Producer channel full / network timeout
			 The Produce() method is asynchronous. It adds the message to an internal queue.
			 If the queue is full (e.g., due to network issues), it will block.
			 We use a timeout to prevent hanging forever.
	*/

	//channel to signal operation completion
	done := make(chan error, 1)

	//using go routine to handle potentially blocking call
	go func() {
		// This will block if the internal producer queue is full.
		// It's a good thing - it applies backpressure!
		err := kp.producer.Produce(message, nil) // second arg is delivery channel ,we can use global channels instead

		done <- err
	}()

	// waits for the produce call to finish or timeout
	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("failed to produce message for TaskID %d:%w", result.TaskID, err)
		}
		return nil
	case <-time.After(10 * time.Second):
		return fmt.Errorf("timeout while trying to produce message for TASKID %d (producer queue may be full)", result.TaskID)
	}
}

// it must be execute before shutting down the application
func (kp *KafkaProducer) Close() {
	if kp.producer != nil {
		log.Println("Flushing pending kafka messages before shutdown..")

		remaining := kp.producer.Flush(15 * 1000)
		if remaining > 0 {
			log.Printf("WARNING: %d messages were not delivered during flush", remaining)
		}
		kp.producer.Close()
		log.Println("kafka producer closed")
	}
}
