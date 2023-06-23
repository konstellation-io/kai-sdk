package common_test

import (
	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/mocks"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/runner/common"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"testing"
)

type RunnerCommonTestSuite struct {
	suite.Suite
	logger logr.Logger
	sdk    sdk.KaiSDK
}

func (suite *RunnerCommonTestSuite) SetupSuite() {
	suite.logger = testr.NewWithOptions(suite.T(), testr.Options{Verbosity: 1})
}

func (suite *RunnerCommonTestSuite) SetupTest() {
	// Reset viper values before each test
	viper.Reset()

	suite.sdk = sdk.KaiSDK{
		Logger:            suite.logger,
		PathUtils:         mocks.NewPathUtilsMock(suite.T()),
		Metadata:          mocks.NewMetadataMock(suite.T()),
		Messaging:         mocks.NewMessagingMock(suite.T()),
		ObjectStore:       mocks.NewObjectStoreMock(suite.T()),
		CentralizedConfig: mocks.NewCentralizedConfigMock(suite.T()),
	}
}

func (suite *RunnerCommonTestSuite) TestInitializeProcessConfiguration_WhenCalledWithValidaData_ExpectOk() {
	//Given
	viper.SetDefault("centralized_configuration.process.config",
		map[string]string{"key1": "value1", "key2": "value2"})
	suite.sdk.CentralizedConfig.(*mocks.CentralizedConfigMock).
		On("SetConfig", "key1", "value1").Return(nil)
	suite.sdk.CentralizedConfig.(*mocks.CentralizedConfigMock).
		On("SetConfig", "key2", "value2").Return(nil)

	// When
	common.InitializeProcessConfiguration(suite.sdk)

	//Then
	suite.sdk.CentralizedConfig.(*mocks.CentralizedConfigMock).AssertNumberOfCalls(suite.T(), "SetConfig", 2)
}

func (suite *RunnerCommonTestSuite) TestInitializeProcessConfiguration_WhenCalledWithNoData_ExpectOk() {
	// When
	common.InitializeProcessConfiguration(suite.sdk)

	//Then
	suite.sdk.CentralizedConfig.(*mocks.CentralizedConfigMock).AssertNotCalled(suite.T(), "SetConfig")
}

func TestRunnerCommonTestSuite(t *testing.T) {
	suite.Run(t, new(RunnerCommonTestSuite))
}
