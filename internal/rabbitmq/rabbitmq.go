package rabbitmq

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

type RabbitmqRepo struct {
	conn  *amqp.Connection
	chann *amqp.Channel
}

func NewRabbitmqRepo(conn *amqp.Connection, chann *amqp.Channel) *RabbitmqRepo {
	return &RabbitmqRepo{
		conn:  conn,
		chann: chann,
	}
}

type Data struct {
	ID  int64    `json:"id"`
	URL []string `json:"url"`
}

func (repo *RabbitmqRepo) SendURL(data Data) error {
	queueName := os.Getenv("QUEUE_NAME")
	if queueName == "" {
		return fmt.Errorf("missing QUEUE_NAME env variable")
	}

	_, err := repo.chann.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal struct: %w", err)
	}

	err = repo.chann.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}
