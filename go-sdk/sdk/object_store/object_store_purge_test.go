package object_store_test

import (
	"github.com/konstellation-io/kre-runners/go-sdk/v1/internal/errors"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk/object_store"
	"github.com/spf13/viper"
)

// TODO to be finished
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

/*func (suite *SdkObjectStoreTestSuite) TestObjectStore_PurgeObject_ExpectOK() {
	// Given
	viper.SetDefault("nats.object_store", "object-store")
	suite.jetstream.On("ObjectStore", "object-store").Return(&suite.objectStore, nil)
	objectStore, err := object_store.NewObjectStore(suite.logger, &suite.jetstream)

	suite.objectStore.On("Delete", "key").Return(nil)

	// When
	err = objectStore.DeleteMany("key")

	// Then
	suite.NoError(err)
	suite.NotNil(objectStore)
	suite.objectStore.AssertNumberOfCalls(suite.T(), "Delete", 1)
}*/
