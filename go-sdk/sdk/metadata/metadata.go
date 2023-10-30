package metadata

import (
	"github.com/spf13/viper"
)

type Metadata struct {
}

func NewMetadata() *Metadata {
	return &Metadata{}
}

func (md Metadata) GetProduct() string {
	return viper.GetString("metadata.product_id")
}

func (md Metadata) GetWorkflow() string {
	return viper.GetString("metadata.workflow_id")
}

func (md Metadata) GetProcess() string {
	return viper.GetString("metadata.process_id")
}

func (md Metadata) GetVersion() string {
	return viper.GetString("metadata.version_id")
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
