package object_store_test

import (
	"fmt"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/internal/errors"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk/object_store"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
)

func (suite *SdkObjectStoreTestSuite) TestObjectStore_PurgeObjectStoreNotInitialized_ExpectError() {
	// Given
	viper.SetDefault("nats.object_store", "")
	objectStore, err := object_store.NewObjectStore(suite.logger, &suite.jetstream)

	// When
	err = objectStore.Purge("key")

	// Then
	suite.ErrorIs(err, errors.ErrUndefinedObjectStore)
	suite.NotNil(objectStore)
}

func (suite *SdkObjectStoreTestSuite) TestObjectStore_ErrorPurgingObject_ExpectError() {
	// Given
	viper.SetDefault("nats.object_store", "object-store")
	suite.jetstream.On("ObjectStore", "object-store").Return(&suite.objectStore, nil)
	objectStore, err := object_store.NewObjectStore(suite.logger, &suite.jetstream)

	keys := []string{"key1", "key2"}
	objects := generateObjectInfoResponse(keys)

	suite.objectStore.On("List").Return(objects, nil)
	suite.objectStore.On("Delete", mock.AnythingOfType("string")).
		Return(fmt.Errorf("error saving object"))

	// When
	err = objectStore.Purge("key")

	// Then
	suite.Error(err)
	suite.NotNil(objectStore)
	suite.objectStore.AssertNumberOfCalls(suite.T(), "Delete", 1)
}

func (suite *SdkObjectStoreTestSuite) TestObjectStore_PurgeObject_ExpectOK() {
	// Given
	viper.SetDefault("nats.object_store", "object-store")
	suite.jetstream.On("ObjectStore", "object-store").Return(&suite.objectStore, nil)
	objectStore, err := object_store.NewObjectStore(suite.logger, &suite.jetstream)

	keys := []string{"key1", "key2"}
	objects := generateObjectInfoResponse(keys)

	suite.objectStore.On("List").Return(objects, nil)
	suite.objectStore.On("Delete", mock.AnythingOfType("string")).Return(nil)

	// When
	err = objectStore.Purge()

	// Then
	suite.NoError(err)
	suite.NotNil(objectStore)
	suite.objectStore.AssertNumberOfCalls(suite.T(), "List", 1)
	suite.objectStore.AssertNumberOfCalls(suite.T(), "Delete", 2)
}

func (suite *SdkObjectStoreTestSuite) TestObjectStore_PurgeObjectWithFilter_ExpectOK() {
	// Given
	viper.SetDefault("nats.object_store", "object-store")
	suite.jetstream.On("ObjectStore", "object-store").Return(&suite.objectStore, nil)
	objectStore, err := object_store.NewObjectStore(suite.logger, &suite.jetstream)

	keys := []string{"key1", "key2", "anotherKey3", "anotherKey4"}
	objects := generateObjectInfoResponse(keys)

	suite.objectStore.On("List").Return(objects, nil)
	suite.objectStore.On("Delete", mock.AnythingOfType("string")).Return(nil)

	// When
	err = objectStore.Purge("anotherKey*")

	// Then
	suite.NoError(err)
	suite.NotNil(objectStore)
	suite.objectStore.AssertNumberOfCalls(suite.T(), "List", 1)
	suite.objectStore.AssertNumberOfCalls(suite.T(), "Delete", 2)
	suite.objectStore.AssertCalled(suite.T(), "Delete", "anotherKey3")
	suite.objectStore.AssertCalled(suite.T(), "Delete", "anotherKey4")
}

func (suite *SdkObjectStoreTestSuite) TestObjectStore_PurgeObjectWithFilter_InvalidRegexp_ExpectError() {
	// Given
	viper.SetDefault("nats.object_store", "object-store")
	suite.jetstream.On("ObjectStore", "object-store").Return(&suite.objectStore, nil)
	objectStore, err := object_store.NewObjectStore(suite.logger, &suite.jetstream)

	suite.objectStore.On("Delete", mock.AnythingOfType("string")).Return(nil)

	// When
	err = objectStore.Purge("anotherKey*[")

	// Then
	suite.Error(err)
	suite.NotNil(objectStore)
	suite.objectStore.AssertNotCalled(suite.T(), "List")
	suite.objectStore.AssertNotCalled(suite.T(), "Delete")
}

func (suite *SdkObjectStoreTestSuite) TestObjectStore_PurgeObjectWithFilter_ErrorListingKeys_ExpectError() {
	// Given
	viper.SetDefault("nats.object_store", "object-store")
	suite.jetstream.On("ObjectStore", "object-store").Return(&suite.objectStore, nil)
	objectStore, err := object_store.NewObjectStore(suite.logger, &suite.jetstream)

	suite.objectStore.On("List").Return(nil, fmt.Errorf("error listing keys"))

	// When
	err = objectStore.Purge()

	// Then
	suite.Error(err)
	suite.NotNil(objectStore)
	suite.objectStore.AssertNumberOfCalls(suite.T(), "List", 1)
	suite.objectStore.AssertNotCalled(suite.T(), "Delete")
}
