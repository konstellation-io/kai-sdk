package common_test

import (
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"

	"github.com/konstellation-io/kre-runners/go-sdk/v1/mocks"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/runner/common"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk"
)

type RunnerCommonTestSuite struct {
	suite.Suite
	logger logr.Logger
	sdk    sdk.KaiSDK
}

func (s *RunnerCommonTestSuite) SetupSuite() {
	s.logger = testr.NewWithOptions(s.T(), testr.Options{Verbosity: 1})
}

func (s *RunnerCommonTestSuite) SetupTest() {
	// Reset viper values before each test
	viper.Reset()

	s.sdk = sdk.KaiSDK{
		Logger:            s.logger,
		PathUtils:         mocks.NewPathUtilsMock(s.T()),
		Metadata:          mocks.NewMetadataMock(s.T()),
		Messaging:         mocks.NewMessagingMock(s.T()),
		ObjectStore:       mocks.NewObjectStoreMock(s.T()),
		CentralizedConfig: mocks.NewCentralizedConfigMock(s.T()),
	}
}

func (s *RunnerCommonTestSuite) TestInitializeProcessConfiguration_WhenCalledWithValidaData_ExpectOk() {
	viper.SetDefault("centralized-configuration.process.config",
		map[string]string{"key1": "value1", "key2": "value2"})
	s.sdk.CentralizedConfig.(*mocks.CentralizedConfigMock).
		On("SetConfig", "key1", "value1").Return(nil)
	s.sdk.CentralizedConfig.(*mocks.CentralizedConfigMock).
		On("SetConfig", "key2", "value2").Return(nil)

	// When
	common.InitializeProcessConfiguration(s.sdk)

	s.sdk.CentralizedConfig.(*mocks.CentralizedConfigMock).AssertNumberOfCalls(s.T(), "SetConfig", 2)
}

func (s *RunnerCommonTestSuite) TestInitializeProcessConfiguration_WhenCalledWithNoData_ExpectOk() {
	// When
	common.InitializeProcessConfiguration(s.sdk)

	s.sdk.CentralizedConfig.(*mocks.CentralizedConfigMock).AssertNotCalled(s.T(), "SetConfig")
}

func TestRunnerCommonTestSuite(t *testing.T) {
	suite.Run(t, new(RunnerCommonTestSuite))
}
