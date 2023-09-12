package centralizedconfiguration

import (
	"github.com/go-logr/logr"
	"github.com/nats-io/nats.go"
)

func NewCentralizedConfigurationBuilder(logger logr.Logger,
	productKv, workflowKv, processKv nats.KeyValue) (*CentralizedConfiguration, error) {
	logger = logger.WithName("[CENTRALIZED CONFIGURATION]")

	return &CentralizedConfiguration{
		logger:     logger,
		productKv:  productKv,
		workflowKv: workflowKv,
		processKv:  processKv,
	}, nil
}
