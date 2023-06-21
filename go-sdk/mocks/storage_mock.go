// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// StorageMock is an autogenerated mock type for the storage type
type StorageMock struct {
	mock.Mock
}

type StorageMock_Expecter struct {
	mock *mock.Mock
}

func (_m *StorageMock) EXPECT() *StorageMock_Expecter {
	return &StorageMock_Expecter{mock: &_m.Mock}
}

// NewStorageMock creates a new instance of StorageMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStorageMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *StorageMock {
	mock := &StorageMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
