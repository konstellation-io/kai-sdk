package object_store_test

import (
	"fmt"

	"github.com/konstellation-io/kre-runners/go-sdk/v1/internal/errors"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk/object_store"
	"github.com/spf13/viper"
)

func (suite *SdkObjectStoreTestSuite) TestObjectStore_ListObjectStoreNotInitialized_ExpectError() {
	// Given
	viper.SetDefault("nats.object_store", "")
	objectStore, err := object_store.NewObjectStore(suite.logger, &suite.jetstream)

	// When
	keyList, err := objectStore.List()

	// Then
	suite.ErrorIs(err, errors.ErrUndefinedObjectStore)
	suite.Empty(keyList)
	suite.NotNil(objectStore)
}

func (suite *SdkObjectStoreTestSuite) TestObjectStore_ErrorListObject_ExpectError() {
	// Given
	viper.SetDefault("nats.object_store", "object-store")
	suite.jetstream.On("ObjectStore", "object-store").Return(&suite.objectStore, nil)
	objectStore, err := object_store.NewObjectStore(suite.logger, &suite.jetstream)

	suite.objectStore.On("List").Return(nil, fmt.Errorf("error saving object"))

	// When
	keyList, err := objectStore.List("key")

	// Then
	suite.Error(err)
	suite.NotNil(objectStore)
	suite.Nil(keyList)
	suite.objectStore.AssertNumberOfCalls(suite.T(), "List", 1)
}

func (suite *SdkObjectStoreTestSuite) TestObjectStore_ListObject_ExpectOK() {
	// Given
	viper.SetDefault("nats.object_store", "object-store")
	suite.jetstream.On("ObjectStore", "object-store").Return(&suite.objectStore, nil)
	objectStore, err := object_store.NewObjectStore(suite.logger, &suite.jetstream)

	keys := []string{"key1", "key2"}
	objects := generateObjectInfoResponse(keys)

	suite.objectStore.On("List").Return(objects, nil)

	// When
	keyList, err := objectStore.List()

	// Then
	suite.NoError(err)
	suite.NotNil(objectStore)
	suite.NotNil(keyList)
	suite.objectStore.AssertNumberOfCalls(suite.T(), "List", 1)
	suite.Equal([]string{"key1", "key2"}, keyList)
}

func (suite *SdkObjectStoreTestSuite) TestObjectStore_ListObjectWithFilter_ExpectOK() {
	// Given
	viper.SetDefault("nats.object_store", "object-store")
	suite.jetstream.On("ObjectStore", "object-store").Return(&suite.objectStore, nil)
	objectStore, err := object_store.NewObjectStore(suite.logger, &suite.jetstream)

	keys := []string{"key1", "key2", "anotherKey3", "anotherKey4"}
	objects := generateObjectInfoResponse(keys)

	suite.objectStore.On("List").Return(objects, nil)

	// When
	keyList, err := objectStore.List("anotherKey*")

	// Then
	suite.NoError(err)
	suite.NotNil(objectStore)
	suite.NotNil(keyList)
	suite.objectStore.AssertNumberOfCalls(suite.T(), "List", 1)
	suite.Equal([]string{"anotherKey3", "anotherKey4"}, keyList)
}

func (suite *SdkObjectStoreTestSuite) TestObjectStore_ListObjectWithFilter_InvalidRegexp_ExpectError() {
	// Given
	viper.SetDefault("nats.object_store", "object-store")
	suite.jetstream.On("ObjectStore", "object-store").Return(&suite.objectStore, nil)
	objectStore, err := object_store.NewObjectStore(suite.logger, &suite.jetstream)

	keys := []string{"key1", "key2", "anotherKey3", "anotherKey4"}
	objects := generateObjectInfoResponse(keys)

	suite.objectStore.On("List").Return(objects, nil)

	// When
	keyList, err := objectStore.List("anotherKey*[")

	// Then
	suite.Error(err)
	suite.NotNil(objectStore)
	suite.Nil(keyList)
	suite.objectStore.AssertNumberOfCalls(suite.T(), "List", 1)
}

func (suite *SdkObjectStoreTestSuite) TestObjectStore_ListObject_ErrorOnList_ExpectError() {
	// Given
	viper.SetDefault("nats.object_store", "object-store")
	suite.jetstream.On("ObjectStore", "object-store").Return(&suite.objectStore, nil)
	objectStore, err := object_store.NewObjectStore(suite.logger, &suite.jetstream)

	suite.objectStore.On("List").Return(nil, fmt.Errorf("error listing objects"))

	// When
	keyList, err := objectStore.List()

	// Then
	suite.Error(err)
	suite.NotNil(objectStore)
	suite.Nil(keyList)
	suite.objectStore.AssertNumberOfCalls(suite.T(), "List", 1)
}
