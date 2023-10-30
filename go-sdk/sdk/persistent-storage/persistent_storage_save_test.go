//go:build integration

package persistentstorage_test

import (
	"fmt"

	"github.com/go-logr/logr/testr"
	"github.com/konstellation-io/kai-sdk/go-sdk/internal/errors"
	persistentstorage "github.com/konstellation-io/kai-sdk/go-sdk/sdk/persistent-storage"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/mock"
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

func (s *SdkPersistentStorageTestSuite) TestPersistentStorage_SaveObject_ExpectError() {
	// GIVEN
	key := "some-key"
	data := []byte("some-data")

	s.minioClientMock.On("PutObject", mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(minio.UploadInfo{}, fmt.Errorf("some-error"))

	// WHEN
	pStorage, err := persistentstorage.NewPersistentStorageBuilder(testr.New(s.T()), &s.minioClientMock)
	s.Require().NoError(err)
	returnedVersion, err := pStorage.Save(key, data)

	// THEN
	s.Assert().Error(err)
	s.Assert().Empty(returnedVersion)
}
