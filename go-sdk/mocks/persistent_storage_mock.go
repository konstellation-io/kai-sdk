// Code generated by mockery v2.32.0. DO NOT EDIT.

package mocks

import (
	persistentstorage "github.com/konstellation-io/kai-sdk/go-sdk/sdk/persistent-storage"
	mock "github.com/stretchr/testify/mock"
)

// PersistentStorageMock is an autogenerated mock type for the persistentStorage type
type PersistentStorageMock struct {
	mock.Mock
}

type PersistentStorageMock_Expecter struct {
	mock *mock.Mock
}

func (_m *PersistentStorageMock) EXPECT() *PersistentStorageMock_Expecter {
	return &PersistentStorageMock_Expecter{mock: &_m.Mock}
}

// Delete provides a mock function with given fields: key, version
func (_m *PersistentStorageMock) Delete(key string, version ...string) error {
	_va := make([]interface{}, len(version))
	for _i := range version {
		_va[_i] = version[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, key)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, ...string) error); ok {
		r0 = rf(key, version...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PersistentStorageMock_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type PersistentStorageMock_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - key string
//   - version ...string
func (_e *PersistentStorageMock_Expecter) Delete(key interface{}, version ...interface{}) *PersistentStorageMock_Delete_Call {
	return &PersistentStorageMock_Delete_Call{Call: _e.mock.On("Delete",
		append([]interface{}{key}, version...)...)}
}

func (_c *PersistentStorageMock_Delete_Call) Run(run func(key string, version ...string)) *PersistentStorageMock_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]string, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(string)
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *PersistentStorageMock_Delete_Call) Return(_a0 error) *PersistentStorageMock_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PersistentStorageMock_Delete_Call) RunAndReturn(run func(string, ...string) error) *PersistentStorageMock_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: key, version
func (_m *PersistentStorageMock) Get(key string, version ...string) (*persistentstorage.Object, error) {
	_va := make([]interface{}, len(version))
	for _i := range version {
		_va[_i] = version[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, key)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *persistentstorage.Object
	var r1 error
	if rf, ok := ret.Get(0).(func(string, ...string) (*persistentstorage.Object, error)); ok {
		return rf(key, version...)
	}
	if rf, ok := ret.Get(0).(func(string, ...string) *persistentstorage.Object); ok {
		r0 = rf(key, version...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*persistentstorage.Object)
		}
	}

	if rf, ok := ret.Get(1).(func(string, ...string) error); ok {
		r1 = rf(key, version...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PersistentStorageMock_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type PersistentStorageMock_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - key string
//   - version ...string
func (_e *PersistentStorageMock_Expecter) Get(key interface{}, version ...interface{}) *PersistentStorageMock_Get_Call {
	return &PersistentStorageMock_Get_Call{Call: _e.mock.On("Get",
		append([]interface{}{key}, version...)...)}
}

func (_c *PersistentStorageMock_Get_Call) Run(run func(key string, version ...string)) *PersistentStorageMock_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]string, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(string)
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *PersistentStorageMock_Get_Call) Return(_a0 *persistentstorage.Object, _a1 error) *PersistentStorageMock_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PersistentStorageMock_Get_Call) RunAndReturn(run func(string, ...string) (*persistentstorage.Object, error)) *PersistentStorageMock_Get_Call {
	_c.Call.Return(run)
	return _c
}

// List provides a mock function with given fields:
func (_m *PersistentStorageMock) List() ([]*persistentstorage.ObjectInfo, error) {
	ret := _m.Called()

	var r0 []*persistentstorage.ObjectInfo
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]*persistentstorage.ObjectInfo, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []*persistentstorage.ObjectInfo); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*persistentstorage.ObjectInfo)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PersistentStorageMock_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type PersistentStorageMock_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
func (_e *PersistentStorageMock_Expecter) List() *PersistentStorageMock_List_Call {
	return &PersistentStorageMock_List_Call{Call: _e.mock.On("List")}
}

func (_c *PersistentStorageMock_List_Call) Run(run func()) *PersistentStorageMock_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *PersistentStorageMock_List_Call) Return(_a0 []*persistentstorage.ObjectInfo, _a1 error) *PersistentStorageMock_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PersistentStorageMock_List_Call) RunAndReturn(run func() ([]*persistentstorage.ObjectInfo, error)) *PersistentStorageMock_List_Call {
	_c.Call.Return(run)
	return _c
}

// ListVersions provides a mock function with given fields: key
func (_m *PersistentStorageMock) ListVersions(key string) ([]*persistentstorage.ObjectInfo, error) {
	ret := _m.Called(key)

	var r0 []*persistentstorage.ObjectInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]*persistentstorage.ObjectInfo, error)); ok {
		return rf(key)
	}
	if rf, ok := ret.Get(0).(func(string) []*persistentstorage.ObjectInfo); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*persistentstorage.ObjectInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PersistentStorageMock_ListVersions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListVersions'
type PersistentStorageMock_ListVersions_Call struct {
	*mock.Call
}

// ListVersions is a helper method to define mock.On call
//   - key string
func (_e *PersistentStorageMock_Expecter) ListVersions(key interface{}) *PersistentStorageMock_ListVersions_Call {
	return &PersistentStorageMock_ListVersions_Call{Call: _e.mock.On("ListVersions", key)}
}

func (_c *PersistentStorageMock_ListVersions_Call) Run(run func(key string)) *PersistentStorageMock_ListVersions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *PersistentStorageMock_ListVersions_Call) Return(_a0 []*persistentstorage.ObjectInfo, _a1 error) *PersistentStorageMock_ListVersions_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PersistentStorageMock_ListVersions_Call) RunAndReturn(run func(string) ([]*persistentstorage.ObjectInfo, error)) *PersistentStorageMock_ListVersions_Call {
	_c.Call.Return(run)
	return _c
}

// Save provides a mock function with given fields: key, value, ttlDays
func (_m *PersistentStorageMock) Save(key string, value []byte, ttlDays ...int) (*persistentstorage.ObjectInfo, error) {
	_va := make([]interface{}, len(ttlDays))
	for _i := range ttlDays {
		_va[_i] = ttlDays[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, key, value)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *persistentstorage.ObjectInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(string, []byte, ...int) (*persistentstorage.ObjectInfo, error)); ok {
		return rf(key, value, ttlDays...)
	}
	if rf, ok := ret.Get(0).(func(string, []byte, ...int) *persistentstorage.ObjectInfo); ok {
		r0 = rf(key, value, ttlDays...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*persistentstorage.ObjectInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(string, []byte, ...int) error); ok {
		r1 = rf(key, value, ttlDays...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PersistentStorageMock_Save_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Save'
type PersistentStorageMock_Save_Call struct {
	*mock.Call
}

// Save is a helper method to define mock.On call
//   - key string
//   - value []byte
//   - ttlDays ...int
func (_e *PersistentStorageMock_Expecter) Save(key interface{}, value interface{}, ttlDays ...interface{}) *PersistentStorageMock_Save_Call {
	return &PersistentStorageMock_Save_Call{Call: _e.mock.On("Save",
		append([]interface{}{key, value}, ttlDays...)...)}
}

func (_c *PersistentStorageMock_Save_Call) Run(run func(key string, value []byte, ttlDays ...int)) *PersistentStorageMock_Save_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]int, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(int)
			}
		}
		run(args[0].(string), args[1].([]byte), variadicArgs...)
	})
	return _c
}

func (_c *PersistentStorageMock_Save_Call) Return(_a0 *persistentstorage.ObjectInfo, _a1 error) *PersistentStorageMock_Save_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PersistentStorageMock_Save_Call) RunAndReturn(run func(string, []byte, ...int) (*persistentstorage.ObjectInfo, error)) *PersistentStorageMock_Save_Call {
	_c.Call.Return(run)
	return _c
}

// NewPersistentStorageMock creates a new instance of PersistentStorageMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPersistentStorageMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *PersistentStorageMock {
	mock := &PersistentStorageMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
