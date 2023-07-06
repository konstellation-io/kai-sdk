package centralizedconfiguration_test

import (
	"errors"
	"testing"

	centralizedConfiguration "github.com/konstellation-io/kre-runners/go-sdk/v1/sdk/centralized-configuration"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/mocks"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk/messaging"
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

func (s *SdkCentralizedConfigurationTestSuite) SetupSuite() {
	s.logger = testr.NewWithOptions(s.T(), testr.Options{Verbosity: 1})
}

func (s *SdkCentralizedConfigurationTestSuite) SetupTest() {
	// Reset viper values before each test
	viper.Reset()

	s.jetstream = *mocks.NewJetStreamContextMock(s.T())
	s.productKv = *mocks.NewKeyValueMock(s.T())
	s.workflowKv = *mocks.NewKeyValueMock(s.T())
	s.processKv = *mocks.NewKeyValueMock(s.T())
}

func (s *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_InitializeConfigurationScopes_ExpectOK() {
	// Given
	viper.SetDefault("centralized-configuration.product.bucket", "product-bucket")
	viper.SetDefault("centralized-configuration.workflow.bucket", "workflow-bucket")
	viper.SetDefault("centralized-configuration.process.bucket", "process-bucket")

	s.jetstream.On("KeyValue", "product-bucket").Return(&s.productKv, nil)
	s.jetstream.On("KeyValue", "workflow-bucket").Return(&s.workflowKv, nil)
	s.jetstream.On("KeyValue", "process-bucket").Return(&s.processKv, nil)

	// When
	conf, err := centralizedConfiguration.NewCentralizedConfiguration(s.logger, &s.jetstream)

	// Then
	s.NoError(err)
	s.NotNil(conf)
}

func (s *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_InitializeConfigurationScopes_ProductConfigNotExist_ExpectError() {
	// Given
	viper.SetDefault("centralized-configuration.product.bucket", "some-product-bucket")
	viper.SetDefault("centralized-configuration.workflow.bucket", "some-workflow-bucket")
	viper.SetDefault("centralized-configuration.process.bucket", "some-process-bucket")

	s.jetstream.On("KeyValue", "some-product-bucket").Return(nil, errors.New("not exist"))
	s.jetstream.On("KeyValue", "some-workflow-bucket").Return(&s.workflowKv, nil)
	s.jetstream.On("KeyValue", "some-process-bucket").Return(&s.processKv, nil)
	// When
	config, err := centralizedConfiguration.NewCentralizedConfiguration(s.logger, &s.jetstream)

	// Then
	s.Error(err)
	s.Nil(config)
}

func (s *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_InitializeConfigurationScopes_WorkflowConfigNotExist_ExpectError() {
	// Given
	viper.SetDefault("centralized-configuration.product.bucket", "some-product-bucket")
	viper.SetDefault("centralized-configuration.workflow.bucket", "some-workflow-bucket")
	viper.SetDefault("centralized-configuration.process.bucket", "some-process-bucket")

	s.jetstream.On("KeyValue", "some-product-bucket").Return(&s.productKv, nil)
	s.jetstream.On("KeyValue", "some-workflow-bucket").Return(nil, errors.New("not exist"))
	s.jetstream.On("KeyValue", "some-process-bucket").Return(&s.processKv, nil)

	// When
	config, err := centralizedConfiguration.NewCentralizedConfiguration(s.logger, &s.jetstream)

	// Then
	s.Error(err)
	s.Nil(config)
}

func (s *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_InitializeConfigurationScopes_ProcessConfigNotExist_ExpectError() {
	// Given
	viper.SetDefault("centralized-configuration.product.bucket", "some-product-bucket")
	viper.SetDefault("centralized-configuration.workflow.bucket", "some-workflow-bucket")
	viper.SetDefault("centralized-configuration.process.bucket", "some-process-bucket")

	s.jetstream.On("KeyValue", "some-product-bucket").Return(&s.productKv, nil)
	s.jetstream.On("KeyValue", "some-workflow-bucket").Return(&s.workflowKv, nil)
	s.jetstream.On("KeyValue", "some-process-bucket").Return(nil, errors.New("not exist"))

	// When
	config, err := centralizedConfiguration.NewCentralizedConfiguration(s.logger, &s.jetstream)

	// Then
	s.Error(err)
	s.Nil(config)
}

func (s *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_DeleteConfigOnProductScope_ExpectOK() {
	// Given
	s.productKv.On("Delete", "key1").Return(nil)

	config, err := centralizedConfiguration.NewCentralizedConfigurationBuilder(
		s.logger,
		&s.productKv,
		&s.workflowKv,
		&s.processKv,
	)

	// When
	err = config.DeleteConfig("key1", messaging.ProductScope)

	// Then
	s.NotNil(config)
	s.NoError(err)
	s.productKv.AssertNumberOfCalls(s.T(), "Delete", 1)
}

func TestSdkCentralizedConfigurationTestSuite(t *testing.T) {
	suite.Run(t, new(SdkCentralizedConfigurationTestSuite))
}
