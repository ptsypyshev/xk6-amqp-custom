package amqp

import (
	"errors"

	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/metrics"
)

// Metrics contains metrics for consumer and publisher.
type Metrics struct {
	ConsumerQueueMessages  *metrics.Metric
	ConsumerQueueLoad      *metrics.Metric
	PublisherQueueMessages *metrics.Metric
	PublisherQueueLoad     *metrics.Metric
}

// registerMetrics registers the metrics for the AMQP module in the metrics registry.
func registerMetrics(vu modules.VU) (Metrics, error) {
	var err error
	registry := vu.InitEnv().Registry
	amqpMetrics := Metrics{}

	if amqpMetrics.ConsumerQueueMessages, err = registry.NewMetric(
		"consumer_message_count", metrics.Counter); err != nil {
		return amqpMetrics, errors.Unwrap(err)
	}
	if amqpMetrics.ConsumerQueueLoad, err = registry.NewMetric(
		"consumer_queue_load", metrics.Gauge); err != nil {
		return amqpMetrics, errors.Unwrap(err)
	}

	if amqpMetrics.PublisherQueueMessages, err = registry.NewMetric(
		"publisher_message_count", metrics.Counter); err != nil {
		return amqpMetrics, errors.Unwrap(err)
	}
	if amqpMetrics.PublisherQueueLoad, err = registry.NewMetric(
		"publisher_queue_load", metrics.Gauge); err != nil {
		return amqpMetrics, errors.Unwrap(err)
	}

	return amqpMetrics, nil
}
