package centralizedconfiguration_test

import (
	"errors"

	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk/centralized-configuration"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk/messaging"
)

func (s *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_SetConfigOnDefaultScope_ExpectOK() {
	// Given
	s.productKv.On("PutString", "key1", "value1").Return(uint64(0), errors.New("not exist"))
	s.workflowKv.On("PutString", "key1", "value1").Return(uint64(0), errors.New("not exist"))
	s.processKv.On("PutString", "key1", "value1").Return(uint64(1), nil)

	config, err := centralizedconfiguration.NewCentralizedConfigurationBuilder(
		s.logger,
		&s.productKv,
		&s.workflowKv,
		&s.processKv,
	)

	// When
	err = config.SetConfig("key1", "value1")

	// Then
	s.NotNil(config)
	s.NoError(err)
	s.productKv.AssertNotCalled(s.T(), "PutString")
	s.workflowKv.AssertNotCalled(s.T(), "PutString")
	s.processKv.AssertNumberOfCalls(s.T(), "PutString", 1)
}

func (s *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_SetConfigOnProductScope_ExpectOK() {
	// Given
	s.productKv.On("PutString", "key1", "value1").Return(uint64(1), nil)
	s.workflowKv.On("PutString", "key1", "value1").Return(uint64(0), errors.New("not exist"))
	s.processKv.On("PutString", "key1", "value1").Return(uint64(0), errors.New("not exist"))

	config, err := centralizedconfiguration.NewCentralizedConfigurationBuilder(
		s.logger,
		&s.productKv,
		&s.workflowKv,
		&s.processKv,
	)

	// When
	err = config.SetConfig("key1", "value1", messaging.ProductScope)

	// Then
	s.NotNil(config)
	s.NoError(err)
	s.productKv.AssertNumberOfCalls(s.T(), "PutString", 1)
	s.productKv.AssertCalled(s.T(), "PutString", "key1", "value1")
	s.workflowKv.AssertNotCalled(s.T(), "PutString")
	s.processKv.AssertNotCalled(s.T(), "PutString")
}

func (s *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_SetConfigOnWorkflowScope_ExpectOK() {
	// Given
	s.productKv.On("PutString", "key1", "value1").Return(uint64(0), errors.New("not exist"))
	s.workflowKv.On("PutString", "key1", "value1").Return(uint64(1), nil)
	s.processKv.On("PutString", "key1", "value1").Return(uint64(0), errors.New("not exist"))

	config, err := centralizedconfiguration.NewCentralizedConfigurationBuilder(
		s.logger,
		&s.productKv,
		&s.workflowKv,
		&s.processKv,
	)

	// When
	err = config.SetConfig("key1", "value1", messaging.WorkflowScope)

	// Then
	s.NotNil(config)
	s.NoError(err)
	s.productKv.AssertNotCalled(s.T(), "PutString")
	s.workflowKv.AssertNumberOfCalls(s.T(), "PutString", 1)
	s.workflowKv.AssertCalled(s.T(), "PutString", "key1", "value1")
	s.processKv.AssertNotCalled(s.T(), "PutString")
}

func (s *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_SetConfigOnProcessScope_ExpectOK() {
	// Given
	s.productKv.On("PutString", "key1", "value1").Return(uint64(0), errors.New("not exist"))
	s.workflowKv.On("PutString", "key1", "value1").Return(uint64(0), errors.New("not exist"))
	s.processKv.On("PutString", "key1", "value1").Return(uint64(1), nil)
	config, err := centralizedconfiguration.NewCentralizedConfigurationBuilder(
		s.logger,
		&s.productKv,
		&s.workflowKv,
		&s.processKv,
	)

	// When
	err = config.SetConfig("key1", "value1", messaging.ProcessScope)

	// Then
	s.NotNil(config)
	s.NoError(err)
	s.productKv.AssertNotCalled(s.T(), "PutString")
	s.workflowKv.AssertNotCalled(s.T(), "PutString")
	s.processKv.AssertNumberOfCalls(s.T(), "PutString", 1)
	s.processKv.AssertCalled(s.T(), "PutString", "key1", "value1")
}

func (s *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_SetConfigOnDefaultScope_UnexpectedError_ExpectError() {
	// Given
	s.productKv.On("PutString", "key1", "value1").Return(uint64(0), errors.New("not exist"))
	s.workflowKv.On("PutString", "key1", "value1").Return(uint64(0), errors.New("not exist"))
	s.processKv.On("PutString", "key1", "value1").Return(uint64(0), errors.New("not exist"))

	config, err := centralizedconfiguration.NewCentralizedConfigurationBuilder(
		s.logger,
		&s.productKv,
		&s.workflowKv,
		&s.processKv,
	)

	// When
	err = config.SetConfig("key1", "value1")

	// Then
	s.NotNil(config)
	s.Error(err)
	s.productKv.AssertNotCalled(s.T(), "PutString")
	s.workflowKv.AssertNotCalled(s.T(), "PutString")
	s.processKv.AssertNumberOfCalls(s.T(), "PutString", 1)
}
