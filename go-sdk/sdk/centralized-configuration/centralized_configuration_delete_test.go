package centralizedconfiguration_test

import (
	centralizedConfiguration "github.com/konstellation-io/kai-sdk/go-sdk/sdk/centralized-configuration"
	"github.com/konstellation-io/kai-sdk/go-sdk/sdk/messaging"
	"github.com/nats-io/nats.go"
)

func (s *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_DeleteConfigOnWorkflowScope_ExpectOK() {
	// Given
	s.workflowKv.On("Delete", "key1").Return(nil)

	config, err := centralizedConfiguration.NewCentralizedConfigurationBuilder(
		s.logger,
		&s.productKv,
		&s.workflowKv,
		&s.processKv,
	)
	s.Require().NoError(err)

	// When
	err = config.DeleteConfig("key1", messaging.WorkflowScope)
	s.Require().NoError(err)

	// Then
	s.NotNil(config)
	s.workflowKv.AssertNumberOfCalls(s.T(), "Delete", 1)
}

func (s *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_DeleteConfigOnProcessScope_ExpectOK() {
	// Given
	s.processKv.On("Delete", "key1").Return(nil)

	config, err := centralizedConfiguration.NewCentralizedConfigurationBuilder(
		s.logger,
		&s.productKv,
		&s.workflowKv,
		&s.processKv,
	)
	s.Require().NoError(err)

	// When
	err = config.DeleteConfig("key1", messaging.ProcessScope)
	s.Require().NoError(err)

	// Then
	s.NotNil(config)
	s.processKv.AssertNumberOfCalls(s.T(), "Delete", 1)
}

func (s *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_DeleteNonExistingConfigOnAllScope_ExpectError() {
	// Given
	s.productKv.On("Delete", "key1").Return(nats.ErrKeyNotFound)
	s.workflowKv.On("Delete", "key1").Return(nats.ErrKeyNotFound)
	s.processKv.On("Delete", "key1").Return(nats.ErrKeyNotFound)

	config, err := centralizedConfiguration.NewCentralizedConfigurationBuilder(
		s.logger,
		&s.productKv,
		&s.workflowKv,
		&s.processKv,
	)
	s.Require().NoError(err)

	// When
	err = config.DeleteConfig("key1", messaging.ProductScope)
	s.Error(err)
	err = config.DeleteConfig("key1", messaging.WorkflowScope)
	s.Error(err)
	err = config.DeleteConfig("key1", messaging.ProcessScope)
	s.Error(err)

	// Then
	s.NotNil(config)
	s.productKv.AssertNumberOfCalls(s.T(), "Delete", 1)
	s.workflowKv.AssertNumberOfCalls(s.T(), "Delete", 1)
	s.processKv.AssertNumberOfCalls(s.T(), "Delete", 1)
}
