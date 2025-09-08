package messaging

import (
	"encoding/json"
	"sync"

	"log"

	"github.com/streadway/amqp"
	"github.com/suhas-developer07/webScraperContentAnalysis/scraper-worker/internal/models.go"
	"github.com/suhas-developer07/webScraperContentAnalysis/scraper-worker/internal/worker"
)

type Client struct {
	conn          *amqp.Connection
	ch            *amqp.Channel
	worker        *worker.Scraper
	kafkaProducer *KafkaProducer
}

func NewClient(conn *amqp.Connection, ch *amqp.Channel, KafkaProducer *KafkaProducer) *Client {
	return &Client{
		conn:          conn,
		ch:            ch,
		kafkaProducer: KafkaProducer,
	}
}

func (c *Client) ConsumeTasks(queueName string) error {
	_, err := c.ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	msgs, err := c.ch.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	log.Printf("waiting for messages on queue :%s\n", queueName)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var job models.Job
			err := json.Unmarshal(d.Body, &job)
			if err != nil {
				log.Printf("Error decoding JSON :%v", err)
				d.Nack(false, true)
				continue
			}

			var wg sync.WaitGroup

			resultCh := make(chan models.ScrapedResult, len(job.URLs)) // creating the channel to collect all results from the goroutines

			//process the task cuncurrently
			for i, url := range job.URLs {
				wg.Add(1)
				go func(taskID int, url string) {
					defer wg.Done()
					task := models.Task{
						JobID:  job.ID,
						TaskID: taskID,
						URL:    url,
					}

					result := c.worker.ProcessTask(task)

					resultCh <- result //sending the result to channel
					log.Printf("Processed URL %d/%d for Job %d", taskID, len(job.URLs), job.ID)
				}(i+1, url)
			}

			//wait for all goroutine to complete in a seperate gorautine
			//so we can close the result channel
			go func() {
				wg.Wait()
				close(resultCh)
			}()

			// Collect and process all results from the channel
			// This happens in the main goroutine for this message
			for result := range resultCh {

				if result.RawText == "" {
					continue
					//TODO:something i should do here
					// rather than sending empty extracted text to kafka drop it here itself
				}
				err := c.kafkaProducer.SendResult(result)
				if err != nil {
					log.Printf("ERROR: Failed tp send result to kafka for Task %d:%v.Might need tp retry or push to DLQ", result.TaskID, err)
				}
				log.Printf("Task successfully sended to kafka:%d", result.TaskID)
			}
			d.Ack(false)

			log.Printf("Completed processing all tasks for Job ID :%d", job.ID)
		}
	}()

	<-forever
	return nil
}
