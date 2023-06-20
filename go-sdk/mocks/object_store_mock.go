// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ObjectStoreMock is an autogenerated mock type for the objectStore type
type ObjectStoreMock struct {
	mock.Mock
}

type ObjectStoreMock_Expecter struct {
	mock *mock.Mock
}

func (_m *ObjectStoreMock) EXPECT() *ObjectStoreMock_Expecter {
	return &ObjectStoreMock_Expecter{mock: &_m.Mock}
}

// Delete provides a mock function with given fields: key
func (_m *ObjectStoreMock) Delete(key string) error {
	ret := _m.Called(key)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ObjectStoreMock_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type ObjectStoreMock_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - key string
func (_e *ObjectStoreMock_Expecter) Delete(key interface{}) *ObjectStoreMock_Delete_Call {
	return &ObjectStoreMock_Delete_Call{Call: _e.mock.On("Delete", key)}
}

func (_c *ObjectStoreMock_Delete_Call) Run(run func(key string)) *ObjectStoreMock_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *ObjectStoreMock_Delete_Call) Return(_a0 error) *ObjectStoreMock_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ObjectStoreMock_Delete_Call) RunAndReturn(run func(string) error) *ObjectStoreMock_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: key
func (_m *ObjectStoreMock) Get(key string) ([]byte, error) {
	ret := _m.Called(key)

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]byte, error)); ok {
		return rf(key)
	}
	if rf, ok := ret.Get(0).(func(string) []byte); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ObjectStoreMock_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type ObjectStoreMock_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - key string
func (_e *ObjectStoreMock_Expecter) Get(key interface{}) *ObjectStoreMock_Get_Call {
	return &ObjectStoreMock_Get_Call{Call: _e.mock.On("Get", key)}
}

func (_c *ObjectStoreMock_Get_Call) Run(run func(key string)) *ObjectStoreMock_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *ObjectStoreMock_Get_Call) Return(_a0 []byte, _a1 error) *ObjectStoreMock_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ObjectStoreMock_Get_Call) RunAndReturn(run func(string) ([]byte, error)) *ObjectStoreMock_Get_Call {
	_c.Call.Return(run)
	return _c
}

// Save provides a mock function with given fields: key, value
func (_m *ObjectStoreMock) Save(key string, value []byte) error {
	ret := _m.Called(key, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []byte) error); ok {
		r0 = rf(key, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ObjectStoreMock_Save_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Save'
type ObjectStoreMock_Save_Call struct {
	*mock.Call
}

// Save is a helper method to define mock.On call
//   - key string
//   - value []byte
func (_e *ObjectStoreMock_Expecter) Save(key interface{}, value interface{}) *ObjectStoreMock_Save_Call {
	return &ObjectStoreMock_Save_Call{Call: _e.mock.On("Save", key, value)}
}

func (_c *ObjectStoreMock_Save_Call) Run(run func(key string, value []byte)) *ObjectStoreMock_Save_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].([]byte))
	})
	return _c
}

func (_c *ObjectStoreMock_Save_Call) Return(_a0 error) *ObjectStoreMock_Save_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ObjectStoreMock_Save_Call) RunAndReturn(run func(string, []byte) error) *ObjectStoreMock_Save_Call {
	_c.Call.Return(run)
	return _c
}

// NewObjectStoreMock creates a new instance of ObjectStoreMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewObjectStoreMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *ObjectStoreMock {
	mock := &ObjectStoreMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
