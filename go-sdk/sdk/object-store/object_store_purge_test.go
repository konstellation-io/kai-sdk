package objectstore_test

import (
	"fmt"

	objectStore2 "github.com/konstellation-io/kai-sdk/go-sdk/v1/sdk/object-store"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"

	"github.com/konstellation-io/kai-sdk/go-sdk/v1/internal/errors"
)

func (s *SdkObjectStoreTestSuite) TestObjectStore_PurgeObjectStoreNotInitialized_ExpectError() {
	// Given
	viper.SetDefault(natsObjectStoreField, "")
	objectStore, err := objectStore2.NewObjectStore(s.logger, &s.jetstream)

	// When
	err = objectStore.Purge("key")

	// Then
	s.ErrorIs(err, errors.ErrUndefinedObjectStore)
	s.NotNil(objectStore)
}

func (s *SdkObjectStoreTestSuite) TestObjectStore_ErrorPurgingObject_ExpectError() {
	// Given
	viper.SetDefault(natsObjectStoreField, natsObjectStoreValue)
	s.jetstream.On("ObjectStore", natsObjectStoreValue).Return(&s.objectStore, nil)
	objectStore, err := objectStore2.NewObjectStore(s.logger, &s.jetstream)

	keys := []string{"key1", "key2"}
	objects := generateObjectInfoResponse(keys)

	s.objectStore.On("List").Return(objects, nil)
	s.objectStore.On("Delete", mock.AnythingOfType("string")).
		Return(fmt.Errorf("error saving object"))

	// When
	err = objectStore.Purge("key")

	// Then
	s.Error(err)
	s.NotNil(objectStore)
	s.objectStore.AssertNumberOfCalls(s.T(), "Delete", 1)
}

func (s *SdkObjectStoreTestSuite) TestObjectStore_PurgeObject_ExpectOK() {
	// Given
	viper.SetDefault(natsObjectStoreField, natsObjectStoreValue)
	s.jetstream.On("ObjectStore", natsObjectStoreValue).Return(&s.objectStore, nil)
	objectStore, err := objectStore2.NewObjectStore(s.logger, &s.jetstream)

	keys := []string{"key1", "key2"}
	objects := generateObjectInfoResponse(keys)

	s.objectStore.On("List").Return(objects, nil)
	s.objectStore.On("Delete", mock.AnythingOfType("string")).Return(nil)

	// When
	err = objectStore.Purge()

	// Then
	s.NoError(err)
	s.NotNil(objectStore)
	s.objectStore.AssertNumberOfCalls(s.T(), "List", 1)
	s.objectStore.AssertNumberOfCalls(s.T(), "Delete", 2)
}

func (s *SdkObjectStoreTestSuite) TestObjectStore_PurgeObjectWithFilter_ExpectOK() {
	// Given
	viper.SetDefault(natsObjectStoreField, natsObjectStoreValue)
	s.jetstream.On("ObjectStore", natsObjectStoreValue).Return(&s.objectStore, nil)
	objectStore, err := objectStore2.NewObjectStore(s.logger, &s.jetstream)

	keys := []string{"key1", "key2", "anotherKey3", "anotherKey4"}
	objects := generateObjectInfoResponse(keys)

	s.objectStore.On("List").Return(objects, nil)
	s.objectStore.On("Delete", mock.AnythingOfType("string")).Return(nil)

	// When
	err = objectStore.Purge("anotherKey*")

	// Then
	s.NoError(err)
	s.NotNil(objectStore)
	s.objectStore.AssertNumberOfCalls(s.T(), "List", 1)
	s.objectStore.AssertNumberOfCalls(s.T(), "Delete", 2)
	s.objectStore.AssertCalled(s.T(), "Delete", "anotherKey3")
	s.objectStore.AssertCalled(s.T(), "Delete", "anotherKey4")
}

func (s *SdkObjectStoreTestSuite) TestObjectStore_PurgeObjectWithFilter_InvalidRegexp_ExpectError() {
	// Given
	viper.SetDefault(natsObjectStoreField, natsObjectStoreValue)
	s.jetstream.On("ObjectStore", natsObjectStoreValue).Return(&s.objectStore, nil)
	objectStore, err := objectStore2.NewObjectStore(s.logger, &s.jetstream)

	s.objectStore.On("Delete", mock.AnythingOfType("string")).Return(nil)

	// When
	err = objectStore.Purge("anotherKey*[")

	// Then
	s.Error(err)
	s.NotNil(objectStore)
	s.objectStore.AssertNotCalled(s.T(), "List")
	s.objectStore.AssertNotCalled(s.T(), "Delete")
}

func (s *SdkObjectStoreTestSuite) TestObjectStore_PurgeObjectWithFilter_ErrorListingKeys_ExpectError() {
	// Given
	viper.SetDefault(natsObjectStoreField, natsObjectStoreValue)
	s.jetstream.On("ObjectStore", natsObjectStoreValue).Return(&s.objectStore, nil)
	objectStore, err := objectStore2.NewObjectStore(s.logger, &s.jetstream)

	s.objectStore.On("List").Return(nil, fmt.Errorf("error listing keys"))

	// When
	err = objectStore.Purge()

	// Then
	s.Error(err)
	s.NotNil(objectStore)
	s.objectStore.AssertNumberOfCalls(s.T(), "List", 1)
	s.objectStore.AssertNotCalled(s.T(), "Delete")
}
