package object_store_test

import (
	"fmt"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/internal/errors"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk/object_store"
	"github.com/spf13/viper"
)

func (suite *SdkObjectStoreTestSuite) TestObjectStore_DeleteObjectStoreNotInitialized_ExpectError() {
	// Given
	viper.SetDefault("nats.object_store", "")
	objectStore, err := object_store.NewObjectStore(suite.logger, &suite.jetstream)

	// When
	err = objectStore.Delete("key")

	// Then
	suite.ErrorIs(err, errors.ErrUndefinedObjectStore)
	suite.NotNil(objectStore)
}

func (suite *SdkObjectStoreTestSuite) TestObjectStore_ErrorDeletingObject_ExpectError() {
	// Given
	viper.SetDefault("nats.object_store", "object-store")
	suite.jetstream.On("ObjectStore", "object-store").Return(&suite.objectStore, nil)
	objectStore, err := object_store.NewObjectStore(suite.logger, &suite.jetstream)

	suite.objectStore.On("Delete", "key").Return(fmt.Errorf("error saving object"))

	// When
	err = objectStore.Delete("key")

	// Then
	suite.Error(err)
	suite.NotNil(objectStore)
	suite.objectStore.AssertNumberOfCalls(suite.T(), "Delete", 1)
}

func (suite *SdkObjectStoreTestSuite) TestObjectStore_DeleteObject_ExpectOK() {
	// Given
	viper.SetDefault("nats.object_store", "object-store")
	suite.jetstream.On("ObjectStore", "object-store").Return(&suite.objectStore, nil)
	objectStore, err := object_store.NewObjectStore(suite.logger, &suite.jetstream)

	suite.objectStore.On("Delete", "key").Return(nil)

	// When
	err = objectStore.Delete("key")

	// Then
	suite.NoError(err)
	suite.NotNil(objectStore)
	suite.objectStore.AssertNumberOfCalls(suite.T(), "Delete", 1)
}
