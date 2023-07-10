package centralizedconfiguration_test

import (
	"errors"

	"github.com/konstellation-io/kre-runners/go-sdk/v1/mocks"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk/centralized-configuration"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk/messaging"
	"github.com/nats-io/nats.go"
)

func (s *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_GetConfigOnAllScopes_ExpectOK() {
	// Given
	productResponse := mocks.NewKeyValueEntryMock(s.T())
	productResponse.On("Value").Return([]byte("value1"))
	workflowResponse := mocks.NewKeyValueEntryMock(s.T())
	workflowResponse.On("Value").Return([]byte("value2"))
	processResponse := mocks.NewKeyValueEntryMock(s.T())
	processResponse.On("Value").Return([]byte("value3"))

	s.productKv.On("Get", "key1").Return(productResponse, nil)
	s.workflowKv.On("Get", "key1").Return(nil, nats.ErrKeyNotFound)
	s.processKv.On("Get", "key1").Return(nil, nats.ErrKeyNotFound)

	s.productKv.On("Get", "key2").Return(nil, nats.ErrKeyNotFound)
	s.workflowKv.On("Get", "key2").Return(workflowResponse, nil)
	s.processKv.On("Get", "key2").Return(nil, nats.ErrKeyNotFound)

	s.productKv.On("Get", "key3").Return(nil, nats.ErrKeyNotFound)
	s.workflowKv.On("Get", "key3").Return(nil, nats.ErrKeyNotFound)
	s.processKv.On("Get", "key3").Return(processResponse, nil)
	config, err := centralizedconfiguration.NewCentralizedConfigurationBuilder(
		s.logger,
		&s.productKv,
		&s.workflowKv,
		&s.processKv,
	)

	// When
	key1Value, err := config.GetConfig("key1")
	key2Value, err := config.GetConfig("key2")
	key3Value, err := config.GetConfig("key3")

	// Then
	s.NotNil(config)
	s.NoError(err)
	s.productKv.AssertNumberOfCalls(s.T(), "Get", 1)
	s.Equal("value1", key1Value)
	s.workflowKv.AssertNumberOfCalls(s.T(), "Get", 2)
	s.Equal("value2", key2Value)
	s.processKv.AssertNumberOfCalls(s.T(), "Get", 3)
	s.Equal("value3", key3Value)
}

func (s *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_GetNonExistingConfigOnAllScopes_ExpectError() {
	// Given
	s.productKv.On("Get", "key1").Return(nil, nats.ErrKeyNotFound)
	s.workflowKv.On("Get", "key1").Return(nil, nats.ErrKeyNotFound)
	s.processKv.On("Get", "key1").Return(nil, nats.ErrKeyNotFound)

	config, err := centralizedconfiguration.NewCentralizedConfigurationBuilder(
		s.logger,
		&s.productKv,
		&s.workflowKv,
		&s.processKv,
	)

	// When
	key1Value, err := config.GetConfig("key1")

	// Then
	s.NotNil(config)
	s.Empty(key1Value)
	s.Error(err)
	s.productKv.AssertNumberOfCalls(s.T(), "Get", 1)
	s.workflowKv.AssertNumberOfCalls(s.T(), "Get", 1)
	s.processKv.AssertNumberOfCalls(s.T(), "Get", 1)
}

func (s *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_GetNonExistingConfigOnAllScopes_UnexpectedError_ExpectError() {
	// Given
	s.productKv.On("Get", "key1").Return(nil, nats.ErrKeyNotFound)
	s.workflowKv.On("Get", "key1").Return(nil, errors.New("unexpected error"))
	s.processKv.On("Get", "key1").Return(nil, nats.ErrKeyNotFound)

	config, err := centralizedconfiguration.NewCentralizedConfigurationBuilder(
		s.logger,
		&s.productKv,
		&s.workflowKv,
		&s.processKv,
	)

	// When
	key1Value, err := config.GetConfig("key1")

	// Then
	s.NotNil(config)
	s.Empty(key1Value)
	s.Error(err)
	s.productKv.AssertNumberOfCalls(s.T(), "Get", 0)
	s.workflowKv.AssertNumberOfCalls(s.T(), "Get", 1)
	s.processKv.AssertNumberOfCalls(s.T(), "Get", 1)
}

func (s *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_GetOverridenConfigOnAllScopes_ExpectOK() {
	// Given
	productResponse := mocks.NewKeyValueEntryMock(s.T())
	workflowResponse := mocks.NewKeyValueEntryMock(s.T())
	workflowResponse.On("Value").Return([]byte("workflowValue"))
	processResponse := mocks.NewKeyValueEntryMock(s.T())
	processResponse.On("Value").Return([]byte("processValue"))

	s.productKv.On("Get", "key1").Return(productResponse, nil)
	s.workflowKv.On("Get", "key1").Return(workflowResponse, nil)
	s.processKv.On("Get", "key1").Return(nil, nats.ErrKeyNotFound)

	s.productKv.On("Get", "key2").Return(productResponse, nil)
	s.workflowKv.On("Get", "key2").Return(workflowResponse, nil)
	s.processKv.On("Get", "key2").Return(processResponse, nil)

	config, err := centralizedconfiguration.NewCentralizedConfigurationBuilder(
		s.logger,
		&s.productKv,
		&s.workflowKv,
		&s.processKv,
	)

	// When
	key1Value, err := config.GetConfig("key1")
	key2Value, err := config.GetConfig("key2")

	// Then
	s.NotNil(config)
	s.NoError(err)
	s.productKv.AssertNotCalled(s.T(), "Get")
	s.workflowKv.AssertNumberOfCalls(s.T(), "Get", 1)
	s.processKv.AssertNumberOfCalls(s.T(), "Get", 2)
	s.Equal("workflowValue", key1Value)
	s.Equal("processValue", key2Value)
}

func (s *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_GetConfigOnProductScope_ExpectOK() {
	// Given
	productResponse := mocks.NewKeyValueEntryMock(s.T())
	productResponse.On("Value").Return([]byte("value1"))

	s.productKv.On("Get", "key1").Return(productResponse, nil)

	config, err := centralizedconfiguration.NewCentralizedConfigurationBuilder(
		s.logger,
		&s.productKv,
		&s.workflowKv,
		&s.processKv,
	)

	// When
	key1Value, err := config.GetConfig("key1", messaging.ProductScope)

	// Then
	s.NotNil(config)
	s.NoError(err)
	s.productKv.AssertNumberOfCalls(s.T(), "Get", 1)
	s.Equal("value1", key1Value)
}

func (s *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_GetNonExistingConfigOnProductScope_ExpectError() {
	// Given
	s.productKv.On("Get", "key1").Return(nil, nats.ErrKeyNotFound)

	config, err := centralizedconfiguration.NewCentralizedConfigurationBuilder(
		s.logger,
		&s.productKv,
		&s.workflowKv,
		&s.processKv,
	)

	// When
	key1Value, err := config.GetConfig("key1", messaging.ProductScope)

	// Then
	s.NotNil(config)
	s.Empty(key1Value)
	s.Error(err)
	s.productKv.AssertNumberOfCalls(s.T(), "Get", 1)
}

func (s *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_GetConfigOnWorkflowScope_ExpectOK() {
	// Given
	workflowResponse := mocks.NewKeyValueEntryMock(s.T())
	workflowResponse.On("Value").Return([]byte("value2"))

	s.workflowKv.On("Get", "key1").Return(workflowResponse, nil)

	config, err := centralizedconfiguration.NewCentralizedConfigurationBuilder(
		s.logger,
		&s.productKv,
		&s.workflowKv,
		&s.processKv,
	)

	// When
	key1Value, err := config.GetConfig("key1", messaging.WorkflowScope)

	// Then
	s.NotNil(config)
	s.NoError(err)
	s.workflowKv.AssertNumberOfCalls(s.T(), "Get", 1)
	s.Equal("value2", key1Value)
}

func (s *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_GetNonExistingConfigOnWorkflowScope_ExpectError() {
	// Given
	s.workflowKv.On("Get", "key1").Return(nil, nats.ErrKeyNotFound)

	config, err := centralizedconfiguration.NewCentralizedConfigurationBuilder(
		s.logger,
		&s.productKv,
		&s.workflowKv,
		&s.processKv,
	)

	// When
	key1Value, err := config.GetConfig("key1", messaging.WorkflowScope)

	// Then
	s.NotNil(config)
	s.Empty(key1Value)
	s.Error(err)
	s.workflowKv.AssertNumberOfCalls(s.T(), "Get", 1)
}

func (s *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_GetConfigOnProcessScope_ExpectOK() {
	// Given
	processResponse := mocks.NewKeyValueEntryMock(s.T())
	processResponse.On("Value").Return([]byte("value3"))

	s.processKv.On("Get", "key1").Return(processResponse, nil)

	config, err := centralizedconfiguration.NewCentralizedConfigurationBuilder(
		s.logger,
		&s.productKv,
		&s.workflowKv,
		&s.processKv,
	)

	// When
	key1Value, err := config.GetConfig("key1", messaging.ProcessScope)

	// Then
	s.NotNil(config)
	s.NoError(err)
	s.processKv.AssertNumberOfCalls(s.T(), "Get", 1)
	s.Equal("value3", key1Value)
}

func (s *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_GetNonExistingConfigOnProcessScope_ExpectError() {
	// Given
	s.processKv.On("Get", "key1").Return(nil, nats.ErrKeyNotFound)

	config, err := centralizedconfiguration.NewCentralizedConfigurationBuilder(
		s.logger,
		&s.productKv,
		&s.workflowKv,
		&s.processKv,
	)

	// When
	key1Value, err := config.GetConfig("key1", messaging.ProcessScope)

	// Then
	s.NotNil(config)
	s.Empty(key1Value)
	s.Error(err)
	s.processKv.AssertNumberOfCalls(s.T(), "Get", 1)
}