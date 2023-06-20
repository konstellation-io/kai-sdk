package centralized_configuration_test

import (
	"errors"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/mocks"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk/centralized_configuration"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk/messaging"
	"github.com/nats-io/nats.go"
)

func (suite *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_GetConfigOnAllScopes_ExpectOK() {
	// Given
	productResponse := mocks.NewKeyValueEntryMock(suite.T())
	productResponse.On("Value").Return([]byte("value1"))
	workflowResponse := mocks.NewKeyValueEntryMock(suite.T())
	workflowResponse.On("Value").Return([]byte("value2"))
	processResponse := mocks.NewKeyValueEntryMock(suite.T())
	processResponse.On("Value").Return([]byte("value3"))

	suite.productKv.On("Get", "key1").Return(productResponse, nil)
	suite.workflowKv.On("Get", "key1").Return(nil, nats.ErrKeyNotFound)
	suite.processKv.On("Get", "key1").Return(nil, nats.ErrKeyNotFound)

	suite.productKv.On("Get", "key2").Return(nil, nats.ErrKeyNotFound)
	suite.workflowKv.On("Get", "key2").Return(workflowResponse, nil)
	suite.processKv.On("Get", "key2").Return(nil, nats.ErrKeyNotFound)

	suite.productKv.On("Get", "key3").Return(nil, nats.ErrKeyNotFound)
	suite.workflowKv.On("Get", "key3").Return(nil, nats.ErrKeyNotFound)
	suite.processKv.On("Get", "key3").Return(processResponse, nil)
	config, err := centralized_configuration.NewCentralizedConfigurationBuilder(
		suite.logger,
		&suite.productKv,
		&suite.workflowKv,
		&suite.processKv,
	)

	// When
	key1Value, err := config.GetConfig("key1")
	key2Value, err := config.GetConfig("key2")
	key3Value, err := config.GetConfig("key3")

	// Then
	suite.NotNil(config)
	suite.NoError(err)
	suite.productKv.AssertNumberOfCalls(suite.T(), "Get", 1)
	suite.Equal("value1", key1Value)
	suite.workflowKv.AssertNumberOfCalls(suite.T(), "Get", 2)
	suite.Equal("value2", key2Value)
	suite.processKv.AssertNumberOfCalls(suite.T(), "Get", 3)
	suite.Equal("value3", key3Value)
}

func (suite *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_GetNonExistingConfigOnAllScopes_ExpectOK() {
	// Given
	suite.productKv.On("Get", "key1").Return(nil, nats.ErrKeyNotFound)
	suite.workflowKv.On("Get", "key1").Return(nil, nats.ErrKeyNotFound)
	suite.processKv.On("Get", "key1").Return(nil, nats.ErrKeyNotFound)

	config, err := centralized_configuration.NewCentralizedConfigurationBuilder(
		suite.logger,
		&suite.productKv,
		&suite.workflowKv,
		&suite.processKv,
	)

	// When
	key1Value, err := config.GetConfig("key1")

	// Then
	suite.NotNil(config)
	suite.Empty(key1Value)
	suite.Error(err)
	suite.productKv.AssertNumberOfCalls(suite.T(), "Get", 1)
	suite.workflowKv.AssertNumberOfCalls(suite.T(), "Get", 1)
	suite.processKv.AssertNumberOfCalls(suite.T(), "Get", 1)
}

func (suite *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_GetNonExistingConfigOnAllScopes_UnexpectedError_ExpectOK() {
	// Given
	suite.productKv.On("Get", "key1").Return(nil, nats.ErrKeyNotFound)
	suite.workflowKv.On("Get", "key1").Return(nil, errors.New("unexpected error"))
	suite.processKv.On("Get", "key1").Return(nil, nats.ErrKeyNotFound)

	config, err := centralized_configuration.NewCentralizedConfigurationBuilder(
		suite.logger,
		&suite.productKv,
		&suite.workflowKv,
		&suite.processKv,
	)

	// When
	key1Value, err := config.GetConfig("key1")

	// Then
	suite.NotNil(config)
	suite.Empty(key1Value)
	suite.Error(err)
	suite.productKv.AssertNumberOfCalls(suite.T(), "Get", 0)
	suite.workflowKv.AssertNumberOfCalls(suite.T(), "Get", 1)
	suite.processKv.AssertNumberOfCalls(suite.T(), "Get", 1)
}

func (suite *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_GetOverridenConfigOnAllScopes_ExpectOK() {
	// Given
	productResponse := mocks.NewKeyValueEntryMock(suite.T())
	workflowResponse := mocks.NewKeyValueEntryMock(suite.T())
	workflowResponse.On("Value").Return([]byte("workflowValue"))
	processResponse := mocks.NewKeyValueEntryMock(suite.T())
	processResponse.On("Value").Return([]byte("processValue"))

	suite.productKv.On("Get", "key1").Return(productResponse, nil)
	suite.workflowKv.On("Get", "key1").Return(workflowResponse, nil)
	suite.processKv.On("Get", "key1").Return(nil, nats.ErrKeyNotFound)

	suite.productKv.On("Get", "key2").Return(productResponse, nil)
	suite.workflowKv.On("Get", "key2").Return(workflowResponse, nil)
	suite.processKv.On("Get", "key2").Return(processResponse, nil)

	config, err := centralized_configuration.NewCentralizedConfigurationBuilder(
		suite.logger,
		&suite.productKv,
		&suite.workflowKv,
		&suite.processKv,
	)

	// When
	key1Value, err := config.GetConfig("key1")
	key2Value, err := config.GetConfig("key2")

	// Then
	suite.NotNil(config)
	suite.NoError(err)
	suite.productKv.AssertNotCalled(suite.T(), "Get")
	suite.workflowKv.AssertNumberOfCalls(suite.T(), "Get", 1)
	suite.processKv.AssertNumberOfCalls(suite.T(), "Get", 2)
	suite.Equal("workflowValue", key1Value)
	suite.Equal("processValue", key2Value)
}

func (suite *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_GetConfigOnProductScope_ExpectOK() {
	// Given
	productResponse := mocks.NewKeyValueEntryMock(suite.T())
	productResponse.On("Value").Return([]byte("value1"))

	suite.productKv.On("Get", "key1").Return(productResponse, nil)

	config, err := centralized_configuration.NewCentralizedConfigurationBuilder(
		suite.logger,
		&suite.productKv,
		&suite.workflowKv,
		&suite.processKv,
	)

	// When
	key1Value, err := config.GetConfig("key1", messaging.ProductScope)

	// Then
	suite.NotNil(config)
	suite.NoError(err)
	suite.productKv.AssertNumberOfCalls(suite.T(), "Get", 1)
	suite.Equal("value1", key1Value)
}

func (suite *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_GetNonExistingConfigOnProductScope_ExpectOK() {
	// Given
	suite.productKv.On("Get", "key1").Return(nil, nats.ErrKeyNotFound)

	config, err := centralized_configuration.NewCentralizedConfigurationBuilder(
		suite.logger,
		&suite.productKv,
		&suite.workflowKv,
		&suite.processKv,
	)

	// When
	key1Value, err := config.GetConfig("key1", messaging.ProductScope)

	// Then
	suite.NotNil(config)
	suite.Empty(key1Value)
	suite.Error(err)
	suite.productKv.AssertNumberOfCalls(suite.T(), "Get", 1)
}

func (suite *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_GetConfigOnWorkflowScope_ExpectOK() {
	// Given
	workflowResponse := mocks.NewKeyValueEntryMock(suite.T())
	workflowResponse.On("Value").Return([]byte("value2"))

	suite.workflowKv.On("Get", "key1").Return(workflowResponse, nil)

	config, err := centralized_configuration.NewCentralizedConfigurationBuilder(
		suite.logger,
		&suite.productKv,
		&suite.workflowKv,
		&suite.processKv,
	)

	// When
	key1Value, err := config.GetConfig("key1", messaging.WorkflowScope)

	// Then
	suite.NotNil(config)
	suite.NoError(err)
	suite.workflowKv.AssertNumberOfCalls(suite.T(), "Get", 1)
	suite.Equal("value2", key1Value)
}

func (suite *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_GetNonExistingConfigOnWorkflowScope_ExpectOK() {
	// Given
	suite.workflowKv.On("Get", "key1").Return(nil, nats.ErrKeyNotFound)

	config, err := centralized_configuration.NewCentralizedConfigurationBuilder(
		suite.logger,
		&suite.productKv,
		&suite.workflowKv,
		&suite.processKv,
	)

	// When
	key1Value, err := config.GetConfig("key1", messaging.WorkflowScope)

	// Then
	suite.NotNil(config)
	suite.Empty(key1Value)
	suite.Error(err)
	suite.workflowKv.AssertNumberOfCalls(suite.T(), "Get", 1)
}

func (suite *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_GetConfigOnProcessScope_ExpectOK() {
	// Given
	processResponse := mocks.NewKeyValueEntryMock(suite.T())
	processResponse.On("Value").Return([]byte("value3"))

	suite.processKv.On("Get", "key1").Return(processResponse, nil)

	config, err := centralized_configuration.NewCentralizedConfigurationBuilder(
		suite.logger,
		&suite.productKv,
		&suite.workflowKv,
		&suite.processKv,
	)

	// When
	key1Value, err := config.GetConfig("key1", messaging.ProcessScope)

	// Then
	suite.NotNil(config)
	suite.NoError(err)
	suite.processKv.AssertNumberOfCalls(suite.T(), "Get", 1)
	suite.Equal("value3", key1Value)
}

func (suite *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_GetNonExistingConfigOnProcessScope_ExpectOK() {
	// Given
	suite.processKv.On("Get", "key1").Return(nil, nats.ErrKeyNotFound)

	config, err := centralized_configuration.NewCentralizedConfigurationBuilder(
		suite.logger,
		&suite.productKv,
		&suite.workflowKv,
		&suite.processKv,
	)

	// When
	key1Value, err := config.GetConfig("key1", messaging.ProcessScope)

	// Then
	suite.NotNil(config)
	suite.Empty(key1Value)
	suite.Error(err)
	suite.processKv.AssertNumberOfCalls(suite.T(), "Get", 1)
}
