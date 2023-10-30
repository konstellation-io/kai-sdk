//go:build integration

package persistentstorage

import (
	"github.com/go-logr/logr"
	"github.com/spf13/viper"
)

func NewPersistentStorageBuilder(logger logr.Logger, mock persistentStorageInterface) (*PersistentStorage, error) {
	persistentStorageBucket := viper.GetString("minio.bucket")

	return &PersistentStorage{
		logger:                  logger,
		persistentStorage:       mock,
		persistentStorageBucket: persistentStorageBucket,
	}, nil
}
