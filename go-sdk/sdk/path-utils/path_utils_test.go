package pathutils_test

import (
	"testing"

	pathUtils2 "github.com/konstellation-io/kre-runners/go-sdk/v1/sdk/path-utils"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type SdkPathUtilsTestSuite struct {
	suite.Suite
	logger logr.Logger
}

func (s *SdkPathUtilsTestSuite) SetupSuite() {
	s.logger = testr.NewWithOptions(s.T(), testr.Options{Verbosity: 1, LogTimestamp: true})
}

func (s *SdkPathUtilsTestSuite) SetupTest() {
	// Reset viper values before each test
	viper.Reset()

	viper.SetDefault("metadata.base_path", "/base/path")
}

func (s *SdkPathUtilsTestSuite) TestPathUtils_GetBasePath_ExpectOK() {
	// Given
	pathUtils := pathUtils2.NewPathUtils(s.logger)

	// When
	basePath := pathUtils.GetBasePath()

	// Then
	s.Equal("/base/path", basePath)
}

func (s *SdkPathUtilsTestSuite) TestPathUtils_ComposePath_NoElements_ExpectOK() {
	// Given
	pathUtils := pathUtils2.NewPathUtils(s.logger)

	// When
	basePath := pathUtils.ComposePath()

	// Then
	s.Equal("/base/path", basePath)
}

func (s *SdkPathUtilsTestSuite) TestPathUtils_ComposePath_OneElements_ExpectOK() {
	// Given
	pathUtils := pathUtils2.NewPathUtils(s.logger)

	// When
	basePath := pathUtils.ComposePath("test")

	// Then
	s.Equal("/base/path/test", basePath)
}

func (s *SdkPathUtilsTestSuite) TestPathUtils_ComposePath_MultipleElements_ExpectOK() {
	// Given
	pathUtils := pathUtils2.NewPathUtils(s.logger)

	// When
	basePath := pathUtils.ComposePath("test1", "test2", "test3")

	// Then
	s.Equal("/base/path/test1/test2/test3", basePath)
}

func TestSdkPathUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(SdkPathUtilsTestSuite))
}
