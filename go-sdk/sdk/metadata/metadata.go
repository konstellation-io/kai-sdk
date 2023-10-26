package metadata

import (
	"github.com/go-logr/logr"
	"github.com/spf13/viper"
)

type Metadata struct {
	logger logr.Logger
}

func NewMetadata(logger logr.Logger) *Metadata {
	return &Metadata{
		logger: logger.WithName("[METADATA]"),
	}
}

func (md Metadata) GetProduct() string {
	return viper.GetString("metadata.product_id")
}

func (md Metadata) GetWorkflow() string {
	return viper.GetString("metadata.workflow_name")
}

func (md Metadata) GetProcess() string {
	return viper.GetString("metadata.process_name")
}

func (md Metadata) GetVersion() string {
	return viper.GetString("metadata.version_tag")
}

func (md Metadata) GetEphemeralStorageName() string {
	return viper.GetString("nats.object_store")
}

func (md Metadata) GetGlobalCentralizedConfigurationName() string {
	return viper.GetString("centralized_configuration.global.bucket")
}

func (md Metadata) GetProductCentralizedConfigurationName() string {
	return viper.GetString("centralized_configuration.product.bucket")
}

func (md Metadata) GetWorkflowCentralizedConfigurationName() string {
	return viper.GetString("centralized_configuration.workflow.bucket")
}

func (md Metadata) GetProcessCentralizedConfigurationName() string {
	return viper.GetString("centralized_configuration.process.bucket")
}
