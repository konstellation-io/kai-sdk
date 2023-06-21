package path_utils_test

import (
	"fmt"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk/path_utils"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
)

type SdkPathUtilsTestSuite struct {
	suite.Suite
	logger logr.Logger
}

func (suite *SdkPathUtilsTestSuite) SetupTest() {
	zapLog, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Sprintf("who watches the watchmen (%v)?", err))
	}

	// Reset viper values before each test
	viper.Reset()

	suite.logger = zapr.NewLogger(zapLog)

	viper.SetDefault("metadata.base_path", "/base/path")
}

func (suite *SdkPathUtilsTestSuite) TestPathUtils_GetBasePath_ExpectOK() {
	// Given
	pathUtils := path_utils.NewPathUtils(suite.logger)

	// When
	basePath := pathUtils.GetBasePath()

	// Then
	suite.Equal("/base/path", basePath)
}

func (suite *SdkPathUtilsTestSuite) TestPathUtils_ComposePath_NoElements_ExpectOK() {
	// Given
	pathUtils := path_utils.NewPathUtils(suite.logger)

	// When
	basePath := pathUtils.ComposePath()

	// Then
	suite.Equal("/base/path", basePath)
}

func (suite *SdkPathUtilsTestSuite) TestPathUtils_ComposePath_OneElements_ExpectOK() {
	// Given
	pathUtils := path_utils.NewPathUtils(suite.logger)

	// When
	basePath := pathUtils.ComposePath("test")

	// Then
	suite.Equal("/base/path/test", basePath)
}

func (suite *SdkPathUtilsTestSuite) TestPathUtils_ComposePath_MultipleElements_ExpectOK() {
	// Given
	pathUtils := path_utils.NewPathUtils(suite.logger)

	// When
	basePath := pathUtils.ComposePath("test1", "test2", "test3")

	// Then
	suite.Equal("/base/path/test1/test2/test3", basePath)
}

func TestSdkPathUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(SdkPathUtilsTestSuite))
}
