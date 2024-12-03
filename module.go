// Package amqp is an extension for K6 to work with RabbitMQ.
package amqp

import (
	"github.com/grafana/sobek"
	"github.com/sirupsen/logrus"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
)

const version = "v0.4.3"

// logger is globally used by the AMQP module.
var logger *logrus.Logger //nolint:gochecknoglobals

// init registers the xk6-amqp-custom module as 'k6/x/amqp'.
func init() {
	// Initialize the global logger.
	logger = logrus.New()

	// Register the module namespace (aka. JS import path).
	modules.Register("k6/x/amqp", New())
}

type (
	// AMQP is a main struct for xk6-amqp extension.
	AMQP struct {
		vu      modules.VU
		metrics Metrics
		exports *sobek.Object
		version string
	}

	// RootModule implements Module interface.
	RootModule struct{}

	// Module implements Instance interface.
	Module struct {
		*AMQP
	}
)

var (
	_ modules.Instance = &Module{}
	_ modules.Module   = &RootModule{}
)

// New creates a new instance of the root module.
func New() *RootModule {
	return &RootModule{}
}

// NewModuleInstance creates a new instance of the AMQP module.
func (*RootModule) NewModuleInstance(virtualUser modules.VU) modules.Instance {
	runtime := virtualUser.Runtime()

	metrics, err := registerMetrics(virtualUser)
	if err != nil {
		common.Throw(virtualUser.Runtime(), err)
	}

	// Create a new AMQP module.
	moduleInstance := &Module{
		AMQP: &AMQP{
			vu:      virtualUser,
			metrics: metrics,
			exports: runtime.NewObject(),
			version: version,
		},
	}

	// Export constants to the JS code.
	// moduleInstance.defineConstants()

	mustExport := func(name string, value interface{}) {
		if err := moduleInstance.exports.Set(name, value); err != nil {
			common.Throw(runtime, err)
		}
	}

	// Export the constructors and functions from the AMQP module to the JS code.
	// The Publisher is a constructor and must be called with new, e.g. new Publisher(...).
	mustExport("Publisher", moduleInstance.publisherClass)
	// The Consumer is a constructor and must be called with new, e.g. new Consumer(...).
	mustExport("Consumer", moduleInstance.consumerClass)
	// The Connection is a constructor and must be called with new, e.g. new Connection(...).
	mustExport("Connection", moduleInstance.connectionClass)

	return moduleInstance
}

// Exports returns the exports of the AMQP module, which are the functions
// that can be called from the JS code.
func (m *Module) Exports() modules.Exports {
	return modules.Exports{
		Default: m.exports,
	}
}
