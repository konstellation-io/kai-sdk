package centralizedconfiguration_test

import (
	"errors"

	centralizedconfiguration "github.com/konstellation-io/kai-sdk/go-sdk/v1/sdk/centralized-configuration"
	"github.com/konstellation-io/kai-sdk/go-sdk/v1/sdk/messaging"
)

const notExist = "not exist"

func (s *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_SetConfigOnDefaultScope_ExpectOK() {
	// Given

	s.productKv.On("PutString", "key1", "value1").Return(uint64(0), errors.New(notExist))
	s.workflowKv.On("PutString", "key1", "value1").Return(uint64(0), errors.New(notExist))
	s.processKv.On("PutString", "key1", "value1").Return(uint64(1), nil)

	config, err := centralizedconfiguration.NewCentralizedConfigurationBuilder(
		s.logger,
		&s.productKv,
		&s.workflowKv,
		&s.processKv,
	)
	s.Require().NoError(err)

	// When
	err = config.SetConfig("key1", "value1")
	s.Require().NoError(err)

	// Then
	s.NotNil(config)
	s.productKv.AssertNotCalled(s.T(), "PutString")
	s.workflowKv.AssertNotCalled(s.T(), "PutString")
	s.processKv.AssertNumberOfCalls(s.T(), "PutString", 1)
}

func (s *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_SetConfigOnProductScope_ExpectOK() {
	// Given
	s.productKv.On("PutString", "key1", "value1").Return(uint64(1), nil)
	s.workflowKv.On("PutString", "key1", "value1").Return(uint64(0), errors.New(notExist))
	s.processKv.On("PutString", "key1", "value1").Return(uint64(0), errors.New(notExist))

	config, err := centralizedconfiguration.NewCentralizedConfigurationBuilder(
		s.logger,
		&s.productKv,
		&s.workflowKv,
		&s.processKv,
	)
	s.Require().NoError(err)

	// When
	err = config.SetConfig("key1", "value1", messaging.ProductScope)
	s.Require().NoError(err)

	// Then
	s.NotNil(config)
	s.productKv.AssertNumberOfCalls(s.T(), "PutString", 1)
	s.productKv.AssertCalled(s.T(), "PutString", "key1", "value1")
	s.workflowKv.AssertNotCalled(s.T(), "PutString")
	s.processKv.AssertNotCalled(s.T(), "PutString")
}

func (s *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_SetConfigOnWorkflowScope_ExpectOK() {
	// Given
	s.productKv.On("PutString", "key1", "value1").Return(uint64(0), errors.New(notExist))
	s.workflowKv.On("PutString", "key1", "value1").Return(uint64(1), nil)
	s.processKv.On("PutString", "key1", "value1").Return(uint64(0), errors.New(notExist))

	config, err := centralizedconfiguration.NewCentralizedConfigurationBuilder(
		s.logger,
		&s.productKv,
		&s.workflowKv,
		&s.processKv,
	)
	s.Require().NoError(err)

	// When
	err = config.SetConfig("key1", "value1", messaging.WorkflowScope)
	s.Require().NoError(err)

	// Then
	s.NotNil(config)
	s.productKv.AssertNotCalled(s.T(), "PutString")
	s.workflowKv.AssertNumberOfCalls(s.T(), "PutString", 1)
	s.workflowKv.AssertCalled(s.T(), "PutString", "key1", "value1")
	s.processKv.AssertNotCalled(s.T(), "PutString")
}

func (s *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_SetConfigOnProcessScope_ExpectOK() {
	// Given
	s.productKv.On("PutString", "key1", "value1").Return(uint64(0), errors.New(notExist))
	s.workflowKv.On("PutString", "key1", "value1").Return(uint64(0), errors.New(notExist))
	s.processKv.On("PutString", "key1", "value1").Return(uint64(1), nil)
	config, err := centralizedconfiguration.NewCentralizedConfigurationBuilder(
		s.logger,
		&s.productKv,
		&s.workflowKv,
		&s.processKv,
	)
	s.Require().NoError(err)

	// When
	err = config.SetConfig("key1", "value1", messaging.ProcessScope)
	s.Require().NoError(err)

	// Then
	s.NotNil(config)
	s.productKv.AssertNotCalled(s.T(), "PutString")
	s.workflowKv.AssertNotCalled(s.T(), "PutString")
	s.processKv.AssertNumberOfCalls(s.T(), "PutString", 1)
	s.processKv.AssertCalled(s.T(), "PutString", "key1", "value1")
}

func (s *SdkCentralizedConfigurationTestSuite) TestCentralizedConfiguration_SetConfigOnDefaultScope_UnexpectedError_ExpectError() {
	// Given
	s.productKv.On("PutString", "key1", "value1").Return(uint64(0), errors.New(notExist))
	s.workflowKv.On("PutString", "key1", "value1").Return(uint64(0), errors.New(notExist))
	s.processKv.On("PutString", "key1", "value1").Return(uint64(0), errors.New(notExist))

	config, err := centralizedconfiguration.NewCentralizedConfigurationBuilder(
		s.logger,
		&s.productKv,
		&s.workflowKv,
		&s.processKv,
	)
	s.Require().NoError(err)

	// When
	err = config.SetConfig("key1", "value1")

	// Then
	s.Error(err)
	s.NotNil(config)
	s.productKv.AssertNotCalled(s.T(), "PutString")
	s.workflowKv.AssertNotCalled(s.T(), "PutString")
	s.processKv.AssertNumberOfCalls(s.T(), "PutString", 1)
}
