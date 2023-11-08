package persistentstorage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/konstellation-io/kai-sdk/go-sdk/sdk/metadata"
	"github.com/minio/minio-go/v7/pkg/lifecycle"

	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai-sdk/go-sdk/internal/auth"
	"github.com/konstellation-io/kai-sdk/go-sdk/internal/common"
	"github.com/konstellation-io/kai-sdk/go-sdk/internal/errors"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
)

type PersistentStorage struct {
	logger                  logr.Logger
	persistentStorage       *minio.Client
	persistentStorageBucket string
	metadata                *metadata.Metadata
}

type ObjectInfo struct {
	Key       string
	VersionID string
	ExpiresIn time.Time
}

type Object struct {
	ObjectInfo
	data []byte
}

func (o Object) GetAsString() string {
	return string(o.data)
}

func (o Object) GetBytes() []byte {
	return o.data
}

func NewPersistentStorage(logger logr.Logger) (*PersistentStorage, error) {
	persistentStorageBucket := viper.GetString(common.ConfigMinioBucketKey)

	storageManager, err := initPersistentStorage(logger)
	if err != nil {
		return nil, err
	}

	return &PersistentStorage{
		logger:                  logger,
		persistentStorage:       storageManager,
		persistentStorageBucket: persistentStorageBucket,
		metadata:                metadata.NewMetadata(),
	}, nil
}

func initPersistentStorage(logger logr.Logger) (*minio.Client, error) {
	endpoint := viper.GetString(common.ConfigMinioEndpointKey)
	useSSL := viper.GetBool(common.ConfigMinioUseSslKey)
	url := ""

	if useSSL {
		url = fmt.Sprintf("%s://%s", "https", viper.GetString(common.ConfigMinioEndpointKey))
	} else {
		url = fmt.Sprintf("%s://%s", "http", viper.GetString(common.ConfigMinioEndpointKey))
	}

	minioCredentials, err := getClientCredentials(logger, url)
	if err != nil {
		return nil, fmt.Errorf("error initializing persistent storage: %w", err)
	}

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:        minioCredentials,
		Secure:       useSSL,
		BucketLookup: minio.BucketLookupPath,
	})
	if err != nil {
		return nil, fmt.Errorf("error initializing persistent storage: %w", err)
	}

	logger.Info("Successfully initialized persistent storage")

	return minioClient, nil
}

func (ps PersistentStorage) Save(key string, payload []byte, ttlDays ...int) (*ObjectInfo, error) {
	ctx := context.Background()

	if key == "" {
		return nil, errors.ErrEmptyKey
	}

	if len(payload) == 0 {
		return nil, errors.ErrEmptyPayload
	}

	err := ps.addLifecycleDeletionRule(key, ttlDays, ctx)
	if err != nil {
		return nil, fmt.Errorf("error adding lifecycle deletion rule: %w", err)
	}

	reader := bytes.NewReader(payload)

	opts := minio.PutObjectOptions{
		UserTags: map[string]string{
			"product":  ps.metadata.GetProduct(),
			"version":  ps.metadata.GetVersion(),
			"workflow": ps.metadata.GetWorkflow(),
			"process":  ps.metadata.GetProcess(),
		},
	}

	info, err := ps.persistentStorage.PutObject(
		ctx,
		ps.persistentStorageBucket,
		key,
		reader,
		int64(reader.Len()),
		opts,
	)
	if err != nil {
		return nil, fmt.Errorf("error storing object to the persistent storage: %w", err)
	}

	ps.logger.V(1).Info(fmt.Sprintf("Object %s successfully stored in persistent storage with version ID %s",
		key, info.VersionID))

	obj := &ObjectInfo{
		Key:       key,
		VersionID: info.VersionID,
		ExpiresIn: info.Expiration,
	}

	return obj, nil
}

func (ps PersistentStorage) Get(key string, version ...string) (*Object, error) {
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

	ps.logger.V(1).Info(fmt.Sprintf("Object %s successfully retrieved from persistent storage ", key))

	defer object.Close()

	data, err := io.ReadAll(object)
	if err != nil {
		return nil, fmt.Errorf("error reading object from the persistent storage: %w", err)
	}

	objStats, err := object.Stat()
	if err != nil {
		return nil, fmt.Errorf("error getting object stats from the persistent storage: %w", err)
	}

	obj := &Object{
		ObjectInfo: ObjectInfo{
			Key:       key,
			VersionID: objStats.VersionID,
			ExpiresIn: objStats.Expiration,
		},
		data: data,
	}

	return obj, nil
}

func (ps PersistentStorage) List() ([]*ObjectInfo, error) {
	var objectList []*ObjectInfo

	objects := ps.persistentStorage.ListObjects(
		context.Background(),
		ps.persistentStorageBucket,
		minio.ListObjectsOptions{
			WithMetadata: true,
			Recursive:    true,
		},
	)

	ps.logger.V(1).Info("Objects successfully retrieved from persistent storage")

	for object := range objects {
		if object.Key != "" && object.VersionID != "" {
			objectList = append(
				objectList,
				&ObjectInfo{
					Key:       object.Key,
					VersionID: object.VersionID,
					ExpiresIn: object.Expiration,
				},
			)
		}
	}

	return objectList, nil
}

func (ps PersistentStorage) ListVersions(key string) ([]*ObjectInfo, error) {
	var objectList []*ObjectInfo

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

	ps.logger.V(1).
		Info(fmt.Sprintf("Object versions successfully retrieved for prefix %s from persistent storage", key))

	for object := range objects {
		if object.VersionID != "" {
			objectList = append(
				objectList,
				&ObjectInfo{
					Key:       object.Key,
					VersionID: object.VersionID,
					ExpiresIn: object.Expiration,
				},
			)
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

	ps.logger.V(1).Info(fmt.Sprintf("Object %s successfully deleted from persistent storage", key))

	return nil
}

func getClientCredentials(logger logr.Logger, url string) (*credentials.Credentials, error) {
	minioCreds, err := credentials.NewSTSClientGrants(
		url,
		func() (*credentials.ClientGrantsToken, error) {
			authClient := auth.New(logger)
			token, err := authClient.GetToken()
			if err != nil {
				return nil, err
			}

			return &credentials.ClientGrantsToken{
				Token:  token.AccessToken,
				Expiry: token.ExpiresIn,
			}, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return minioCreds, err
}

func (ps PersistentStorage) addLifecycleDeletionRule(key string, ttlDays []int, ctx context.Context) error {
	if len(ttlDays) > 0 && ttlDays[0] > 0 {
		lc, err := ps.persistentStorage.GetBucketLifecycle(ctx, ps.persistentStorageBucket)
		if err != nil {
			lc = lifecycle.NewConfiguration()
		}

		lc.Rules = append(lc.Rules, lifecycle.Rule{
			ID:     fmt.Sprintf("ttl-%s", key),
			Status: minio.Enabled,
			//Prefix: key,
			RuleFilter: lifecycle.Filter{
				Prefix: key,
			},
			Expiration: lifecycle.Expiration{
				Days: lifecycle.ExpirationDays(ttlDays[0]),
			},
		})

		err = ps.persistentStorage.SetBucketLifecycle(ctx, ps.persistentStorageBucket, lc)
		if err != nil {
			return err
		}
	}

	return nil
}
