package objectstore_test

import (
	"fmt"
	objectstore "github.com/konstellation-io/kre-runners/go-sdk/v1/sdk/object-store"

	"github.com/konstellation-io/kre-runners/go-sdk/v1/internal/errors"
	"github.com/spf13/viper"
)

func (s *SdkObjectStoreTestSuite) TestObjectStore_GetObjectStoreNotInitialized_ExpectError() {
	// Given
	viper.SetDefault("nats.object-store", "")
	objectStore, err := objectstore.NewObjectStore(s.logger, &s.jetstream)

	// When
	value, err := objectStore.Get("key")

	// Then
	s.ErrorIs(err, errors.ErrUndefinedObjectStore)
	s.NotNil(objectStore)
	s.Nil(value)
}

func (s *SdkObjectStoreTestSuite) TestObjectStore_ErrorRetrievingObject_ExpectError() {
	// Given
	viper.SetDefault("nats.object-store", "object-store")
	s.jetstream.On("ObjectStore", "object-store").Return(&s.objectStore, nil)
	objectStore, err := objectstore.NewObjectStore(s.logger, &s.jetstream)

	s.objectStore.On("GetBytes", "key").Return(nil, fmt.Errorf("error saving object"))

	// When
	value, err := objectStore.Get("key")

	// Then
	s.Error(err)
	s.NotNil(objectStore)
	s.Empty(value)
	s.objectStore.AssertNumberOfCalls(s.T(), "GetBytes", 1)
}

func (s *SdkObjectStoreTestSuite) TestObjectStore_GetObject_ExpectOK() {
	// Given
	viper.SetDefault("nats.object-store", "object-store")
	s.jetstream.On("ObjectStore", "object-store").Return(&s.objectStore, nil)
	objectStore, err := objectstore.NewObjectStore(s.logger, &s.jetstream)

	s.objectStore.On("GetBytes", "key").Return([]byte("value"), nil)

	// When
	value, err := objectStore.Get("key")

	// Then
	s.NoError(err)
	s.NotNil(objectStore)
	s.Equal("value", string(value))
	s.objectStore.AssertNumberOfCalls(s.T(), "GetBytes", 1)
}
