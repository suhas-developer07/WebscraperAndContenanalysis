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
	conn   *amqp.Connection
	ch     *amqp.Channel
	worker *worker.Scraper
}

func NewClient(conn *amqp.Connection, ch *amqp.Channel) *Client {
	return &Client{
		conn: conn,
		ch:   ch,
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
			log.Println(job)

			var wg sync.WaitGroup

			for i, url := range job.URLs {
				wg.Add(1)
				go func(i int, url string) {
					defer wg.Done()
					task := models.Task{
						JobID:  job.ID,
						TaskID: i + 1,
						URL:    url,
					}

					result := c.worker.ProcessTask(task)
					log.Println(result)
				}(i, url)
			}
			wg.Wait()
			d.Ack(false)

			log.Printf("Completed processing all tasks for Job ID :%d", job.ID)
		}
	}()

	<-forever
	return nil
}
