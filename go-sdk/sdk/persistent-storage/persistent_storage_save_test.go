//go:build integration

package persistentstorage_test

import (
	"github.com/konstellation-io/kai-sdk/go-sdk/internal/errors"
)

func (s *SdkPersistentStorageTestSuite) TestPersistentStorage_SaveObject_ExpectOK() {
	// GIVEN
	key := "some-object"
	data := []byte("some-data")

	// WHEN
	returnedVersion, err := s.persistentStorage.Save(key, data)

	// THEN
	s.Assert().NoError(err)
	s.Assert().NotEmpty(returnedVersion)
	s.Assert().NotEmpty(returnedVersion)
}

func (s *SdkPersistentStorageTestSuite) TestPersistentStorage_SaveObject_WithTTL_ExpectOK() {
	// GIVEN
	key := "some-object"
	data := []byte("some-data")
	ttlDays := 1

	// WHEN
	returnedVersion, err := s.persistentStorage.Save(key, data, ttlDays)

	// THEN
	s.Assert().NoError(err)
	s.Assert().NotEmpty(returnedVersion)
	s.Assert().NotEmpty(returnedVersion)
}

func (s *SdkPersistentStorageTestSuite) TestPersistentStorage_SaveObjectWithNoKey_ExpectError() {
	// GIVEN
	key := ""
	data := []byte("some-data")

	// WHEN
	returnedVersion, err := s.persistentStorage.Save(key, data)

	// THEN
	s.Assert().Error(err)
	s.Assert().ErrorIs(err, errors.ErrEmptyKey)
	s.Assert().Empty(returnedVersion)
}

func (s *SdkPersistentStorageTestSuite) TestPersistentStorage_SaveObjectWithNoPayload_ExpectError() {
	// GIVEN
	key := "some-key"

	// WHEN
	returnedVersion, err := s.persistentStorage.Save(key, nil)

	// THEN
	s.Assert().Error(err)
	s.Assert().ErrorIs(err, errors.ErrEmptyPayload)
	s.Assert().Empty(returnedVersion)
}
