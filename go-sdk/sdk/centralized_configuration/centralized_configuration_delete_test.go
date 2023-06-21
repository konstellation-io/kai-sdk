package centralized_configuration_test

import (
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk/centralized_configuration"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk/messaging"
	"github.com/nats-io/nats.go"
)

func (suite *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_DeleteConfigOnWorkflowScope_ExpectOK() {
	// Given
	suite.workflowKv.On("Delete", "key1").Return(nil)

	config, err := centralized_configuration.NewCentralizedConfigurationBuilder(
		suite.logger,
		&suite.productKv,
		&suite.workflowKv,
		&suite.processKv,
	)

	// When
	err = config.DeleteConfig("key1", messaging.WorkflowScope)

	// Then
	suite.NotNil(config)
	suite.NoError(err)
	suite.workflowKv.AssertNumberOfCalls(suite.T(), "Delete", 1)
}

func (suite *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_DeleteConfigOnProcessScope_ExpectOK() {
	// Given
	suite.processKv.On("Delete", "key1").Return(nil)

	config, err := centralized_configuration.NewCentralizedConfigurationBuilder(
		suite.logger,
		&suite.productKv,
		&suite.workflowKv,
		&suite.processKv,
	)

	// When
	err = config.DeleteConfig("key1", messaging.ProcessScope)

	// Then
	suite.NotNil(config)
	suite.NoError(err)
	suite.processKv.AssertNumberOfCalls(suite.T(), "Delete", 1)
}

func (suite *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_DeleteNonExistingConfigOnAllScope_ExpectError() {
	// Given
	suite.productKv.On("Delete", "key1").Return(nats.ErrKeyNotFound)
	suite.workflowKv.On("Delete", "key1").Return(nats.ErrKeyNotFound)
	suite.processKv.On("Delete", "key1").Return(nats.ErrKeyNotFound)

	config, err := centralized_configuration.NewCentralizedConfigurationBuilder(
		suite.logger,
		&suite.productKv,
		&suite.workflowKv,
		&suite.processKv,
	)

	// When
	err = config.DeleteConfig("key1", messaging.ProductScope)
	err = config.DeleteConfig("key1", messaging.WorkflowScope)
	err = config.DeleteConfig("key1", messaging.ProcessScope)

	// Then
	suite.NotNil(config)
	suite.Error(err)
	suite.productKv.AssertNumberOfCalls(suite.T(), "Delete", 1)
	suite.workflowKv.AssertNumberOfCalls(suite.T(), "Delete", 1)
	suite.processKv.AssertNumberOfCalls(suite.T(), "Delete", 1)
}
