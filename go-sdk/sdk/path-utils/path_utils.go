package pathutils

import (
	path2 "path"

	"github.com/konstellation-io/kai-sdk/go-sdk/internal/common"

	"github.com/go-logr/logr"
	"github.com/spf13/viper"
)

type PathUtils struct {
	logger logr.Logger
}

func NewPathUtils(logger logr.Logger) *PathUtils {
	return &PathUtils{
		logger: logger.WithName("[PATH UTILS]"),
	}
}

func (pu PathUtils) GetBasePath() string {
	return viper.GetString(common.ConfigMetadataBasePathKey)
}

func (pu PathUtils) ComposePath(relativePath ...string) string {
	path := pu.GetBasePath()
	if len(relativePath) == 0 {
		return path
	}

	for _, elem := range relativePath {
		path = path2.Join(path, elem)
	}

	return path
}
