//go:build integration

package persistentstorage

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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

func NewPersistentStorageIntegration(logger logr.Logger) (*PersistentStorage, error) {
	persistentStorageBucket := viper.GetString("minio.bucket")

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
	endpoint := viper.GetString("minio.endpoint")
	useSSL := viper.GetBool("minio.ssl")
	user := viper.GetString("minio.client_user")
	password := viper.GetString("minio.client_password")

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
