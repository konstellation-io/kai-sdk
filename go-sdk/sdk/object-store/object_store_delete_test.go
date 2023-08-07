package objectstore_test

import (
	"fmt"

	objectstore "github.com/konstellation-io/kai-sdk/go-sdk/sdk/object-store"

	"github.com/spf13/viper"

	"github.com/konstellation-io/kai-sdk/go-sdk/internal/errors"
)

const (
	natsObjectStoreField = "nats.object_store"
	natsObjectStoreValue = "object-store"
)

func (s *SdkObjectStoreTestSuite) TestObjectStore_DeleteObjectStoreNotInitialized_ExpectError() {
	// Given
	viper.SetDefault(natsObjectStoreField, "")
	
	objectStore, _ := objectstore.NewObjectStore(s.logger, &s.jetstream)

	// When
	err := objectStore.Delete("key")

	// Then
	s.ErrorIs(err, errors.ErrUndefinedObjectStore)
	s.NotNil(objectStore)
}

func (s *SdkObjectStoreTestSuite) TestObjectStore_ErrorDeletingObject_ExpectError() {
	// Given
	viper.SetDefault(natsObjectStoreField, natsObjectStoreValue)
	s.jetstream.On("ObjectStore", natsObjectStoreValue).Return(&s.objectStore, nil)
	objectStore, _ := objectstore.NewObjectStore(s.logger, &s.jetstream)
	
	s.objectStore.On("Delete", "key").Return(fmt.Errorf("error saving object"))

	// When
	err := objectStore.Delete("key")

	// Then
	s.Error(err)
	s.NotNil(objectStore)
	s.objectStore.AssertNumberOfCalls(s.T(), "Delete", 1)
}

func (s *SdkObjectStoreTestSuite) TestObjectStore_DeleteObject_ExpectOK() {
	// Given
	viper.SetDefault(natsObjectStoreField, natsObjectStoreValue)
	s.jetstream.On("ObjectStore", natsObjectStoreValue).Return(&s.objectStore, nil)
	objectStore, _ := objectstore.NewObjectStore(s.logger, &s.jetstream)
	
	s.objectStore.On("Delete", "key").Return(nil)

	// When
	err := objectStore.Delete("key")

	// Then
	s.NoError(err)
	s.NotNil(objectStore)
	s.objectStore.AssertNumberOfCalls(s.T(), "Delete", 1)
}
