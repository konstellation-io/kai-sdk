package objectstore_test

import (
	"fmt"

	objectstore "github.com/konstellation-io/kre-runners/go-sdk/v1/sdk/object-store"

	"github.com/konstellation-io/kre-runners/go-sdk/v1/internal/errors"
	"github.com/spf13/viper"
)

func (s *SdkObjectStoreTestSuite) TestObjectStore_DeleteObjectStoreNotInitialized_ExpectError() {
	// Given
	viper.SetDefault("nats.object-store", "")
	objectStore, err := objectstore.NewObjectStore(s.logger, &s.jetstream)

	// When
	err = objectStore.Delete("key")

	// Then
	s.ErrorIs(err, errors.ErrUndefinedObjectStore)
	s.NotNil(objectStore)
}

func (s *SdkObjectStoreTestSuite) TestObjectStore_ErrorDeletingObject_ExpectError() {
	// Given
	viper.SetDefault("nats.object-store", "object-store")
	s.jetstream.On("ObjectStore", "object-store").Return(&s.objectStore, nil)
	objectStore, err := objectstore.NewObjectStore(s.logger, &s.jetstream)

	s.objectStore.On("Delete", "key").Return(fmt.Errorf("error saving object"))

	// When
	err = objectStore.Delete("key")

	// Then
	s.Error(err)
	s.NotNil(objectStore)
	s.objectStore.AssertNumberOfCalls(s.T(), "Delete", 1)
}

func (s *SdkObjectStoreTestSuite) TestObjectStore_DeleteObject_ExpectOK() {
	// Given
	viper.SetDefault("nats.object-store", "object-store")
	s.jetstream.On("ObjectStore", "object-store").Return(&s.objectStore, nil)
	objectStore, err := objectstore.NewObjectStore(s.logger, &s.jetstream)

	s.objectStore.On("Delete", "key").Return(nil)

	// When
	err = objectStore.Delete("key")

	// Then
	s.NoError(err)
	s.NotNil(objectStore)
	s.objectStore.AssertNumberOfCalls(s.T(), "Delete", 1)
}
