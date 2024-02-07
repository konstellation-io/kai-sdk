//go:build unit

package centralizedconfiguration_test

import (
	centralizedconfiguration "github.com/konstellation-io/kai-sdk/go-sdk/sdk/centralized-configuration"
	"github.com/nats-io/nats.go"
)

func (s *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_DeleteConfigOnWorkflowScope_ExpectOK() {
	// Given
	s.workflowKv.On("Delete", "key1").Return(nil)

	config, err := centralizedconfiguration.NewBuilder(
		s.logger,
		&s.globalKv,
		&s.productKv,
		&s.workflowKv,
		&s.processKv,
	)
	s.Require().NoError(err)

	// When
	err = config.DeleteConfig("key1", centralizedconfiguration.WorkflowScope)
	s.Require().NoError(err)

	// Then
	s.NotNil(config)
	s.workflowKv.AssertNumberOfCalls(s.T(), "Delete", 1)
}

func (s *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_DeleteConfigOnProcessScope_ExpectOK() {
	// Given
	s.processKv.On("Delete", "key1").Return(nil)

	config, err := centralizedconfiguration.NewBuilder(
		s.logger,
		&s.globalKv,
		&s.productKv,
		&s.workflowKv,
		&s.processKv,
	)
	s.Require().NoError(err)

	// When
	err = config.DeleteConfig("key1", centralizedconfiguration.ProcessScope)
	s.Require().NoError(err)

	// Then
	s.NotNil(config)
	s.processKv.AssertNumberOfCalls(s.T(), "Delete", 1)
}

func (s *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_DeleteNonExistingConfigOnAllScope_ExpectError() {
	// Given
	s.globalKv.On("Delete", "key1").Return(nats.ErrKeyNotFound)
	s.productKv.On("Delete", "key1").Return(nats.ErrKeyNotFound)
	s.workflowKv.On("Delete", "key1").Return(nats.ErrKeyNotFound)
	s.processKv.On("Delete", "key1").Return(nats.ErrKeyNotFound)

	config, err := centralizedconfiguration.NewBuilder(
		s.logger,
		&s.globalKv,
		&s.productKv,
		&s.workflowKv,
		&s.processKv,
	)
	s.Require().NoError(err)

	// When
	err = config.DeleteConfig("key1", centralizedconfiguration.GlobalScope)
	s.Error(err)
	err = config.DeleteConfig("key1", centralizedconfiguration.ProductScope)
	s.Error(err)
	err = config.DeleteConfig("key1", centralizedconfiguration.WorkflowScope)
	s.Error(err)
	err = config.DeleteConfig("key1", centralizedconfiguration.ProcessScope)
	s.Error(err)

	// Then
	s.NotNil(config)
	s.globalKv.AssertNumberOfCalls(s.T(), "Delete", 1)
	s.productKv.AssertNumberOfCalls(s.T(), "Delete", 1)
	s.workflowKv.AssertNumberOfCalls(s.T(), "Delete", 1)
	s.processKv.AssertNumberOfCalls(s.T(), "Delete", 1)
}
