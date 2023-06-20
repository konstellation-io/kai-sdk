package object_store_test

import (
	"fmt"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/internal/errors"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk/object_store"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
)

func (suite *SdkObjectStoreTestSuite) TestObjectStore_SaveObjectStoreNotInitialized_ExpectError() {
	// Given
	viper.SetDefault("nats.object_store", "")
	objectStore, err := object_store.NewObjectStore(suite.logger, &suite.jetstream)

	// When
	err = objectStore.Save("key", []byte("value"))

	// Then
	suite.ErrorIs(err, errors.ErrUndefinedObjectStore)
	suite.NotNil(objectStore)
}

func (suite *SdkObjectStoreTestSuite) TestObjectStore_SaveEmptyPayload_ExpectError() {
	// Given
	viper.SetDefault("nats.object_store", "object-store")
	suite.jetstream.On("ObjectStore", "object-store").Return(&suite.objectStore, nil)
	objectStore, err := object_store.NewObjectStore(suite.logger, &suite.jetstream)

	suite.objectStore.On("PutBytes", "key", []byte("value")).Return(&nats.ObjectInfo{}, nil)

	// When
	err = objectStore.Save("key", nil)

	// Then
	suite.ErrorIs(err, errors.ErrEmptyPayload)
	suite.NotNil(objectStore)
}

func (suite *SdkObjectStoreTestSuite) TestObjectStore_ErrorSavingPayload_ExpectError() {
	// Given
	viper.SetDefault("nats.object_store", "object-store")
	suite.jetstream.On("ObjectStore", "object-store").Return(&suite.objectStore, nil)
	objectStore, err := object_store.NewObjectStore(suite.logger, &suite.jetstream)

	suite.objectStore.On("PutBytes", "key", []byte("value")).Return(nil, fmt.Errorf("error saving payload"))

	// When
	err = objectStore.Save("key", []byte("value"))

	// Then
	suite.Error(err)
	suite.NotNil(objectStore)
}

func (suite *SdkObjectStoreTestSuite) TestObjectStore_SaveObject_ExpectOK() {
	// Given
	viper.SetDefault("nats.object_store", "object-store")
	suite.jetstream.On("ObjectStore", "object-store").Return(&suite.objectStore, nil)
	objectStore, err := object_store.NewObjectStore(suite.logger, &suite.jetstream)

	suite.objectStore.On("PutBytes", "key", []byte("value")).Return(&nats.ObjectInfo{}, nil)

	// When
	err = objectStore.Save("key", []byte("value"))

	// Then
	suite.NoError(err)
	suite.NotNil(objectStore)
	suite.objectStore.AssertNumberOfCalls(suite.T(), "PutBytes", 1)
}
