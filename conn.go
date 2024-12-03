package amqp

import (
	"encoding/json"

	"github.com/grafana/sobek"
	amqpDriver "github.com/rabbitmq/amqp091-go"
	"go.k6.io/k6/js/common"
)

// ConnectionConfig used to pass connection url to constructor.
type ConnectionConfig struct {
	ConnectionURL string `json:"connection_url"`
}

// connectionClass is a constructor for the Connection object in JS
// that creates a new connection for creating, listing and deleting topics,
// e.g. new Connection(...).
func (a *AMQP) connectionClass(call sobek.ConstructorCall) *sobek.Object { //nolint: funlen,gocognit,cyclop
	runtime := a.vu.Runtime()
	var connectionConfig *ConnectionConfig
	if len(call.Arguments) == 0 {
		common.Throw(runtime, ErrNotEnoughArguments)
	}

	if params, ok := call.Argument(0).Export().(map[string]interface{}); ok {
		if b, err := json.Marshal(params); err != nil {
			common.Throw(runtime, err)
		} else {
			if err = json.Unmarshal(b, &connectionConfig); err != nil {
				common.Throw(runtime, err)
			}
		}
	}

	connection := a.getAMQPConnection(connectionConfig)

	connectionObject := runtime.NewObject()
	// This is the connection object itself
	if err := connectionObject.Set("This", connection); err != nil {
		common.Throw(runtime, err)
	}

	err := connectionObject.Set("declareExchange", func(call sobek.FunctionCall) sobek.Value {
		var exchangeConfig *ExchangeDeclareOptions
		if len(call.Arguments) == 0 {
			common.Throw(runtime, ErrNotEnoughArguments)
		}

		if params, ok := call.Argument(0).Export().(map[string]interface{}); ok {
			if b, err := json.Marshal(params); err != nil {
				common.Throw(runtime, err)
			} else {
				if err = json.Unmarshal(b, &exchangeConfig); err != nil {
					common.Throw(runtime, err)
				}
			}
		}

		a.declareExchange(connection, exchangeConfig)
		return sobek.Undefined()
	})
	if err != nil {
		common.Throw(runtime, err)
	}

	err = connectionObject.Set("bindExchange", func(call sobek.FunctionCall) sobek.Value {
		var exchangeConfig *ExchangeBindOptions
		if len(call.Arguments) == 0 {
			common.Throw(runtime, ErrNotEnoughArguments)
		}

		if params, ok := call.Argument(0).Export().(map[string]interface{}); ok {
			if b, errInt := json.Marshal(params); errInt != nil {
				common.Throw(runtime, errInt)
			} else {
				if errInt = json.Unmarshal(b, &exchangeConfig); errInt != nil {
					common.Throw(runtime, errInt)
				}
			}
		}

		a.bindExchange(connection, exchangeConfig)
		return sobek.Undefined()
	})
	if err != nil {
		common.Throw(runtime, err)
	}

	err = connectionObject.Set("unbindExchange", func(call sobek.FunctionCall) sobek.Value {
		var exchangeConfig *ExchangeUnbindOptions
		if len(call.Arguments) == 0 {
			common.Throw(runtime, ErrNotEnoughArguments)
		}

		if params, ok := call.Argument(0).Export().(map[string]interface{}); ok {
			if b, errInt := json.Marshal(params); errInt != nil {
				common.Throw(runtime, errInt)
			} else {
				if errInt = json.Unmarshal(b, &exchangeConfig); errInt != nil {
					common.Throw(runtime, errInt)
				}
			}
		}

		a.unbindExchange(connection, exchangeConfig)
		return sobek.Undefined()
	})
	if err != nil {
		common.Throw(runtime, err)
	}

	err = connectionObject.Set("deleteExchange", func(call sobek.FunctionCall) sobek.Value {
		var exchangeConfig *ExchangeDeleteOptions
		if len(call.Arguments) == 0 {
			common.Throw(runtime, ErrNotEnoughArguments)
		}

		if params, ok := call.Argument(0).Export().(map[string]interface{}); ok {
			if b, errInt := json.Marshal(params); errInt != nil {
				common.Throw(runtime, errInt)
			} else {
				if errInt = json.Unmarshal(b, &exchangeConfig); errInt != nil {
					common.Throw(runtime, errInt)
				}
			}
		}

		a.deleteExchange(connection, exchangeConfig)
		return sobek.Undefined()
	})
	if err != nil {
		common.Throw(runtime, err)
	}

	err = connectionObject.Set("declareQueue", func(call sobek.FunctionCall) sobek.Value {
		var queueConfig *QueueDeclareOptions
		if len(call.Arguments) == 0 {
			common.Throw(runtime, ErrNotEnoughArguments)
		}

		if params, ok := call.Argument(0).Export().(map[string]interface{}); ok {
			if b, errInt := json.Marshal(params); errInt != nil {
				common.Throw(runtime, errInt)
			} else {
				if errInt = json.Unmarshal(b, &queueConfig); errInt != nil {
					common.Throw(runtime, errInt)
				}
			}
		}

		a.declareQueue(connection, queueConfig)
		return sobek.Undefined()
	})
	if err != nil {
		common.Throw(runtime, err)
	}

	err = connectionObject.Set("bindQueue", func(call sobek.FunctionCall) sobek.Value {
		var queueConfig *QueueBindOptions
		if len(call.Arguments) == 0 {
			common.Throw(runtime, ErrNotEnoughArguments)
		}

		if params, ok := call.Argument(0).Export().(map[string]interface{}); ok {
			if b, errInt := json.Marshal(params); errInt != nil {
				common.Throw(runtime, errInt)
			} else {
				if errInt = json.Unmarshal(b, &queueConfig); errInt != nil {
					common.Throw(runtime, errInt)
				}
			}
		}

		a.bindQueue(connection, queueConfig)
		return sobek.Undefined()
	})
	if err != nil {
		common.Throw(runtime, err)
	}

	err = connectionObject.Set("unbindQueue", func(call sobek.FunctionCall) sobek.Value {
		var queueConfig *QueueUnbindOptions
		if len(call.Arguments) == 0 {
			common.Throw(runtime, ErrNotEnoughArguments)
		}

		if params, ok := call.Argument(0).Export().(map[string]interface{}); ok {
			if b, errInt := json.Marshal(params); errInt != nil {
				common.Throw(runtime, errInt)
			} else {
				if errInt = json.Unmarshal(b, &queueConfig); errInt != nil {
					common.Throw(runtime, errInt)
				}
			}
		}

		a.unbindQueue(connection, queueConfig)
		return sobek.Undefined()
	})
	if err != nil {
		common.Throw(runtime, err)
	}

	err = connectionObject.Set("deleteQueue", func(call sobek.FunctionCall) sobek.Value {
		var queueConfig *QueueDeleteOptions
		if len(call.Arguments) == 0 {
			common.Throw(runtime, ErrNotEnoughArguments)
		}

		if params, ok := call.Argument(0).Export().(map[string]interface{}); ok {
			if b, errInt := json.Marshal(params); errInt != nil {
				common.Throw(runtime, errInt)
			} else {
				if errInt = json.Unmarshal(b, &queueConfig); errInt != nil {
					common.Throw(runtime, errInt)
				}
			}
		}

		a.deleteQueue(connection, queueConfig)
		return sobek.Undefined()
	})
	if err != nil {
		common.Throw(runtime, err)
	}

	err = connectionObject.Set("inspectQueue", func(call sobek.FunctionCall) sobek.Value {
		var queueConfig *QueueInspectOptions
		if len(call.Arguments) == 0 {
			common.Throw(runtime, ErrNotEnoughArguments)
		}

		if params, ok := call.Argument(0).Export().(map[string]interface{}); ok {
			if b, errInt := json.Marshal(params); errInt != nil {
				common.Throw(runtime, errInt)
			} else {
				if errInt = json.Unmarshal(b, &queueConfig); errInt != nil {
					common.Throw(runtime, errInt)
				}
			}
		}

		return runtime.ToValue(a.inspectQueue(connection, queueConfig))
	})
	if err != nil {
		common.Throw(runtime, err)
	}

	err = connectionObject.Set("purgeQueue", func(call sobek.FunctionCall) sobek.Value {
		var queue *QueuePurgeOptions
		if len(call.Arguments) == 0 {
			common.Throw(runtime, ErrNotEnoughArguments)
		}

		if params, ok := call.Argument(0).Export().(map[string]interface{}); ok {
			if b, errInt := json.Marshal(params); errInt != nil {
				common.Throw(runtime, errInt)
			} else {
				if errInt = json.Unmarshal(b, &queue); errInt != nil {
					common.Throw(runtime, errInt)
				}
			}
		}

		return runtime.ToValue(a.purgeQueue(connection, queue))
	})
	if err != nil {
		common.Throw(runtime, err)
	}

	err = connectionObject.Set("close", func(_ sobek.FunctionCall) sobek.Value {
		if errInt := connection.Close(); errInt != nil {
			common.Throw(runtime, errInt)
		}

		return sobek.Undefined()
	})
	if err != nil {
		common.Throw(runtime, err)
	}

	return connectionObject
}

// getAMQPConnection returns an amqp091-go connection for the given ConnectionURL.
func (a *AMQP) getAMQPConnection(config *ConnectionConfig) *amqpDriver.Connection {
	ctx := a.vu.Context()
	if ctx == nil {
		err := ErrEmptyContext
		common.Throw(a.vu.Runtime(), err)
		return nil
	}

	conn, err := amqpDriver.Dial(config.ConnectionURL)
	if err != nil {
		if conn == nil {
			common.Throw(a.vu.Runtime(), err)
			return nil
		}
	}

	return conn
}
