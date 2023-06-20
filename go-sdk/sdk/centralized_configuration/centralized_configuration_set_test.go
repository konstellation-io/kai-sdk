package centralized_configuration_test

import (
	"errors"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk/centralized_configuration"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk/messaging"
)

func (suite *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_SetConfigOnDefaultScope_ExpectOK() {
	// Given
	suite.productKv.On("PutString", "key1", "value1").Return(uint64(0), errors.New("not exist"))
	suite.workflowKv.On("PutString", "key1", "value1").Return(uint64(0), errors.New("not exist"))
	suite.processKv.On("PutString", "key1", "value1").Return(uint64(1), nil)

	config, err := centralized_configuration.NewCentralizedConfigurationBuilder(
		suite.logger,
		&suite.productKv,
		&suite.workflowKv,
		&suite.processKv,
	)

	// When
	err = config.SetConfig("key1", "value1")

	// Then
	suite.NotNil(config)
	suite.NoError(err)
	suite.productKv.AssertNotCalled(suite.T(), "PutString")
	suite.workflowKv.AssertNotCalled(suite.T(), "PutString")
	suite.processKv.AssertNumberOfCalls(suite.T(), "PutString", 1)
}

func (suite *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_SetConfigOnProductScope_ExpectOK() {
	// Given
	suite.productKv.On("PutString", "key1", "value1").Return(uint64(1), nil)
	suite.workflowKv.On("PutString", "key1", "value1").Return(uint64(0), errors.New("not exist"))
	suite.processKv.On("PutString", "key1", "value1").Return(uint64(0), errors.New("not exist"))

	config, err := centralized_configuration.NewCentralizedConfigurationBuilder(
		suite.logger,
		&suite.productKv,
		&suite.workflowKv,
		&suite.processKv,
	)

	// When
	err = config.SetConfig("key1", "value1", messaging.ProductScope)

	// Then
	suite.NotNil(config)
	suite.NoError(err)
	suite.productKv.AssertNumberOfCalls(suite.T(), "PutString", 1)
	suite.productKv.AssertCalled(suite.T(), "PutString", "key1", "value1")
	suite.workflowKv.AssertNotCalled(suite.T(), "PutString")
	suite.processKv.AssertNotCalled(suite.T(), "PutString")
}

func (suite *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_SetConfigOnWorkflowScope_ExpectOK() {
	// Given
	suite.productKv.On("PutString", "key1", "value1").Return(uint64(0), errors.New("not exist"))
	suite.workflowKv.On("PutString", "key1", "value1").Return(uint64(1), nil)
	suite.processKv.On("PutString", "key1", "value1").Return(uint64(0), errors.New("not exist"))

	config, err := centralized_configuration.NewCentralizedConfigurationBuilder(
		suite.logger,
		&suite.productKv,
		&suite.workflowKv,
		&suite.processKv,
	)

	// When
	err = config.SetConfig("key1", "value1", messaging.WorkflowScope)

	// Then
	suite.NotNil(config)
	suite.NoError(err)
	suite.productKv.AssertNotCalled(suite.T(), "PutString")
	suite.workflowKv.AssertNumberOfCalls(suite.T(), "PutString", 1)
	suite.workflowKv.AssertCalled(suite.T(), "PutString", "key1", "value1")
	suite.processKv.AssertNotCalled(suite.T(), "PutString")
}

func (suite *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_SetConfigOnProcessScope_ExpectOK() {
	// Given
	suite.productKv.On("PutString", "key1", "value1").Return(uint64(0), errors.New("not exist"))
	suite.workflowKv.On("PutString", "key1", "value1").Return(uint64(0), errors.New("not exist"))
	suite.processKv.On("PutString", "key1", "value1").Return(uint64(1), nil)
	config, err := centralized_configuration.NewCentralizedConfigurationBuilder(
		suite.logger,
		&suite.productKv,
		&suite.workflowKv,
		&suite.processKv,
	)

	// When
	err = config.SetConfig("key1", "value1", messaging.ProcessScope)

	// Then
	suite.NotNil(config)
	suite.NoError(err)
	suite.productKv.AssertNotCalled(suite.T(), "PutString")
	suite.workflowKv.AssertNotCalled(suite.T(), "PutString")
	suite.processKv.AssertNumberOfCalls(suite.T(), "PutString", 1)
	suite.processKv.AssertCalled(suite.T(), "PutString", "key1", "value1")
}

func (suite *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_SetConfigOnDefaultScope_UnexpectedError_ExpectOK() {
	// Given
	suite.productKv.On("PutString", "key1", "value1").Return(uint64(0), errors.New("not exist"))
	suite.workflowKv.On("PutString", "key1", "value1").Return(uint64(0), errors.New("not exist"))
	suite.processKv.On("PutString", "key1", "value1").Return(uint64(0), errors.New("not exist"))

	config, err := centralized_configuration.NewCentralizedConfigurationBuilder(
		suite.logger,
		&suite.productKv,
		&suite.workflowKv,
		&suite.processKv,
	)

	// When
	err = config.SetConfig("key1", "value1")

	// Then
	suite.NotNil(config)
	suite.Error(err)
	suite.productKv.AssertNotCalled(suite.T(), "PutString")
	suite.workflowKv.AssertNotCalled(suite.T(), "PutString")
	suite.processKv.AssertNumberOfCalls(suite.T(), "PutString", 1)
}
