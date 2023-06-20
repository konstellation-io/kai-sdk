package path_utils

import (
	"github.com/go-logr/logr"
	"github.com/spf13/viper"
	"path"
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
	return viper.GetString("metadata.base_path")
}
func (pu PathUtils) ComposeBasePath(relativePath string) string {
	return path.Join(pu.GetBasePath(), relativePath)
}
