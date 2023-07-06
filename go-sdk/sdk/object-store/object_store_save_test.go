package objectstore_test

import (
	"fmt"
	objectstore "github.com/konstellation-io/kre-runners/go-sdk/v1/sdk/object-store"

	"github.com/konstellation-io/kre-runners/go-sdk/v1/internal/errors"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
)

func (s *SdkObjectStoreTestSuite) TestObjectStore_SaveObjectStoreNotInitialized_ExpectError() {
	// Given
	viper.SetDefault("nats.object-store", "")
	objectStore, err := objectstore.NewObjectStore(s.logger, &s.jetstream)

	// When
	err = objectStore.Save("key", []byte("value"))

	// Then
	s.ErrorIs(err, errors.ErrUndefinedObjectStore)
	s.NotNil(objectStore)
}

func (s *SdkObjectStoreTestSuite) TestObjectStore_SaveEmptyPayload_ExpectError() {
	// Given
	viper.SetDefault("nats.object-store", "object-store")
	s.jetstream.On("ObjectStore", "object-store").Return(&s.objectStore, nil)
	objectStore, err := objectstore.NewObjectStore(s.logger, &s.jetstream)

	s.objectStore.On("PutBytes", "key", []byte("value")).Return(&nats.ObjectInfo{}, nil)

	// When
	err = objectStore.Save("key", nil)

	// Then
	s.ErrorIs(err, errors.ErrEmptyPayload)
	s.NotNil(objectStore)
}

func (s *SdkObjectStoreTestSuite) TestObjectStore_ErrorSavingPayload_ExpectError() {
	// Given
	viper.SetDefault("nats.object-store", "object-store")
	s.jetstream.On("ObjectStore", "object-store").Return(&s.objectStore, nil)
	objectStore, err := objectstore.NewObjectStore(s.logger, &s.jetstream)

	s.objectStore.On("PutBytes", "key", []byte("value")).Return(nil, fmt.Errorf("error saving payload"))

	// When
	err = objectStore.Save("key", []byte("value"))

	// Then
	s.Error(err)
	s.NotNil(objectStore)
}

func (s *SdkObjectStoreTestSuite) TestObjectStore_SaveObject_ExpectOK() {
	// Given
	viper.SetDefault("nats.object-store", "object-store")
	s.jetstream.On("ObjectStore", "object-store").Return(&s.objectStore, nil)
	objectStore, err := objectstore.NewObjectStore(s.logger, &s.jetstream)

	s.objectStore.On("PutBytes", "key", []byte("value")).Return(&nats.ObjectInfo{}, nil)

	// When
	err = objectStore.Save("key", []byte("value"))

	// Then
	s.NoError(err)
	s.NotNil(objectStore)
	s.objectStore.AssertNumberOfCalls(s.T(), "PutBytes", 1)
}
