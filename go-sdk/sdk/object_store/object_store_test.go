package object_store_test

import (
	"fmt"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/internal/errors"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/mocks"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk/object_store"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
)

//go:generate mockery --dir $GOPATH/pkg/mod/github.com/nats-io/nats.go@v1.26.0 --output ../../mocks --name JetStreamContext --structname JetStreamContextMock --filename jetstream_context_mock.go
//go:generate mockery --dir $GOPATH/pkg/mod/github.com/nats-io/nats.go@v1.26.0 --output ../../mocks --name ObjectStore --structname NatsObjectStoreMock --filename nats_object_store_mock.go
type SdkObjectStoreTestSuite struct {
	suite.Suite
	logger      logr.Logger
	jetstream   mocks.JetStreamContextMock
	objectStore mocks.NatsObjectStoreMock
}

func (suite *SdkObjectStoreTestSuite) SetupTest() {
	zapLog, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Sprintf("who watches the watchmen (%v)?", err))
	}

	// Reset viper values before each test
	viper.Reset()

	suite.logger = zapr.NewLogger(zapLog)
	suite.jetstream = *mocks.NewJetStreamContextMock(suite.T())
	suite.objectStore = *mocks.NewNatsObjectStoreMock(suite.T())
}

func (suite *SdkObjectStoreTestSuite) TestObjectStore_InitializeObjectStore_ExpectOK() {
	// Given
	viper.SetDefault("nats.object_store", "object-store")

	suite.jetstream.On("ObjectStore", "object-store").Return(&suite.objectStore, nil)

	// When
	objectStore, err := object_store.NewObjectStore(suite.logger, &suite.jetstream)

	// Then
	suite.NoError(err)
	suite.NotNil(objectStore)
}

func (suite *SdkObjectStoreTestSuite) TestObjectStore_InitializeObjectStore_EmptyOsName_ExpectOK() {
	// Given
	viper.SetDefault("nats.object_store", "")

	suite.jetstream.On("ObjectStore", "object-store").Return(&suite.objectStore, nil)

	// When
	objectStore, osErr := object_store.NewObjectStore(suite.logger, &suite.jetstream)
	_, getErr := objectStore.Get("key")

	// Then
	suite.NoError(osErr)
	suite.NotNil(objectStore)
	suite.ErrorIs(getErr, errors.ErrUndefinedObjectStore)
}

func (suite *SdkObjectStoreTestSuite) TestObjectStore_InitializeObjectStore_ErrorOnOsReturn_ExpectError() {
	// Given
	viper.SetDefault("nats.object_store", "object-store")

	suite.jetstream.On("ObjectStore", "object-store").Return(nil, fmt.Errorf("not found"))

	// When
	objectStore, err := object_store.NewObjectStore(suite.logger, &suite.jetstream)

	// Then
	suite.Error(err)
	suite.Nil(objectStore)
}

func TestSdkCentralizedConfigurationTestSuite(t *testing.T) {
	suite.Run(t, new(SdkObjectStoreTestSuite))
}
