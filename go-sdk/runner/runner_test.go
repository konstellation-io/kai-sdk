package runner_test

import (
	"fmt"
	"testing"

	"github.com/konstellation-io/kre-runners/go-sdk/v1/mocks"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/runner"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type SdkRunnerTestSuite struct {
	suite.Suite
	js mocks.JetStreamContextMock
}

func (s *SdkRunnerTestSuite) SetupSuite() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../../testdata")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error initializing configuration: %w", err))
	}
}

func (s *SdkRunnerTestSuite) SetupTest() {
	s.js = *mocks.NewJetStreamContextMock(s.T())
}

func (s *SdkRunnerTestSuite) TestNewTriggerRunnerInitialization_ExpectOK() {
	// Given
	s.js.On("KeyValue", mock.AnythingOfType("string")).Return(mocks.NewKeyValueMock(s.T()), nil)
	s.js.On("ObjectStore", mock.AnythingOfType("string")).Return(mocks.NewNatsObjectStoreMock(s.T()), nil)

	// When
	triggerRunner := runner.NewTestRunner(nil, &s.js).TriggerRunner()
	// Then
	s.NotNil(triggerRunner)
}

func (s *SdkRunnerTestSuite) TestNewTriggerRunner_WithInitializer_ExpectPanic() {
	// Given
	s.js.On("KeyValue", mock.AnythingOfType("string")).Return(mocks.NewKeyValueMock(s.T()), nil)
	s.js.On("ObjectStore", mock.AnythingOfType("string")).Return(mocks.NewNatsObjectStoreMock(s.T()), nil)

	// Then
	s.Panicsf(func() {
		// When
		runner.NewTestRunner(nil, &s.js).
			TriggerRunner().
			WithInitializer(func(sdk sdk.KaiSDK) {}).
			Run()
	}, "No runner function defined")
}

func (s *SdkRunnerTestSuite) TestNewTaskRunnerInitialization_ExpectOK() {
	// Given
	s.js.On("KeyValue", mock.AnythingOfType("string")).Return(mocks.NewKeyValueMock(s.T()), nil)
	s.js.On("ObjectStore", mock.AnythingOfType("string")).Return(mocks.NewNatsObjectStoreMock(s.T()), nil)

	// When
	triggerRunner := runner.NewTestRunner(nil, &s.js).TaskRunner()
	// Then
	s.NotNil(triggerRunner)
}

func (s *SdkRunnerTestSuite) TestNewTaskRunner_WithInitializer_ExpectPanic() {
	// Given
	s.js.On("KeyValue", mock.AnythingOfType("string")).Return(mocks.NewKeyValueMock(s.T()), nil)
	s.js.On("ObjectStore", mock.AnythingOfType("string")).Return(mocks.NewNatsObjectStoreMock(s.T()), nil)

	// Then
	s.Panicsf(func() {
		// When
		runner.NewTestRunner(nil, &s.js).
			TaskRunner().
			WithInitializer(func(sdk sdk.KaiSDK) {}).
			Run()
	}, "No default handler defined")
}

func (s *SdkRunnerTestSuite) TestNewExitRunnerInitialization_ExpectOK() {
	// Given
	s.js.On("KeyValue", mock.AnythingOfType("string")).Return(mocks.NewKeyValueMock(s.T()), nil)
	s.js.On("ObjectStore", mock.AnythingOfType("string")).Return(mocks.NewNatsObjectStoreMock(s.T()), nil)

	// When
	triggerRunner := runner.NewTestRunner(nil, &s.js).ExitRunner()
	// Then
	s.NotNil(triggerRunner)
}

func (s *SdkRunnerTestSuite) TestNewExitRunner_WithInitializer_ExpectPanic() {
	// Given
	s.js.On("KeyValue", mock.AnythingOfType("string")).Return(mocks.NewKeyValueMock(s.T()), nil)
	s.js.On("ObjectStore", mock.AnythingOfType("string")).Return(mocks.NewNatsObjectStoreMock(s.T()), nil)

	// Then
	s.Panicsf(func() {
		// When
		runner.NewTestRunner(nil, &s.js).
			ExitRunner().
			WithInitializer(func(sdk sdk.KaiSDK) {}).
			Run()
	}, "No default handler defined")
}

func TestRunnerTestSuite(t *testing.T) {
	suite.Run(t, new(SdkRunnerTestSuite))
}
