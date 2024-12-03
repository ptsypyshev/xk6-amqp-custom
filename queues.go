package amqp

import (
	"github.com/pkg/errors"
	amqpDriver "github.com/rabbitmq/amqp091-go"
	"go.k6.io/k6/js/common"
)

// QueueOptions defines configuration settings for accessing a queue.
type QueueOptions struct {
	ConnectionURL string `json:"connection_url"`
}

// QueueDeclareOptions provides queue options when declaring (creating) a queue.
type QueueDeclareOptions struct {
	Args             amqpDriver.Table `json:"args"`
	RoutingKey       string           `json:"routing_key"`
	Name             string           `json:"name"`
	Durable          bool             `json:"durable"`
	DeleteWhenUnused bool             `json:"delete_when_unused"`
	Exclusive        bool             `json:"exclusive"`
	NoWait           bool             `json:"no_wait"`
}

// QueueInspectOptions provide options when inspecting a queue.
type QueueInspectOptions struct {
	Args             amqpDriver.Table `json:"args"`
	Name             string           `json:"name"`
	Durable          bool             `json:"durable"`
	DeleteWhenUnused bool             `json:"delete_when_unused"`
	Exclusive        bool             `json:"exclusive"`
	NoWait           bool             `json:"no_wait"`
}

// QueueDeleteOptions provide options when deleting a queue.
type QueueDeleteOptions struct {
	Name string `json:"name"`
}

// QueueBindOptions provides options when binding a queue to an exchange in order to receive message(s).
type QueueBindOptions struct {
	Args         amqpDriver.Table `json:"args"`
	QueueName    string           `json:"queue_name"`
	ExchangeName string           `json:"exchange_name"`
	RoutingKey   string           `json:"routing_key"`
	NoWait       bool             `json:"no_wait"`
}

// QueueUnbindOptions provides options when unbinding a queue from an exchange to stop receiving message(s).
type QueueUnbindOptions struct {
	Args         amqpDriver.Table `json:"args"`
	QueueName    string           `json:"queue_name"`
	ExchangeName string           `json:"exchange_name"`
	RoutingKey   string           `json:"routing_key"`
}

// QueuePurgeOptions provide options when purging (emptying) a queue.
type QueuePurgeOptions struct {
	QueueName string `json:"queue_name"`
	NoWait    bool   `json:"no_wait"`
}

// declareQueue creates a new queue given the provided options.
func (a *AMQP) declareQueue(conn *amqpDriver.Connection, options *QueueDeclareOptions) {
	ch, err := conn.Channel()
	if err != nil {
		wrappedError := errors.Wrap(ErrGetChannel, err.Error())
		logger.WithField("connection", conn).Error(wrappedError)
		common.Throw(a.vu.Runtime(), wrappedError)
	}

	defer func() {
		_ = ch.Close()
	}()

	_, err = ch.QueueDeclare(
		options.Name,
		options.Durable,
		options.DeleteWhenUnused,
		options.Exclusive,
		options.NoWait,
		options.Args,
	)
	if err != nil {
		wrappedError := errors.Wrap(ErrDeclareQueue, err.Error())
		logger.WithField("channel", ch).Error(wrappedError)
		common.Throw(a.vu.Runtime(), wrappedError)
	}
}

// bindQueue subscribes a queue to an exchange in order to receive message(s).
func (a *AMQP) bindQueue(conn *amqpDriver.Connection, options *QueueBindOptions) {
	ch, err := conn.Channel()
	if err != nil {
		wrappedError := errors.Wrap(ErrGetChannel, err.Error())
		logger.WithField("connection", conn).Error(wrappedError)
		common.Throw(a.vu.Runtime(), wrappedError)
	}

	defer func() {
		_ = ch.Close()
	}()

	err = ch.QueueBind(
		options.QueueName,
		options.RoutingKey,
		options.ExchangeName,
		options.NoWait,
		options.Args,
	)
	if err != nil {
		wrappedError := errors.Wrap(ErrBindQueue, err.Error())
		logger.WithField("channel", ch).Error(wrappedError)
		common.Throw(a.vu.Runtime(), wrappedError)
	}
}

// unbindQueue removes a queue subscription from an exchange to discontinue receiving message(s).
func (a *AMQP) unbindQueue(conn *amqpDriver.Connection, options *QueueUnbindOptions) {
	ch, err := conn.Channel()
	if err != nil {
		wrappedError := errors.Wrap(ErrGetChannel, err.Error())
		logger.WithField("connection", conn).Error(wrappedError)
		common.Throw(a.vu.Runtime(), wrappedError)
	}

	defer func() {
		_ = ch.Close()
	}()

	err = ch.QueueUnbind(
		options.QueueName,
		options.RoutingKey,
		options.ExchangeName,
		options.Args,
	)
	if err != nil {
		wrappedError := errors.Wrap(ErrUnbindQueue, err.Error())
		logger.WithField("channel", ch).Error(wrappedError)
		common.Throw(a.vu.Runtime(), wrappedError)
	}
}

// deleteQueue removes a queue from the remote server given the queue name.
func (a *AMQP) deleteQueue(conn *amqpDriver.Connection, options *QueueDeleteOptions) {
	ch, err := conn.Channel()
	if err != nil {
		wrappedError := errors.Wrap(ErrGetChannel, err.Error())
		logger.WithField("connection", conn).Error(wrappedError)
		common.Throw(a.vu.Runtime(), wrappedError)
	}

	defer func() {
		_ = ch.Close()
	}()

	_, err = ch.QueueDelete(
		options.Name,
		false, // ifUnused
		false, // ifEmpty
		false, // noWait
	)
	if err != nil {
		wrappedError := errors.Wrap(ErrDeleteExchange, err.Error())
		logger.WithField("channel", ch).Error(wrappedError)
		common.Throw(a.vu.Runtime(), wrappedError)
	}
}

// inspectQueue provides queue metadata given queue name.
func (a *AMQP) inspectQueue(conn *amqpDriver.Connection, options *QueueInspectOptions) map[string]interface{} {
	ch, err := conn.Channel()
	if err != nil {
		wrappedError := errors.Wrap(ErrGetChannel, err.Error())
		logger.WithField("connection", conn).Error(wrappedError)
		common.Throw(a.vu.Runtime(), wrappedError)
	}

	defer func() {
		_ = ch.Close()
	}()

	q, err := ch.QueueDeclarePassive(
		options.Name,
		options.Durable,
		options.DeleteWhenUnused,
		options.Exclusive,
		options.NoWait,
		options.Args,
	)
	if err != nil {
		wrappedError := errors.Wrap(ErrInspectQueue, err.Error())
		logger.WithField("channel", ch).Error(wrappedError)
		common.Throw(a.vu.Runtime(), wrappedError)
	}

	return map[string]interface{}{
		"name":      q.Name,
		"consumers": q.Consumers,
		"messages":  q.Messages,
	}
}

// purgeQueue removes all non-consumed message(s) from the specified queue.
func (a *AMQP) purgeQueue(conn *amqpDriver.Connection, options *QueuePurgeOptions) int {
	ch, err := conn.Channel()
	if err != nil {
		wrappedError := errors.Wrap(ErrGetChannel, err.Error())
		logger.WithField("connection", conn).Error(wrappedError)
		common.Throw(a.vu.Runtime(), wrappedError)
	}

	defer func() {
		_ = ch.Close()
	}()

	deleted, err := ch.QueuePurge(options.QueueName, options.NoWait)
	if err != nil {
		wrappedError := errors.Wrap(ErrPurgeQueue, err.Error())
		logger.WithField("channel", ch).Error(wrappedError)
		common.Throw(a.vu.Runtime(), wrappedError)
	}

	return deleted
}
