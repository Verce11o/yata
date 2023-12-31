package rabbitmq

import (
	"github.com/Verce11o/yata/internal/http/websocket"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type NotificationConsumer struct {
	AmqpConn *amqp.Connection
	log      *zap.SugaredLogger
	trace    trace.Tracer
}

func NewNotificationConsumer(amqpConn *amqp.Connection, log *zap.SugaredLogger, trace trace.Tracer) *NotificationConsumer {
	return &NotificationConsumer{AmqpConn: amqpConn, log: log, trace: trace}
}

func (c *NotificationConsumer) createChannel(exchangeName, queueName, bindingKey string) *amqp.Channel {
	ch, err := c.AmqpConn.Channel()

	if err != nil {
		panic(err)
	}

	err = ch.ExchangeDeclare(
		exchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	queue, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	err = ch.QueueBind(
		queue.Name,
		bindingKey,
		exchangeName,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	return ch

}

func (c *NotificationConsumer) StartConsumer(queueName, consumerTag, exchangeName, bindingKey string, clients websocket.WsClients) error {
	ch := c.createChannel(exchangeName, queueName, bindingKey)
	defer ch.Close()

	deliveries, err := ch.Consume(
		queueName,
		consumerTag,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	for i := 0; i < 5; i++ {
		i := i
		go c.worker(i, deliveries, clients)
	}
	chanErr := <-ch.NotifyClose(make(chan *amqp.Error))
	c.log.Infof("Notify close: %v", chanErr)

	return chanErr

}

func (c *NotificationConsumer) worker(index int, messages <-chan amqp.Delivery, clients websocket.WsClients) {
	//for message := range messages {
	//	c.log.Infof("Worker #%d: %v", index, string(message.Body))
	//
	//	var request domain.IncomingNewTweetNotification // under change
	//
	//	err := json.Unmarshal(message.Body, &request)
	//
	//	if err != nil {
	//		c.log.Errorf("failed to unmarshal request: %v", err)
	//	}
	//
	//	// get all user subscribers
	//
	//	conn, ok := clients[request.FromUserID]
	//
	//	if ok {
	//		// mark notification as read
	//	}
	//
	//	err = conn.WriteJSON(request)
	//
	//	if err != nil {
	//		c.log.Errorf("error sending notification: %v", err.Error())
	//	}
	//
	//	err = message.Ack(false)
	//
	//	if err != nil {
	//		c.log.Errorf("failed to acknowledge delivery: %v", err)
	//	}
	//
	//}
	c.log.Info("Channel closed")
}
