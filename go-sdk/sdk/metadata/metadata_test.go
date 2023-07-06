package metadata_test

import (
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk/metadata"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
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
	viper.SetDefault("nats.object-store", "my-object-store")
	viper.SetDefault("centralized-configuration.product.bucket", "product-bucket")
	viper.SetDefault("centralized-configuration.workflow.bucket", "workflow-bucket")
	viper.SetDefault("centralized-configuration.process.bucket", "process-bucket")

	// When
	sdkMetadata := metadata.NewMetadata(s.logger)

	productName := sdkMetadata.GetProduct()
	workflowName := sdkMetadata.GetWorkflow()
	processName := sdkMetadata.GetProcess()
	versionName := sdkMetadata.GetVersion()
	objectStoreName := sdkMetadata.GetObjectStoreName()
	productKvName := sdkMetadata.GetKeyValueStoreProductName()
	workflowKvName := sdkMetadata.GetKeyValueStoreWorkflowName()
	processKvName := sdkMetadata.GetKeyValueStoreProcessName()

	// Then
	s.NotNil(sdkMetadata)
	s.Equal("product-name", productName)
	s.Equal("workflow-name", workflowName)
	s.Equal("process-name", processName)
	s.Equal("version-name", versionName)
	s.Equal("my-object-store", objectStoreName)
	s.Equal("product-bucket", productKvName)
	s.Equal("workflow-bucket", workflowKvName)
	s.Equal("process-bucket", processKvName)
}

func TestSdkMetadataTestSuite(t *testing.T) {
	suite.Run(t, new(SdkMetadataTestSuite))
}
