package rabbit

import (
	"dots-api/lib/utils"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func Connect(host string) (*amqp.Connection, *amqp.Channel, error) {
	var (
		err         error
		conn        *amqp.Connection
		amqpChannel *amqp.Channel
	)

	conn, err = amqp.Dial(host)
	handleError(err, utils.ErrConnectAMQP)

	amqpChannel, err = conn.Channel()
	handleError(err, utils.ErrCreateChannelAMQP)

	return conn, amqpChannel, err
}
