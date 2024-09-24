package rabbit

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	QueueUserBadge = "dots_user_badge"
	QueueBadges    = "dots_badges"
)

type QueueData struct {
	QueueName string
	Data      interface{}
}

// constructor for QueueData
func QueueDataPayload(queueName string, data interface{}) QueueData {
	return QueueData{
		QueueName: queueName,
		Data:      data,
	}
}

// PublishQueue ...
func PublishQueue(ctx context.Context, host string, data QueueData) error {
	// Init Connect Queue
	conn, ch, err := Connect(host)
	if err != nil {
		return fmt.Errorf("[%s] %s: %v", data.QueueName, "Error when connect to queue", err)
	}
	defer conn.Close()
	defer ch.Close()

	// Declare Queue
	queue, err := ch.QueueDeclare(data.QueueName, false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("[%s] %s: %v", data.QueueName, "Could not declare queue", err)
	}

	// Unmarshall Data
	dataJSON, err := json.Marshal(data.Data)
	if err != nil {
		return fmt.Errorf("[%s] %s: %v", data.QueueName, "Error encoding JSON", err)
	}

	// Publish Queue
	err = ch.PublishWithContext(ctx, "", queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         dataJSON,
	})
	if err != nil {
		return fmt.Errorf("[%s] %s: %v", data.QueueName, "Something error when publish the messages", err)
	}

	return nil
}
