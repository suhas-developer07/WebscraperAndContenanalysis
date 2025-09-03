package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type Rabbitmq struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMQRepo(conn *amqp.Connection, ch *amqp.Channel) *Rabbitmq {
	return &Rabbitmq{
		conn: conn,
		ch:   ch,
	}
}

type Data struct {
	ID  int64    `json:"id"`
	URL []string `json:"url"`
}

func (rmq *Rabbitmq) Consumer(queueName string) {
	_, err := rmq.ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Panicf("Failed to declare queue: %s", err)
	}

	msgs, err := rmq.ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Panicf("Failed to Consume messages : %s", err)
	}

	log.Printf("waiting for messages on queue :%s\n", queueName)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var data Data
			err := json.Unmarshal(d.Body, &data)
			if err != nil {
				log.Printf("Error decoding JSON :%v", err)
				continue
			}
			processData(data)
		}
	}()

	<-forever
}

func processData(data Data) {
	fmt.Printf("Recieved Data -ID: %d , URL:%v\n", data.ID, data.URL)
	// here i need to send one by one url to worker using for loop
}
