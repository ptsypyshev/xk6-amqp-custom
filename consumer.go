package amqp

import (
	"context"
	"encoding/json"
	"time"

	"github.com/grafana/sobek"
	"github.com/pkg/errors"
	amqpDriver "github.com/rabbitmq/amqp091-go"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/metrics"
)

// Consumer is a main struct for consume messages to RabbitMQ.
type Consumer struct {
	Channel
	conn   *amqpDriver.Connection
	config ConsumerConfig
}

// ConsumerConfig is used for construct Consumer.
type ConsumerConfig struct {
	ConnectionURL string                 `json:"connection_url"`
	Exchange      ExchangeDeclareOptions `json:"exchange"`
	Queue         QueueDeclareOptions    `json:"queue"`
}

// ConsumeConfig is used for pass params to consume method.
type ConsumeConfig struct {
	ConsumeLimit int      `json:"consume_limit"`
	ReadTimeout  Duration `json:"read_timeout"`
}

// Close is a method of Publisher for safely shutdown it.
func (c *Consumer) Close() error {
	err := c.Channel.amqpChannel.Close()
	if err != nil {
		return err
	}
	return c.conn.Close()
}

// Duration is a custom wrapper for time.Duration type.
type Duration struct {
	time.Duration
}

// MarshalJSON is a wrapper of json.Marshal.
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// UnmarshalJSON is a wrapper of json.Unmarshal.
func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	switch value := v.(type) {
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}

// consumerClass is a wrapper around amqp091-go consuming methods and acts as a JS constructor
// for this extension, thus it must be called with new operator, e.g. new Consumer(...).
func (a *AMQP) consumerClass(call sobek.ConstructorCall) *sobek.Object {
	runtime := a.vu.Runtime()
	var consumerConfig ConsumerConfig
	if len(call.Arguments) == 0 {
		common.Throw(runtime, ErrNotEnoughArguments)
	}

	if params, ok := call.Argument(0).Export().(map[string]interface{}); ok {
		if b, err := json.Marshal(params); err != nil {
			common.Throw(runtime, err)
		} else {
			if err = json.Unmarshal(b, &consumerConfig); err != nil {
				common.Throw(runtime, err)
			}
		}
	}

	consumer := a.consumer(consumerConfig)

	consumerObject := runtime.NewObject()
	// This is the reader object itself
	if err := consumerObject.Set("This", consumer); err != nil {
		common.Throw(runtime, err)
	}

	err := consumerObject.Set("consume", func(call sobek.FunctionCall) sobek.Value {
		var consumeConfig *ConsumeConfig
		if len(call.Arguments) == 0 {
			common.Throw(runtime, ErrNotEnoughArguments)
		}

		if params, ok := call.Argument(0).Export().(map[string]interface{}); ok {
			if b, err := json.Marshal(params); err != nil {
				common.Throw(runtime, err)
			} else {
				if err = json.Unmarshal(b, &consumeConfig); err != nil {
					common.Throw(runtime, err)
				}
			}
		}

		return runtime.ToValue(a.consume(consumer, consumeConfig))
	})
	if err != nil {
		common.Throw(runtime, err)
	}

	// This is unnecessary, but it's here for reference purposes
	err = consumerObject.Set("close", func(_ sobek.FunctionCall) sobek.Value {
		if errInt := consumer.Close(); errInt != nil {
			common.Throw(runtime, errInt)
		}

		return sobek.Undefined()
	})
	if err != nil {
		common.Throw(runtime, err)
	}

	freeze(consumerObject)

	return runtime.ToValue(consumerObject).ToObject(runtime)
}

// consumer creates a RabbitMQ consumer with the given configuration
func (a *AMQP) consumer(config ConsumerConfig) *Consumer {
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
		common.Throw(a.vu.Runtime(), err)
		return nil
	}

	return &Consumer{
		conn:   conn,
		config: config,
		Channel: Channel{
			amqpChannel: ch,
			exchange:    config.Exchange.Name,
			queue:       q.Name,
		},
	}
}

// consume consumes messages from the given consumer.
func (a *AMQP) consume(consumer *Consumer, config *ConsumeConfig) []map[string]interface{} {
	if state := a.vu.State(); state == nil {
		logger.WithField("state", a.vu.State()).Error(ErrForbiddenInInitContext)
		common.Throw(a.vu.Runtime(), ErrForbiddenInInitContext)
	}

	ctx := a.vu.Context()
	if ctx == nil {
		err := ErrEmptyContext
		logger.WithField("vu.Context()", ctx).Error(err)
		common.Throw(a.vu.Runtime(), err)
	}

	if config.ConsumeLimit <= 0 {
		config.ConsumeLimit = 1
	}

	messages := make([]map[string]interface{}, 0)

	ctxWithTimeout, cancel := context.WithTimeout(ctx, config.ReadTimeout.Duration) // TODO: Maybe it is not necessary
	defer cancel()

	queueChannel, err := consumer.amqpChannel.ConsumeWithContext(
		ctxWithTimeout,
		consumer.config.Queue.Name,
		"",    // TODO: Use config variables
		false, // TODO: Use config variables
		false, // TODO: Use config variables
		false, // TODO: Use config variables
		false, // TODO: Use config variables
		nil,   // TODO: Use config variables
	)
	if err != nil {
		logger.WithField("error", err).Error(err)
		common.Throw(a.vu.Runtime(), err)
	}

	for msg := range queueChannel {
		// Rest of the fields of a given message
		message := map[string]interface{}{
			"exchange": msg.Exchange,
			"headers":  msg.Headers,
			"body":     string(msg.Body),
		}

		messages = append(messages, message)
		err = msg.Ack(false)
		if err != nil {
			logger.WithField("error", err).Error(err)
			common.Throw(a.vu.Runtime(), err)
		}

		if config.ConsumeLimit > 0 && len(messages) >= config.ConsumeLimit {
			break
		}
	}

	q, err := consumer.amqpChannel.QueueDeclarePassive(
		consumer.config.Queue.Name,
		consumer.config.Queue.Durable,
		consumer.config.Queue.DeleteWhenUnused,
		consumer.config.Queue.Exclusive,
		consumer.config.Queue.NoWait,
		consumer.config.Queue.Args,
	)
	if err != nil {
		wrappedError := errors.Wrap(ErrGetQueueStats, err.Error())
		logger.WithField("channel", consumer.amqpChannel).Error(wrappedError)
		common.Throw(a.vu.Runtime(), wrappedError)
	}

	a.reportReaderStats(q, len(messages))

	return messages
}

// reportReaderStats reports the reader stats.
func (a *AMQP) reportReaderStats(queue amqpDriver.Queue, counter int) {
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
					Metric: a.metrics.ConsumerQueueMessages,
					Tags:   sampleTags,
				},
				Value:    float64(counter),
				Metadata: ctm.Metadata,
			},
			{
				Time: now,
				TimeSeries: metrics.TimeSeries{
					Metric: a.metrics.ConsumerQueueLoad,
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
