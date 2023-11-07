//go:build integration

package persistentstorage

import (
	"fmt"
	"github.com/konstellation-io/kai-sdk/go-sdk/internal/common"

	"github.com/go-logr/logr"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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

func NewPersistentStorageIntegration(logger logr.Logger) (*PersistentStorage, error) {
	persistentStorageBucket := viper.GetString(common.ConfigMinioBucketKey)

	storageManager, err := initPersistentStorageIntegration(logger)
	if err != nil {
		return nil, err
	}

	return &PersistentStorage{
		logger:                  logger,
		persistentStorage:       storageManager,
		persistentStorageBucket: persistentStorageBucket,
	}, nil
}

func initPersistentStorageIntegration(logger logr.Logger) (*minio.Client, error) {
	endpoint := viper.GetString(common.ConfigMinioEndpointKey)
	useSSL := viper.GetBool(common.ConfigMinioUseSslKey)
	user := viper.GetString(common.ConfigMinioClientUserKey)
	password := viper.GetString(common.ConfigMinioClientPasswordKey)

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:        credentials.NewStaticV4(user, password, ""),
		Secure:       useSSL,
		BucketLookup: minio.BucketLookupPath,
	})
	if err != nil {
		return nil, fmt.Errorf("error initializing persistent storage: %w", err)
	}

	logger.Info("Successfully initialized persistent storage")

	return minioClient, nil
}
