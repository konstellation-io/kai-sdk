package metadata_test

import (
	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk/metadata"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SdkMetadataTestSuite struct {
	suite.Suite
	logger logr.Logger
}

func (suite *SdkMetadataTestSuite) SetupSuite() {
	suite.logger = testr.NewWithOptions(suite.T(), testr.Options{Verbosity: 1})
}

func (suite *SdkMetadataTestSuite) SetupTest() {
	// Reset viper values before each test
	viper.Reset()
}

func (suite *SdkMetadataTestSuite) TestMetadata_GetMetadata_ExpectOK() {
	// Given
	viper.SetDefault("metadata.product_id", "product-name")
	viper.SetDefault("metadata.workflow_id", "workflow-name")
	viper.SetDefault("metadata.process_id", "process-name")
	viper.SetDefault("metadata.version_id", "version-name")
	viper.SetDefault("nats.object_store", "my-object-store")
	viper.SetDefault("centralized_configuration.product.bucket", "product-bucket")
	viper.SetDefault("centralized_configuration.workflow.bucket", "workflow-bucket")
	viper.SetDefault("centralized_configuration.process.bucket", "process-bucket")

	// When
	sdkMetadata := metadata.NewMetadata(suite.logger)

	productName := sdkMetadata.GetProduct()
	workflowName := sdkMetadata.GetWorkflow()
	processName := sdkMetadata.GetProcess()
	versionName := sdkMetadata.GetVersion()
	objectStoreName := sdkMetadata.GetObjectStoreName()
	productKvName := sdkMetadata.GetKeyValueStoreProductName()
	workflowKvName := sdkMetadata.GetKeyValueStoreWorkflowName()
	processKvName := sdkMetadata.GetKeyValueStoreProcessName()

	// Then
	suite.NotNil(sdkMetadata)
	suite.Equal("product-name", productName)
	suite.Equal("workflow-name", workflowName)
	suite.Equal("process-name", processName)
	suite.Equal("version-name", versionName)
	suite.Equal("my-object-store", objectStoreName)
	suite.Equal("product-bucket", productKvName)
	suite.Equal("workflow-bucket", workflowKvName)
	suite.Equal("process-bucket", processKvName)
}

func TestSdkMetadataTestSuite(t *testing.T) {
	suite.Run(t, new(SdkMetadataTestSuite))
}
