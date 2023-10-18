package metadata_test

import (
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"

	"github.com/konstellation-io/kai-sdk/go-sdk/sdk/metadata"
)

type SdkMetadataTestSuite struct {
	suite.Suite
	logger logr.Logger
}

func (s *SdkMetadataTestSuite) SetupSuite() {
	s.logger = testr.NewWithOptions(s.T(), testr.Options{Verbosity: 1, LogTimestamp: true})
}

func (s *SdkMetadataTestSuite) SetupTest() {
	// Reset viper values before each test
	viper.Reset()
}

func (s *SdkMetadataTestSuite) TestMetadata_GetMetadata_ExpectOK() {
	// Given
	viper.SetDefault("metadata.product_id", "product-name")
	viper.SetDefault("metadata.workflow_id", "workflow-name")
	viper.SetDefault("metadata.process_id", "process-name")
	viper.SetDefault("metadata.version_id", "version-name")
	viper.SetDefault("nats.object_store", "my-object-store")
	viper.SetDefault("centralized_configuration.global.bucket", "global-bucket")
	viper.SetDefault("centralized_configuration.product.bucket", "product-bucket")
	viper.SetDefault("centralized_configuration.workflow.bucket", "workflow-bucket")
	viper.SetDefault("centralized_configuration.process.bucket", "process-bucket")

	// When
	sdkMetadata := metadata.NewMetadata(s.logger)

	productName := sdkMetadata.GetProduct()
	workflowName := sdkMetadata.GetWorkflow()
	processName := sdkMetadata.GetProcess()
	versionName := sdkMetadata.GetVersion()
	objectStoreName := sdkMetadata.GetObjectStoreName()
	globalKvName := sdkMetadata.GetGlobalCentralizedConfigurationName()
	productKvName := sdkMetadata.GetProductCentralizedConfigurationName()
	workflowKvName := sdkMetadata.GetWorkflowCentralizedConfigurationName()
	processKvName := sdkMetadata.GetProcessCentralizedConfigurationName()

	// Then
	s.NotNil(sdkMetadata)
	s.Equal("product-name", productName)
	s.Equal("workflow-name", workflowName)
	s.Equal("process-name", processName)
	s.Equal("version-name", versionName)
	s.Equal("my-object-store", objectStoreName)
	s.Equal("global-bucket", globalKvName)
	s.Equal("product-bucket", productKvName)
	s.Equal("workflow-bucket", workflowKvName)
	s.Equal("process-bucket", processKvName)
}

func TestSdkMetadataTestSuite(t *testing.T) {
	suite.Run(t, new(SdkMetadataTestSuite))
}
