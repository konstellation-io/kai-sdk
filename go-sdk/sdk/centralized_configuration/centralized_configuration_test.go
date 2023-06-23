package centralized_configuration_test

import (
	"errors"
	"github.com/go-logr/logr/testr"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk/messaging"
	"testing"

	"github.com/go-logr/logr"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/mocks"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk/centralized_configuration"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

//go:generate mockery --dir $GOPATH/pkg/mod/github.com/nats-io/nats.go@v1.26.0 --output ../../mocks --name KeyValue --structname KeyValueMock --filename key_value_mock.go
//go:generate mockery --dir $GOPATH/pkg/mod/github.com/nats-io/nats.go@v1.26.0 --output ../../mocks --name KeyValueEntry --structname KeyValueEntryMock --filename key_value_entry_mock.go
//go:generate mockery --dir $GOPATH/pkg/mod/github.com/nats-io/nats.go@v1.26.0 --output ../../mocks --name JetStreamContext --structname JetStreamContextMock --filename jetstream_context_mock.go
type SdkCentralizedConfigurationTestSuite struct {
	suite.Suite
	logger     logr.Logger
	jetstream  mocks.JetStreamContextMock
	productKv  mocks.KeyValueMock
	workflowKv mocks.KeyValueMock
	processKv  mocks.KeyValueMock
}

func (suite *SdkCentralizedConfigurationTestSuite) SetupTest() {
	// Reset viper values before each test
	viper.Reset()

	suite.logger = testr.NewWithOptions(suite.T(), testr.Options{Verbosity: 1})
	suite.jetstream = *mocks.NewJetStreamContextMock(suite.T())
	suite.productKv = *mocks.NewKeyValueMock(suite.T())
	suite.workflowKv = *mocks.NewKeyValueMock(suite.T())
	suite.processKv = *mocks.NewKeyValueMock(suite.T())
}

func (suite *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_InitializeConfigurationScopes_ExpectOK() {
	// Given
	viper.SetDefault("centralized_configuration.product.bucket", "product-bucket")
	viper.SetDefault("centralized_configuration.workflow.bucket", "workflow-bucket")
	viper.SetDefault("centralized_configuration.process.bucket", "process-bucket")

	suite.jetstream.On("KeyValue", "product-bucket").Return(&suite.productKv, nil)
	suite.jetstream.On("KeyValue", "workflow-bucket").Return(&suite.workflowKv, nil)
	suite.jetstream.On("KeyValue", "process-bucket").Return(&suite.processKv, nil)

	// When
	conf, err := centralized_configuration.NewCentralizedConfiguration(suite.logger, &suite.jetstream)

	// Then
	suite.NoError(err)
	suite.NotNil(conf)
}

func (suite *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_InitializeConfigurationScopes_ProductConfigNotExist_ExpectError() {
	// Given
	viper.SetDefault("centralized_configuration.product.bucket", "some-product-bucket")
	viper.SetDefault("centralized_configuration.workflow.bucket", "some-workflow-bucket")
	viper.SetDefault("centralized_configuration.process.bucket", "some-process-bucket")

	suite.jetstream.On("KeyValue", "some-product-bucket").Return(nil, errors.New("not exist"))
	suite.jetstream.On("KeyValue", "some-workflow-bucket").Return(&suite.workflowKv, nil)
	suite.jetstream.On("KeyValue", "some-process-bucket").Return(&suite.processKv, nil)
	// When
	config, err := centralized_configuration.NewCentralizedConfiguration(suite.logger, &suite.jetstream)

	// Then
	suite.Error(err)
	suite.Nil(config)
}

func (suite *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_InitializeConfigurationScopes_WorkflowConfigNotExist_ExpectError() {
	// Given
	viper.SetDefault("centralized_configuration.product.bucket", "some-product-bucket")
	viper.SetDefault("centralized_configuration.workflow.bucket", "some-workflow-bucket")
	viper.SetDefault("centralized_configuration.process.bucket", "some-process-bucket")

	suite.jetstream.On("KeyValue", "some-product-bucket").Return(&suite.productKv, nil)
	suite.jetstream.On("KeyValue", "some-workflow-bucket").Return(nil, errors.New("not exist"))
	suite.jetstream.On("KeyValue", "some-process-bucket").Return(&suite.processKv, nil)

	// When
	config, err := centralized_configuration.NewCentralizedConfiguration(suite.logger, &suite.jetstream)

	// Then
	suite.Error(err)
	suite.Nil(config)
}

func (suite *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_InitializeConfigurationScopes_ProcessConfigNotExist_ExpectError() {
	// Given
	viper.SetDefault("centralized_configuration.product.bucket", "some-product-bucket")
	viper.SetDefault("centralized_configuration.workflow.bucket", "some-workflow-bucket")
	viper.SetDefault("centralized_configuration.process.bucket", "some-process-bucket")

	suite.jetstream.On("KeyValue", "some-product-bucket").Return(&suite.productKv, nil)
	suite.jetstream.On("KeyValue", "some-workflow-bucket").Return(&suite.workflowKv, nil)
	suite.jetstream.On("KeyValue", "some-process-bucket").Return(nil, errors.New("not exist"))

	// When
	config, err := centralized_configuration.NewCentralizedConfiguration(suite.logger, &suite.jetstream)

	// Then
	suite.Error(err)
	suite.Nil(config)
}

func (suite *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_DeleteConfigOnProductScope_ExpectOK() {
	// Given
	suite.productKv.On("Delete", "key1").Return(nil)

	config, err := centralized_configuration.NewCentralizedConfigurationBuilder(
		suite.logger,
		&suite.productKv,
		&suite.workflowKv,
		&suite.processKv,
	)

	// When
	err = config.DeleteConfig("key1", messaging.ProductScope)

	// Then
	suite.NotNil(config)
	suite.NoError(err)
	suite.productKv.AssertNumberOfCalls(suite.T(), "Delete", 1)
}

func TestSdkCentralizedConfigurationTestSuite(t *testing.T) {
	suite.Run(t, new(SdkCentralizedConfigurationTestSuite))
}
