package rabbit

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func Connect(host string) (*amqp.Connection, *amqp.Channel, error) {
	var (
		err         error
		conn        *amqp.Connection
		amqpChannel *amqp.Channel
	)

	conn, err = amqp.Dial(host)
	// handleError(err, utils.ErrConnectAMQP)

	if err != nil {
		return nil, nil, err
	}

	amqpChannel, err = conn.Channel()
	// handleError(err, utils.ErrCreateChannelAMQP)

	if err != nil {
		return conn, nil, err
	}

	return conn, amqpChannel, err
}
