package amqp

import (
	"github.com/pkg/errors"
	amqpDriver "github.com/rabbitmq/amqp091-go"
	"go.k6.io/k6/js/common"
)

// ExchangeOptions defines configuration settings for accessing an exchange.
type ExchangeOptions struct {
	ConnectionURL string `json:"connection_url"`
}

// ExchangeDeclareOptions provides options when declaring (creating) an exchange.
type ExchangeDeclareOptions struct {
	Args       amqpDriver.Table `json:"args"`
	Name       string           `json:"name"`
	Kind       string           `json:"kind"`
	Durable    bool             `json:"durable"`
	AutoDelete bool             `json:"auto_delete"`
	Internal   bool             `json:"internal"`
	NoWait     bool             `json:"no_wait"`
}

// ExchangeDeleteOptions provides options when deleting an exchange.
type ExchangeDeleteOptions struct {
	Name string `json:"name"`
}

// ExchangeBindOptions provides options when binding (subscribing) one exchange to another.
type ExchangeBindOptions struct {
	Args                    amqpDriver.Table `json:"args"`
	DestinationExchangeName string           `json:"destination_exchange_name"`
	SourceExchangeName      string           `json:"source_exchange_name"`
	RoutingKey              string           `json:"routing_key"`
	NoWait                  bool             `json:"no_wait"`
}

// ExchangeUnbindOptions provides options when unbinding (unsubscribing) one exchange from another.
type ExchangeUnbindOptions struct {
	Args                    amqpDriver.Table `json:"args"`
	DestinationExchangeName string           `json:"destination_exchange_name"`
	SourceExchangeName      string           `json:"source_exchange_name"`
	RoutingKey              string           `json:"routing_key"`
	NoWait                  bool             `json:"no_wait"`
}

// declareExchange creates a new exchange given the provided options.
func (a *AMQP) declareExchange(conn *amqpDriver.Connection, options *ExchangeDeclareOptions) {
	ch, err := conn.Channel()
	if err != nil {
		wrappedError := errors.Wrap(ErrGetChannel, err.Error())
		logger.WithField("connection", conn).Error(wrappedError)
		common.Throw(a.vu.Runtime(), wrappedError)
	}

	defer func() {
		_ = ch.Close()
	}()

	err = ch.ExchangeDeclare(
		options.Name,
		options.Kind,
		options.Durable,
		options.AutoDelete,
		options.Internal,
		options.NoWait,
		options.Args,
	)
	if err != nil {
		wrappedError := errors.Wrap(ErrDeclareExchange, err.Error())
		logger.WithField("channel", ch).Error(wrappedError)
		common.Throw(a.vu.Runtime(), wrappedError)
	}
}

// bindExchange subscribes one exchange to another.
func (a *AMQP) bindExchange(conn *amqpDriver.Connection, options *ExchangeBindOptions) {
	ch, err := conn.Channel()
	if err != nil {
		wrappedError := errors.Wrap(ErrGetChannel, err.Error())
		logger.WithField("connection", conn).Error(wrappedError)
		common.Throw(a.vu.Runtime(), wrappedError)
	}

	defer func() {
		_ = ch.Close()
	}()

	err = ch.ExchangeBind(
		options.DestinationExchangeName,
		options.RoutingKey,
		options.SourceExchangeName,
		options.NoWait,
		options.Args,
	)
	if err != nil {
		wrappedError := errors.Wrap(ErrBindExchange, err.Error())
		logger.WithField("channel", ch).Error(wrappedError)
		common.Throw(a.vu.Runtime(), wrappedError)
	}
}

// unbindExchange removes a subscription from one exchange to another.
func (a *AMQP) unbindExchange(conn *amqpDriver.Connection, options *ExchangeUnbindOptions) {
	ch, err := conn.Channel()
	if err != nil {
		wrappedError := errors.Wrap(ErrGetChannel, err.Error())
		logger.WithField("connection", conn).Error(wrappedError)
		common.Throw(a.vu.Runtime(), wrappedError)
	}

	defer func() {
		_ = ch.Close()
	}()

	err = ch.ExchangeUnbind(
		options.DestinationExchangeName,
		options.RoutingKey,
		options.SourceExchangeName,
		options.NoWait,
		options.Args,
	)
	if err != nil {
		wrappedError := errors.Wrap(ErrUnbindExchange, err.Error())
		logger.WithField("channel", ch).Error(wrappedError)
		common.Throw(a.vu.Runtime(), wrappedError)
	}
}

// deleteExchange removes an exchange from the remote server given the exchange name.
func (a *AMQP) deleteExchange(conn *amqpDriver.Connection, options *ExchangeDeleteOptions) {
	ch, err := conn.Channel()
	if err != nil {
		wrappedError := errors.Wrap(ErrGetChannel, err.Error())
		logger.WithField("connection", conn).Error(wrappedError)
		common.Throw(a.vu.Runtime(), wrappedError)
	}

	defer func() {
		_ = ch.Close()
	}()

	err = ch.ExchangeDelete(
		options.Name,
		false, // ifUnused
		false, // noWait
	)
	if err != nil {
		wrappedError := errors.Wrap(ErrDeleteExchange, err.Error())
		logger.WithField("channel", ch).Error(wrappedError)
		common.Throw(a.vu.Runtime(), wrappedError)
	}
}
