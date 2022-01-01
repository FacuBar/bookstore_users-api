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

	if err = ch.ExchangeDeclare(
		"users",
		"direct",
		true,
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

func (r *RabbitMQ) Publish(routingKey string, e interface{}) error {
	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(e); err != nil {
		return err
	}

	err := r.Channel.Publish(
		"users",
		routingKey,
		false,
		false,
		amqp.Publishing{
			AppId:       "users-rest-server",
			ContentType: "application/x-encoding-gob",
			Body:        b.Bytes(),
			Timestamp:   time.Now(),
		})
	if err != nil {
		return err
	}

	return nil
}

func (r *RabbitMQ) Close() {
	r.Connection.Close()
}
