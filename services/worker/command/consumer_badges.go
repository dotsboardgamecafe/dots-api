package command

import (
	"context"
	"dots-api/lib/rabbit"
	"dots-api/services/worker/model"
	"encoding/json"
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/urfave/cli/v2"
)

// Consumer ...
func (app Contract) ConsumerBadges(c *cli.Context) error {
	var (
		err  error
		name = rabbit.QueueBadges
		m    = model.Contract{App: app.App}
	)

	// listen rabbit mq
	conn, err := amqp.Dial(app.App.Config.GetString("queue.rabbitmq.host"))
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %s", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		true,  // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %s", err)
	}

	err = ch.Qos(1, 0, false)
	if err != nil {
		return fmt.Errorf("failed to configure Qos: %s", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		true,   // no-wait
		nil,    // args
	)

	if err != nil {
		return fmt.Errorf("failed to register a consumer: %s", err)
	}

	// Handle messages received on the channel
	stopChan := make(chan bool)
	go func() {
		log.Printf("Consumer ready, PID: %d", os.Getpid())
		for d := range msgs {
			// progress to db
			var data rabbit.QueueBadgeData
			ctx := context.Background()
			if err := json.Unmarshal(d.Body, &data); err != nil {
				log.Printf("Unmarshal body err: %s", err)
			}
			if err := m.CheckBadges(ctx, data.BadgeCode); err != nil {
				log.Printf("CheckBadge err: %s", err)
			}

			// successful processing
			d.Ack(true)
		}
	}()
	<-stopChan
	return nil
}
