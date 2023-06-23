package runner_test

import (
	"fmt"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/mocks"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/konstellation-io/kre-runners/go-sdk/v1/runner"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type SdkRunnerTestSuite struct {
	suite.Suite
	js mocks.JetStreamContextMock
}

func (suite *SdkRunnerTestSuite) SetupSuite() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../../testdata")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error initializing configuration: %w", err))

	}
}

func (suite *SdkRunnerTestSuite) SetupTest() {
	suite.js = *mocks.NewJetStreamContextMock(suite.T())
}

func (suite *SdkRunnerTestSuite) TestNewTriggerRunnerInitialization_ExpectOK() {
	// Given
	suite.js.On("KeyValue", mock.AnythingOfType("string")).Return(mocks.NewKeyValueMock(suite.T()), nil)
	suite.js.On("ObjectStore", mock.AnythingOfType("string")).Return(mocks.NewNatsObjectStoreMock(suite.T()), nil)

	// When
	triggerRunner := runner.NewTestRunner(nil, &suite.js).TriggerRunner()
	// Then
	suite.NotNil(triggerRunner)
}

func (suite *SdkRunnerTestSuite) TestNewTriggerRunner_WithInitializer_ExpectPanic() {
	// Given
	suite.js.On("KeyValue", mock.AnythingOfType("string")).Return(mocks.NewKeyValueMock(suite.T()), nil)
	suite.js.On("ObjectStore", mock.AnythingOfType("string")).Return(mocks.NewNatsObjectStoreMock(suite.T()), nil)

	// Then
	suite.Panicsf(func() {
		// When
		runner.NewTestRunner(nil, &suite.js).
			TriggerRunner().
			WithInitializer(func(sdk sdk.KaiSDK) {}).
			Run()
	}, "No runner function defined")
}

func (suite *SdkRunnerTestSuite) TestNewTaskRunnerInitialization_ExpectOK() {
	// Given
	suite.js.On("KeyValue", mock.AnythingOfType("string")).Return(mocks.NewKeyValueMock(suite.T()), nil)
	suite.js.On("ObjectStore", mock.AnythingOfType("string")).Return(mocks.NewNatsObjectStoreMock(suite.T()), nil)

	// When
	triggerRunner := runner.NewTestRunner(nil, &suite.js).TaskRunner()
	// Then
	suite.NotNil(triggerRunner)
}

func (suite *SdkRunnerTestSuite) TestNewTaskRunner_WithInitializer_ExpectPanic() {
	// Given
	suite.js.On("KeyValue", mock.AnythingOfType("string")).Return(mocks.NewKeyValueMock(suite.T()), nil)
	suite.js.On("ObjectStore", mock.AnythingOfType("string")).Return(mocks.NewNatsObjectStoreMock(suite.T()), nil)

	// Then
	suite.Panicsf(func() {
		// When
		runner.NewTestRunner(nil, &suite.js).
			TaskRunner().
			WithInitializer(func(sdk sdk.KaiSDK) {}).
			Run()
	}, "No default handler defined")
}

func (suite *SdkRunnerTestSuite) TestNewExitRunnerInitialization_ExpectOK() {
	// Given
	suite.js.On("KeyValue", mock.AnythingOfType("string")).Return(mocks.NewKeyValueMock(suite.T()), nil)
	suite.js.On("ObjectStore", mock.AnythingOfType("string")).Return(mocks.NewNatsObjectStoreMock(suite.T()), nil)

	// When
	triggerRunner := runner.NewTestRunner(nil, &suite.js).ExitRunner()
	// Then
	suite.NotNil(triggerRunner)
}

func (suite *SdkRunnerTestSuite) TestNewExitRunner_WithInitializer_ExpectPanic() {
	// Given
	suite.js.On("KeyValue", mock.AnythingOfType("string")).Return(mocks.NewKeyValueMock(suite.T()), nil)
	suite.js.On("ObjectStore", mock.AnythingOfType("string")).Return(mocks.NewNatsObjectStoreMock(suite.T()), nil)

	// Then
	suite.Panicsf(func() {
		// When
		runner.NewTestRunner(nil, &suite.js).
			ExitRunner().
			WithInitializer(func(sdk sdk.KaiSDK) {}).
			Run()
	}, "No default handler defined")
}

func TestRunnerTestSuite(t *testing.T) {
	suite.Run(t, new(SdkRunnerTestSuite))
}
