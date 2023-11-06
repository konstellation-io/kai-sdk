//go:build integration

package persistentstorage

import (
	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai-sdk/go-sdk/internal/common"
	"github.com/spf13/viper"
)

func NewPersistentStorageBuilder(logger logr.Logger, mock persistentStorageInterface) *PersistentStorage {
	persistentStorageBucket := viper.GetString(common.ConfigMinioBucketKey)

	return &PersistentStorage{
		logger:                  logger,
		persistentStorage:       mock,
		persistentStorageBucket: persistentStorageBucket,
	}
}
