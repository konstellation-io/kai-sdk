package object_store_test

import (
	"fmt"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/internal/errors"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk/object_store"
	"github.com/spf13/viper"
)

func (suite *SdkObjectStoreTestSuite) TestObjectStore_GetObjectStoreNotInitialized_ExpectError() {
	// Given
	viper.SetDefault("nats.object_store", "")
	objectStore, err := object_store.NewObjectStore(suite.logger, &suite.jetstream)

	// When
	value, err := objectStore.Get("key")

	// Then
	suite.ErrorIs(err, errors.ErrUndefinedObjectStore)
	suite.NotNil(objectStore)
	suite.Nil(value)
}

func (suite *SdkObjectStoreTestSuite) TestObjectStore_ErrorRetrievingObject_ExpectError() {
	// Given
	viper.SetDefault("nats.object_store", "object-store")
	suite.jetstream.On("ObjectStore", "object-store").Return(&suite.objectStore, nil)
	objectStore, err := object_store.NewObjectStore(suite.logger, &suite.jetstream)

	suite.objectStore.On("GetBytes", "key").Return(nil, fmt.Errorf("error saving object"))

	// When
	value, err := objectStore.Get("key")

	// Then
	suite.Error(err)
	suite.NotNil(objectStore)
	suite.Empty(value)
	suite.objectStore.AssertNumberOfCalls(suite.T(), "GetBytes", 1)
}

func (suite *SdkObjectStoreTestSuite) TestObjectStore_GetObject_ExpectOK() {
	// Given
	viper.SetDefault("nats.object_store", "object-store")
	suite.jetstream.On("ObjectStore", "object-store").Return(&suite.objectStore, nil)
	objectStore, err := object_store.NewObjectStore(suite.logger, &suite.jetstream)

	suite.objectStore.On("GetBytes", "key").Return([]byte("value"), nil)

	// When
	value, err := objectStore.Get("key")

	// Then
	suite.NoError(err)
	suite.NotNil(objectStore)
	suite.Equal("value", string(value))
	suite.objectStore.AssertNumberOfCalls(suite.T(), "GetBytes", 1)
}
