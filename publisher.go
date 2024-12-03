package amqp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/grafana/sobek"
	"github.com/pkg/errors"
	amqpDriver "github.com/rabbitmq/amqp091-go"
	"github.com/vmihailenco/msgpack/v5"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/metrics"
)

const (
	jsonContentType = "application/json"
	messagepackType = "application/x-msgpack"
)

// Message is a struct for a publishing message.
type Message struct {
	Headers     map[string]any `json:"headers"`
	Exchange    string         `json:"exchange"`
	RoutingKey  string         `json:"routing_key"`
	Body        string         `json:"body"`
	ContentType string         `json:"content_type"`
	Mandatory   bool           `json:"mandatory"`
	Immediate   bool           `json:"immediate"`
}

// Publisher is a main struct for publish messages to RabbitMQ.
type Publisher struct {
	Channel
	conn   *amqpDriver.Connection
	config PublisherConfig
}

// PublisherConfig is used for construct Publisher.
type PublisherConfig struct {
	ConnectionURL string                 `json:"connection_url"`
	Exchange      ExchangeDeclareOptions `json:"exchange"`
	Queue         QueueDeclareOptions    `json:"queue"`
}

// Channel is a struct used to prepare publishing.
type Channel struct {
	amqpChannel *amqpDriver.Channel
	exchange    string
	queue       string
}

// Close is a method of Publisher for safely shutdown it.
func (p *Publisher) Close() error {
	err := p.Channel.amqpChannel.Close()
	if err != nil {
		return err
	}
	return p.conn.Close()
}

// publisherClass is a wrapper around amqp091-go publishing methods and acts as a JS constructor
// for this extension, thus it must be called with new operator, e.g. new Publisher(...).
func (a *AMQP) publisherClass(call sobek.ConstructorCall) *sobek.Object {
	runtime := a.vu.Runtime()
	var publisherConfig PublisherConfig
	if len(call.Arguments) == 0 {
		common.Throw(runtime, ErrNotEnoughArguments)
	}

	if params, ok := call.Argument(0).Export().(map[string]interface{}); ok {
		if b, err := json.Marshal(params); err != nil {
			common.Throw(runtime, err)
		} else {
			if err = json.Unmarshal(b, &publisherConfig); err != nil {
				common.Throw(runtime, err)
			}
		}
	}

	publisher := a.publisher(publisherConfig)

	publisherObject := runtime.NewObject()
	// This is the publisher object itself.
	if err := publisherObject.Set("This", publisher); err != nil {
		common.Throw(runtime, err)
	}

	err := publisherObject.Set("publish", func(call sobek.FunctionCall) sobek.Value {
		var message *Message
		if len(call.Arguments) == 0 {
			common.Throw(runtime, ErrNotEnoughArguments)
		}

		if params, ok := call.Argument(0).Export().(map[string]interface{}); ok {
			b, err := json.Marshal(params)
			if err != nil {
				common.Throw(runtime, err)
			}
			err = json.Unmarshal(b, &message)
			if err != nil {
				common.Throw(runtime, err)
			}
		}

		a.publish(publisher, message)
		return sobek.Undefined()
	})
	if err != nil {
		common.Throw(runtime, err)
	}

	// This is unnecessary, but it's here for reference purposes.
	err = publisherObject.Set("close", func(_ sobek.FunctionCall) sobek.Value {
		if errInt := publisher.Close(); errInt != nil {
			common.Throw(runtime, errInt)
		}

		return sobek.Undefined()
	})
	if err != nil {
		common.Throw(runtime, err)
	}

	freeze(publisherObject)

	return runtime.ToValue(publisherObject).ToObject(runtime)
}

// publisher creates a new AMQP publisher.
func (a *AMQP) publisher(config PublisherConfig) *Publisher {
	conn, err := amqpDriver.Dial(config.ConnectionURL)
	if err != nil {
		common.Throw(a.vu.Runtime(), err)
		return nil
	}

	ch, err := conn.Channel()
	if err != nil {
		common.Throw(a.vu.Runtime(), err)
		return nil
	}

	err = ch.ExchangeDeclare(
		config.Exchange.Name,
		config.Exchange.Kind,
		config.Exchange.Durable,
		config.Exchange.AutoDelete,
		config.Exchange.Internal,
		config.Exchange.NoWait,
		config.Exchange.Args,
	)
	if err != nil {
		common.Throw(a.vu.Runtime(), err)
		return nil
	}

	q, err := ch.QueueDeclare(
		config.Queue.Name,
		config.Queue.Durable,
		config.Queue.DeleteWhenUnused,
		config.Queue.Exclusive,
		config.Queue.NoWait,
		config.Queue.Args,
	)
	if err != nil {
		common.Throw(a.vu.Runtime(), err)
		return nil
	}

	err = ch.QueueBind(
		config.Queue.Name,
		config.Queue.RoutingKey,
		config.Exchange.Name,
		config.Queue.NoWait,
		config.Queue.Args,
	)
	if err != nil {
		wrappedError := errors.Wrap(ErrBindQueue, err.Error())
		logger.WithField("channel", ch).Error(wrappedError)
		common.Throw(a.vu.Runtime(), wrappedError)
	}

	return &Publisher{
		conn:   conn,
		config: config,
		Channel: Channel{
			amqpChannel: ch,
			exchange:    config.Exchange.Name,
			queue:       q.Name,
		},
	}
}

// publish sends messages to RabbitMQ with the given configuration.
func (a *AMQP) publish(publisher *Publisher, message *Message) { //nolint:funlen
	if state := a.vu.State(); state == nil {
		logger.WithField("state", a.vu.State()).Error(ErrForbiddenInInitContext)
		common.Throw(a.vu.Runtime(), ErrForbiddenInInitContext)
	}

	var ctx context.Context
	if ctx = a.vu.Context(); ctx == nil {
		err := ErrEmptyContext
		logger.WithField("vu.Context()", ctx).Error(err)
		common.Throw(a.vu.Runtime(), err)
	}

	publishing := amqpDriver.Publishing{
		Headers: message.Headers,
	}

	if message.ContentType != "" {
		publishing.ContentType = message.ContentType
	} else {
		publishing.ContentType = jsonContentType
	}

	exch := publisher.exchange
	rkey := publisher.queue

	if message.Exchange != "" {
		exch = message.Exchange
	}

	if message.RoutingKey != "" {
		rkey = message.RoutingKey
	}

	if message.ContentType == messagepackType {
		var (
			jsonParsedBody any
			err            error
		)

		if err = json.Unmarshal([]byte(message.Body), &jsonParsedBody); err != nil {
			logger.WithField("body", message.Body).Error(err)
			common.Throw(a.vu.Runtime(), err)
		}

		publishing.Body, err = msgpack.Marshal(jsonParsedBody)
		if err != nil {
			logger.WithField("body", message.Body).Error(err)
			common.Throw(a.vu.Runtime(), err)
		}
	} else {
		publishing.Body = []byte(message.Body)
	}

	publishing.Body = []byte(message.Body)

	originalErr := publisher.amqpChannel.PublishWithContext(
		a.vu.Context(),
		exch,
		rkey,
		message.Mandatory,
		message.Immediate,
		publishing,
	)

	if originalErr != nil {
		err := fmt.Errorf("error writing messages: %w", originalErr)
		logger.WithField("message", publishing).Error(err)
		common.Throw(a.vu.Runtime(), err)
	}

	q, err := publisher.amqpChannel.QueueDeclarePassive(
		publisher.config.Queue.Name,
		publisher.config.Queue.Durable,
		publisher.config.Queue.DeleteWhenUnused,
		publisher.config.Queue.Exclusive,
		publisher.config.Queue.NoWait,
		publisher.config.Queue.Args,
	)
	if err != nil {
		wrappedError := errors.Wrap(ErrGetQueueStats, err.Error())
		logger.WithField("channel", publisher.amqpChannel).Error(wrappedError)
		common.Throw(a.vu.Runtime(), wrappedError)
	}

	a.reportPublisherStats(q)
}

// reportPublisherStats reports the writer stats to the state.
func (a *AMQP) reportPublisherStats(queue amqpDriver.Queue) {
	state := a.vu.State()
	if state == nil {
		logger.WithField("state", a.vu.State()).Error(ErrForbiddenInInitContext)
		common.Throw(a.vu.Runtime(), ErrForbiddenInInitContext)
	}

	ctx := a.vu.Context()
	if ctx == nil {
		err := ErrEmptyContext
		logger.WithField("vu.Context()", ctx).Error(err)
		common.Throw(a.vu.Runtime(), err)
	}

	ctm := a.vu.State().Tags.GetCurrentValues()
	sampleTags := ctm.Tags.With("queue", queue.Name)

	now := time.Now()

	metrics.PushIfNotDone(ctx, state.Samples, metrics.ConnectedSamples{
		Samples: []metrics.Sample{
			{
				Time: now,
				TimeSeries: metrics.TimeSeries{
					Metric: a.metrics.PublisherQueueMessages,
					Tags:   sampleTags,
				},
				Value:    float64(1),
				Metadata: ctm.Metadata,
			},
			{
				Time: now,
				TimeSeries: metrics.TimeSeries{
					Metric: a.metrics.PublisherQueueLoad,
					Tags:   sampleTags,
				},
				Value:    float64(queue.Messages),
				Metadata: ctm.Metadata,
			},
		},
		Tags: sampleTags,
		Time: now,
	})
}
