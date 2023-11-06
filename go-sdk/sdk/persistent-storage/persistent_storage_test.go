//go:build integration

package persistentstorage_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-logr/logr/testr"
	"github.com/konstellation-io/kai-sdk/go-sdk/mocks"
	persistentstorage "github.com/konstellation-io/kai-sdk/go-sdk/sdk/persistent-storage"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type SdkPersistentStorageTestSuite struct {
	suite.Suite
	persistentStorage       *persistentstorage.PersistentStorage
	minioContainer          testcontainers.Container
	client                  *minio.Client
	persistentStorageBucket string
	minioClientMock         mocks.MinioClientMock
}

func (s *SdkPersistentStorageTestSuite) getContainer() testcontainers.Container {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "minio/minio:latest",
		ExposedPorts: []string{"9000/tcp"},
		Cmd: []string{
			"server",
			"/data",
		},
		Env:        map[string]string{},
		WaitingFor: wait.ForLog("Status:         1 Online, 0 Offline."),
	}

	minioContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	s.Require().NoError(err)

	return minioContainer
}

func (s *SdkPersistentStorageTestSuite) getClient(url, keyID, keySecret string) *minio.Client {
	client, err := minio.New(url, &minio.Options{
		Creds:        credentials.NewStaticV4(keyID, keySecret, ""),
		BucketLookup: minio.BucketLookupPath,
	})
	s.Require().NoError(err)

	return client
}

func (s *SdkPersistentStorageTestSuite) SetupSuite() {
	s.minioContainer = s.getContainer()
}

func (s *SdkPersistentStorageTestSuite) SetupTest() {
	// Reset viper values before each test
	viper.Reset()

	ctx := context.Background()

	s.persistentStorageBucket = "test-bucket"

	host, err := s.minioContainer.Host(ctx)
	s.Require().NoError(err)

	port, err := s.minioContainer.MappedPort(ctx, "9000/tcp")
	s.Require().NoError(err)

	minioEndpoint := fmt.Sprintf("%s:%d", host, port.Int())

	viper.Set("minio.bucket", s.persistentStorageBucket)
	viper.Set("minio.bucket", s.persistentStorageBucket)
	viper.Set("minio.endpoint", minioEndpoint)
	viper.Set("minio.client_user", "minioadmin")
	viper.Set("minio.client_password", "minioadmin")
	viper.Set("minio.ssl", false)

	s.client = s.getClient(minioEndpoint, "minioadmin", "minioadmin")

	err = s.client.MakeBucket(ctx, s.persistentStorageBucket, minio.MakeBucketOptions{})
	s.Require().NoError(err)

	err = s.client.EnableVersioning(ctx, s.persistentStorageBucket)
	s.Require().NoError(err)

	storage, err := persistentstorage.NewPersistentStorage(testr.New(s.T()))
	s.Assert().NoError(err)
	s.persistentStorage = storage

	s.minioClientMock = *mocks.NewMinioClientMock(s.T())
}

func (s *SdkPersistentStorageTestSuite) TearDownTest() {
	ctx := context.Background()

	err := s.client.RemoveBucketWithOptions(
		ctx,
		s.persistentStorageBucket,
		minio.RemoveBucketOptions{ForceDelete: true},
	)
	if err != nil {
		var minioErr minio.ErrorResponse

		errors.As(err, &minioErr)

		if minioErr.StatusCode != http.StatusNotFound {
			s.Failf(err.Error(), "Error deleting bucket %q", s.persistentStorageBucket)
		}
	}
}

func (s *SdkPersistentStorageTestSuite) TestPersistentStorage_Initialize_ExpectOK() {
	// WHEN
	persistentStorage, err := persistentstorage.NewPersistentStorage(testr.New(s.T()))
	s.Require().NoError(err)
	_, err = persistentStorage.List()

	// THEN
	s.Assert().NoError(err)
	s.Assert().NotNil(persistentStorage)
}

func (s *SdkPersistentStorageTestSuite) TestPersistentStorage_Initialize_ExpectError() {
	// GIVEN
	viper.Set("minio.endpoint", "invalid-url")

	// WHEN
	persistentStorage, err := persistentstorage.NewPersistentStorage(testr.New(s.T()))
	s.Require().NoError(err)
	objList, err := persistentStorage.List()

	// THEN
	s.Assert().NoError(err)
	s.Assert().Len(objList, 0)
	s.Assert().NotNil(persistentStorage)
}

func TestSdkPersistentStorageTestSuite(t *testing.T) {
	suite.Run(t, new(SdkPersistentStorageTestSuite))
}
