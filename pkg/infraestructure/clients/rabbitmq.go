package clients

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	if _, err = ch.QueueDeclare(
		"users",
		false,
		false,
		false,
		false,
		nil,
	); err != nil {
		return nil, err
	}

	return &RabbitMQ{
		Connection: conn,
		Channel:    ch,
	}, nil
}

func (t *RabbitMQ) Publish(routingKey string, e interface{}) error {
	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(e); err != nil {
		return err
	}

	err := t.Channel.Publish(
		"users",
		routingKey,
		false,
		false,
		amqp.Publishing{
			AppId:       "tasks-rest-server",
			ContentType: "application/x-encoding-gob",
			Body:        b.Bytes(),
			Timestamp:   time.Now(),
		})
	if err != nil {
		return err
	}

	return nil
}
