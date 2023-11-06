package persistentstorage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai-sdk/go-sdk/internal/errors"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
)

type PersistentStorage struct {
	logger                  logr.Logger
	persistentStorage       persistentStorageInterface
	persistentStorageBucket string
}

//go:generate mockery --name persistentStorageInterface --output ../../mocks --structname MinioClientMock --filename minio_client_mock.go
type persistentStorageInterface interface {
	PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64,
		opts minio.PutObjectOptions) (minio.UploadInfo, error)
	GetObject(ctx context.Context, bucketName, objectName string,
		opts minio.GetObjectOptions) (*minio.Object, error)
	ListObjects(ctx context.Context, bucketName string, opts minio.ListObjectsOptions) <-chan minio.ObjectInfo
	RemoveObject(ctx context.Context, bucketName, objectName string, opts minio.RemoveObjectOptions) error
}

func NewPersistentStorage(logger logr.Logger) (*PersistentStorage, error) {
	persistentStorageBucket := viper.GetString("minio.bucket")

	storageManager, err := initPersistentStorage(logger)
	if err != nil {
		return nil, err
	}

	return &PersistentStorage{
		logger:                  logger,
		persistentStorage:       storageManager,
		persistentStorageBucket: persistentStorageBucket,
	}, nil
}

func initPersistentStorage(logger logr.Logger) (*minio.Client, error) {
	endpoint := viper.GetString("minio.endpoint")
	accessKeyID := viper.GetString("minio.client_user")
	secretAccessKey := viper.GetString("minio.client_password")
	useSSL := viper.GetBool("minio.ssl")

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:        credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure:       useSSL,
		BucketLookup: minio.BucketLookupPath,
	})
	if err != nil {
		return nil, fmt.Errorf("error initializing persistent storage: %w", err)
	}

	logger.Info("Successfully initialized persistent storage")

	return minioClient, nil
}

func (ps PersistentStorage) Save(key string, payload []byte, ttlDays ...int) (string, error) {
	if key == "" {
		return "", errors.ErrEmptyKey
	}

	if payload == nil {
		return "", errors.ErrEmptyPayload
	}

	reader := bytes.NewReader(payload)

	opts := minio.PutObjectOptions{}

	if len(ttlDays) > 0 && ttlDays[0] > 0 {
		opts.RetainUntilDate = time.Now().AddDate(0, 0, ttlDays[0])
	}

	info, err := ps.persistentStorage.PutObject(
		context.Background(),
		ps.persistentStorageBucket,
		key,
		reader,
		int64(reader.Len()),
		opts,
	)
	if err != nil {
		return "", fmt.Errorf("error storing object to the persistent storage: %w", err)
	}

	ps.logger.V(1).Info("Object %s successfully stored in persistent storage with version ID %s",
		key, info.VersionID)

	return info.VersionID, nil
}

func (ps PersistentStorage) Get(key string, version ...string) ([]byte, error) {
	if key == "" {
		return nil, errors.ErrEmptyKey
	}

	opts := minio.GetObjectOptions{}

	if len(version) > 0 && version[0] != "" {
		opts = minio.GetObjectOptions{
			VersionID: version[0],
		}
	}

	object, err := ps.persistentStorage.GetObject(
		context.Background(),
		ps.persistentStorageBucket,
		key,
		opts,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving object from the persistent storage: %w", err)
	}

	ps.logger.V(1).Info("Object %s successfully retrieved from persistent storage ", key)

	defer object.Close()

	data, err := io.ReadAll(object)
	if err != nil {
		ps.logger.V(1).Error(err, "Error reading object from persistent storage ", key)
	}

	return data, nil
}

func (ps PersistentStorage) List() ([]string, error) {
	var objectList []string

	objects := ps.persistentStorage.ListObjects(
		context.Background(),
		ps.persistentStorageBucket,
		minio.ListObjectsOptions{
			WithVersions: true,
			WithMetadata: true,
			Recursive:    true,
		},
	)

	ps.logger.V(1).Info("Objects successfully retrieved from persistent storage")

	for object := range objects {
		if object.Key != "" && object.VersionID != "" {
			objectList = append(objectList, fmt.Sprintf("%s - %s", object.Key, object.VersionID))
		}
	}

	return objectList, nil
}

func (ps PersistentStorage) ListVersions(key string) ([]string, error) {
	var objectList []string

	if key == "" {
		return nil, errors.ErrEmptyKey
	}

	objects := ps.persistentStorage.ListObjects(
		context.Background(),
		ps.persistentStorageBucket,
		minio.ListObjectsOptions{
			WithVersions: true,
			WithMetadata: true,
			Recursive:    true,
			Prefix:       key,
		},
	)

	ps.logger.V(1).Info("Object versions successfully retrieved for prefix %s from persistent storage", key)

	for object := range objects {
		if object.Key != "" && object.VersionID != "" {
			objectList = append(objectList, fmt.Sprintf("%s - %s", object.Key, object.VersionID))
		}
	}

	return objectList, nil
}

func (ps PersistentStorage) Delete(key string, version ...string) error {
	if key == "" {
		return errors.ErrEmptyKey
	}

	opts := minio.RemoveObjectOptions{
		ForceDelete: true,
	}

	if len(version) > 0 && version[0] != "" {
		opts.VersionID = version[0]
	}

	err := ps.persistentStorage.RemoveObject(
		context.Background(),
		ps.persistentStorageBucket,
		key,
		opts,
	)
	if err != nil {
		return fmt.Errorf("error deleting object from the persistent storage: %w", err)
	}

	ps.logger.V(1).Info("Object %s successfully deleted from persistent storage", key)

	return nil
}
