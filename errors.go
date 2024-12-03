package amqp

import "errors"

var (
	// ErrForbiddenInInitContext is used when a function is called during VU initialization.
	ErrForbiddenInInitContext = errors.New("publishing RabbitMQ messages in the init context is not supported")

	// ErrNotEnoughArguments is used when a function is called with too few arguments.
	ErrNotEnoughArguments = errors.New("not enough arguments")

	// ErrEmptyContext is used when a function is called with empty vu context.
	ErrEmptyContext = errors.New("empty vu context")

	// ErrGetChannel is used when a function cannot get a channel from connection.
	ErrGetChannel = errors.New("failed to get channel from connection")

	// ErrDeclareExchange is used when a function cannot declare exchange.
	ErrDeclareExchange = errors.New("failed to declare exchange")

	// ErrDeleteExchange is used when a function cannot delete exchange.
	ErrDeleteExchange = errors.New("failed to delete exchange")

	// ErrBindExchange is used when a function cannot bind exchange.
	ErrBindExchange = errors.New("failed to bind exchange")

	// ErrUnbindExchange is used when a function cannot unbind exchange.
	ErrUnbindExchange = errors.New("failed to unbind exchange")

	// ErrDeclareQueue is used when a function cannot declare queue.
	ErrDeclareQueue = errors.New("failed to declare queue")

	// ErrDeleteQueue is used when a function cannot delete queue.
	ErrDeleteQueue = errors.New("failed to delete queue")

	// ErrBindQueue is used when a function cannot bind queue.
	ErrBindQueue = errors.New("failed to bind queue")

	// ErrUnbindQueue is used when a function cannot unbind queue.
	ErrUnbindQueue = errors.New("failed to unbind queue")

	// ErrInspectQueue is used when a function cannot inspect queue.
	ErrInspectQueue = errors.New("failed to inspect queue")

	// ErrPurgeQueue is used when a function cannot purge queue.
	ErrPurgeQueue = errors.New("failed to purge queue")

	// ErrGetQueueStats is used when a function cannot get queue stats.
	ErrGetQueueStats = errors.New("failed to get queue stats")
)
