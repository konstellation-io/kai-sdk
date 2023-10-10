// Code generated by mockery v2.33.1. DO NOT EDIT.

package mocks

import (
	io "io"

	nats "github.com/nats-io/nats.go"
	mock "github.com/stretchr/testify/mock"
)

// NatsObjectStoreMock is an autogenerated mock type for the ObjectStore type
type NatsObjectStoreMock struct {
	mock.Mock
}

type NatsObjectStoreMock_Expecter struct {
	mock *mock.Mock
}

func (_m *NatsObjectStoreMock) EXPECT() *NatsObjectStoreMock_Expecter {
	return &NatsObjectStoreMock_Expecter{mock: &_m.Mock}
}

// AddBucketLink provides a mock function with given fields: name, bucket
func (_m *NatsObjectStoreMock) AddBucketLink(name string, bucket nats.ObjectStore) (*nats.ObjectInfo, error) {
	ret := _m.Called(name, bucket)

	var r0 *nats.ObjectInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(string, nats.ObjectStore) (*nats.ObjectInfo, error)); ok {
		return rf(name, bucket)
	}
	if rf, ok := ret.Get(0).(func(string, nats.ObjectStore) *nats.ObjectInfo); ok {
		r0 = rf(name, bucket)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*nats.ObjectInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(string, nats.ObjectStore) error); ok {
		r1 = rf(name, bucket)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NatsObjectStoreMock_AddBucketLink_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddBucketLink'
type NatsObjectStoreMock_AddBucketLink_Call struct {
	*mock.Call
}

// AddBucketLink is a helper method to define mock.On call
//   - name string
//   - bucket nats.ObjectStore
func (_e *NatsObjectStoreMock_Expecter) AddBucketLink(name interface{}, bucket interface{}) *NatsObjectStoreMock_AddBucketLink_Call {
	return &NatsObjectStoreMock_AddBucketLink_Call{Call: _e.mock.On("AddBucketLink", name, bucket)}
}

func (_c *NatsObjectStoreMock_AddBucketLink_Call) Run(run func(name string, bucket nats.ObjectStore)) *NatsObjectStoreMock_AddBucketLink_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(nats.ObjectStore))
	})
	return _c
}

func (_c *NatsObjectStoreMock_AddBucketLink_Call) Return(_a0 *nats.ObjectInfo, _a1 error) *NatsObjectStoreMock_AddBucketLink_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *NatsObjectStoreMock_AddBucketLink_Call) RunAndReturn(run func(string, nats.ObjectStore) (*nats.ObjectInfo, error)) *NatsObjectStoreMock_AddBucketLink_Call {
	_c.Call.Return(run)
	return _c
}

// AddLink provides a mock function with given fields: name, obj
func (_m *NatsObjectStoreMock) AddLink(name string, obj *nats.ObjectInfo) (*nats.ObjectInfo, error) {
	ret := _m.Called(name, obj)

	var r0 *nats.ObjectInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(string, *nats.ObjectInfo) (*nats.ObjectInfo, error)); ok {
		return rf(name, obj)
	}
	if rf, ok := ret.Get(0).(func(string, *nats.ObjectInfo) *nats.ObjectInfo); ok {
		r0 = rf(name, obj)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*nats.ObjectInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(string, *nats.ObjectInfo) error); ok {
		r1 = rf(name, obj)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NatsObjectStoreMock_AddLink_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddLink'
type NatsObjectStoreMock_AddLink_Call struct {
	*mock.Call
}

// AddLink is a helper method to define mock.On call
//   - name string
//   - obj *nats.ObjectInfo
func (_e *NatsObjectStoreMock_Expecter) AddLink(name interface{}, obj interface{}) *NatsObjectStoreMock_AddLink_Call {
	return &NatsObjectStoreMock_AddLink_Call{Call: _e.mock.On("AddLink", name, obj)}
}

func (_c *NatsObjectStoreMock_AddLink_Call) Run(run func(name string, obj *nats.ObjectInfo)) *NatsObjectStoreMock_AddLink_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(*nats.ObjectInfo))
	})
	return _c
}

func (_c *NatsObjectStoreMock_AddLink_Call) Return(_a0 *nats.ObjectInfo, _a1 error) *NatsObjectStoreMock_AddLink_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *NatsObjectStoreMock_AddLink_Call) RunAndReturn(run func(string, *nats.ObjectInfo) (*nats.ObjectInfo, error)) *NatsObjectStoreMock_AddLink_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: name
func (_m *NatsObjectStoreMock) Delete(name string) error {
	ret := _m.Called(name)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NatsObjectStoreMock_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type NatsObjectStoreMock_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - name string
func (_e *NatsObjectStoreMock_Expecter) Delete(name interface{}) *NatsObjectStoreMock_Delete_Call {
	return &NatsObjectStoreMock_Delete_Call{Call: _e.mock.On("Delete", name)}
}

func (_c *NatsObjectStoreMock_Delete_Call) Run(run func(name string)) *NatsObjectStoreMock_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *NatsObjectStoreMock_Delete_Call) Return(_a0 error) *NatsObjectStoreMock_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *NatsObjectStoreMock_Delete_Call) RunAndReturn(run func(string) error) *NatsObjectStoreMock_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: name, opts
func (_m *NatsObjectStoreMock) Get(name string, opts ...nats.GetObjectOpt) (nats.ObjectResult, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, name)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 nats.ObjectResult
	var r1 error
	if rf, ok := ret.Get(0).(func(string, ...nats.GetObjectOpt) (nats.ObjectResult, error)); ok {
		return rf(name, opts...)
	}
	if rf, ok := ret.Get(0).(func(string, ...nats.GetObjectOpt) nats.ObjectResult); ok {
		r0 = rf(name, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(nats.ObjectResult)
		}
	}

	if rf, ok := ret.Get(1).(func(string, ...nats.GetObjectOpt) error); ok {
		r1 = rf(name, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NatsObjectStoreMock_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type NatsObjectStoreMock_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - name string
//   - opts ...nats.GetObjectOpt
func (_e *NatsObjectStoreMock_Expecter) Get(name interface{}, opts ...interface{}) *NatsObjectStoreMock_Get_Call {
	return &NatsObjectStoreMock_Get_Call{Call: _e.mock.On("Get",
		append([]interface{}{name}, opts...)...)}
}

func (_c *NatsObjectStoreMock_Get_Call) Run(run func(name string, opts ...nats.GetObjectOpt)) *NatsObjectStoreMock_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]nats.GetObjectOpt, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(nats.GetObjectOpt)
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *NatsObjectStoreMock_Get_Call) Return(_a0 nats.ObjectResult, _a1 error) *NatsObjectStoreMock_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *NatsObjectStoreMock_Get_Call) RunAndReturn(run func(string, ...nats.GetObjectOpt) (nats.ObjectResult, error)) *NatsObjectStoreMock_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetBytes provides a mock function with given fields: name, opts
func (_m *NatsObjectStoreMock) GetBytes(name string, opts ...nats.GetObjectOpt) ([]byte, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, name)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(string, ...nats.GetObjectOpt) ([]byte, error)); ok {
		return rf(name, opts...)
	}
	if rf, ok := ret.Get(0).(func(string, ...nats.GetObjectOpt) []byte); ok {
		r0 = rf(name, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(string, ...nats.GetObjectOpt) error); ok {
		r1 = rf(name, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NatsObjectStoreMock_GetBytes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetBytes'
type NatsObjectStoreMock_GetBytes_Call struct {
	*mock.Call
}

// GetBytes is a helper method to define mock.On call
//   - name string
//   - opts ...nats.GetObjectOpt
func (_e *NatsObjectStoreMock_Expecter) GetBytes(name interface{}, opts ...interface{}) *NatsObjectStoreMock_GetBytes_Call {
	return &NatsObjectStoreMock_GetBytes_Call{Call: _e.mock.On("GetBytes",
		append([]interface{}{name}, opts...)...)}
}

func (_c *NatsObjectStoreMock_GetBytes_Call) Run(run func(name string, opts ...nats.GetObjectOpt)) *NatsObjectStoreMock_GetBytes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]nats.GetObjectOpt, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(nats.GetObjectOpt)
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *NatsObjectStoreMock_GetBytes_Call) Return(_a0 []byte, _a1 error) *NatsObjectStoreMock_GetBytes_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *NatsObjectStoreMock_GetBytes_Call) RunAndReturn(run func(string, ...nats.GetObjectOpt) ([]byte, error)) *NatsObjectStoreMock_GetBytes_Call {
	_c.Call.Return(run)
	return _c
}

// GetFile provides a mock function with given fields: name, file, opts
func (_m *NatsObjectStoreMock) GetFile(name string, file string, opts ...nats.GetObjectOpt) error {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, name, file)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, ...nats.GetObjectOpt) error); ok {
		r0 = rf(name, file, opts...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NatsObjectStoreMock_GetFile_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetFile'
type NatsObjectStoreMock_GetFile_Call struct {
	*mock.Call
}

// GetFile is a helper method to define mock.On call
//   - name string
//   - file string
//   - opts ...nats.GetObjectOpt
func (_e *NatsObjectStoreMock_Expecter) GetFile(name interface{}, file interface{}, opts ...interface{}) *NatsObjectStoreMock_GetFile_Call {
	return &NatsObjectStoreMock_GetFile_Call{Call: _e.mock.On("GetFile",
		append([]interface{}{name, file}, opts...)...)}
}

func (_c *NatsObjectStoreMock_GetFile_Call) Run(run func(name string, file string, opts ...nats.GetObjectOpt)) *NatsObjectStoreMock_GetFile_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]nats.GetObjectOpt, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(nats.GetObjectOpt)
			}
		}
		run(args[0].(string), args[1].(string), variadicArgs...)
	})
	return _c
}

func (_c *NatsObjectStoreMock_GetFile_Call) Return(_a0 error) *NatsObjectStoreMock_GetFile_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *NatsObjectStoreMock_GetFile_Call) RunAndReturn(run func(string, string, ...nats.GetObjectOpt) error) *NatsObjectStoreMock_GetFile_Call {
	_c.Call.Return(run)
	return _c
}

// GetInfo provides a mock function with given fields: name, opts
func (_m *NatsObjectStoreMock) GetInfo(name string, opts ...nats.GetObjectInfoOpt) (*nats.ObjectInfo, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, name)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *nats.ObjectInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(string, ...nats.GetObjectInfoOpt) (*nats.ObjectInfo, error)); ok {
		return rf(name, opts...)
	}
	if rf, ok := ret.Get(0).(func(string, ...nats.GetObjectInfoOpt) *nats.ObjectInfo); ok {
		r0 = rf(name, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*nats.ObjectInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(string, ...nats.GetObjectInfoOpt) error); ok {
		r1 = rf(name, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NatsObjectStoreMock_GetInfo_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetInfo'
type NatsObjectStoreMock_GetInfo_Call struct {
	*mock.Call
}

// GetInfo is a helper method to define mock.On call
//   - name string
//   - opts ...nats.GetObjectInfoOpt
func (_e *NatsObjectStoreMock_Expecter) GetInfo(name interface{}, opts ...interface{}) *NatsObjectStoreMock_GetInfo_Call {
	return &NatsObjectStoreMock_GetInfo_Call{Call: _e.mock.On("GetInfo",
		append([]interface{}{name}, opts...)...)}
}

func (_c *NatsObjectStoreMock_GetInfo_Call) Run(run func(name string, opts ...nats.GetObjectInfoOpt)) *NatsObjectStoreMock_GetInfo_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]nats.GetObjectInfoOpt, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(nats.GetObjectInfoOpt)
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *NatsObjectStoreMock_GetInfo_Call) Return(_a0 *nats.ObjectInfo, _a1 error) *NatsObjectStoreMock_GetInfo_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *NatsObjectStoreMock_GetInfo_Call) RunAndReturn(run func(string, ...nats.GetObjectInfoOpt) (*nats.ObjectInfo, error)) *NatsObjectStoreMock_GetInfo_Call {
	_c.Call.Return(run)
	return _c
}

// GetString provides a mock function with given fields: name, opts
func (_m *NatsObjectStoreMock) GetString(name string, opts ...nats.GetObjectOpt) (string, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, name)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string, ...nats.GetObjectOpt) (string, error)); ok {
		return rf(name, opts...)
	}
	if rf, ok := ret.Get(0).(func(string, ...nats.GetObjectOpt) string); ok {
		r0 = rf(name, opts...)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string, ...nats.GetObjectOpt) error); ok {
		r1 = rf(name, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NatsObjectStoreMock_GetString_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetString'
type NatsObjectStoreMock_GetString_Call struct {
	*mock.Call
}

// GetString is a helper method to define mock.On call
//   - name string
//   - opts ...nats.GetObjectOpt
func (_e *NatsObjectStoreMock_Expecter) GetString(name interface{}, opts ...interface{}) *NatsObjectStoreMock_GetString_Call {
	return &NatsObjectStoreMock_GetString_Call{Call: _e.mock.On("GetString",
		append([]interface{}{name}, opts...)...)}
}

func (_c *NatsObjectStoreMock_GetString_Call) Run(run func(name string, opts ...nats.GetObjectOpt)) *NatsObjectStoreMock_GetString_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]nats.GetObjectOpt, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(nats.GetObjectOpt)
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *NatsObjectStoreMock_GetString_Call) Return(_a0 string, _a1 error) *NatsObjectStoreMock_GetString_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *NatsObjectStoreMock_GetString_Call) RunAndReturn(run func(string, ...nats.GetObjectOpt) (string, error)) *NatsObjectStoreMock_GetString_Call {
	_c.Call.Return(run)
	return _c
}

// List provides a mock function with given fields: opts
func (_m *NatsObjectStoreMock) List(opts ...nats.ListObjectsOpt) ([]*nats.ObjectInfo, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []*nats.ObjectInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(...nats.ListObjectsOpt) ([]*nats.ObjectInfo, error)); ok {
		return rf(opts...)
	}
	if rf, ok := ret.Get(0).(func(...nats.ListObjectsOpt) []*nats.ObjectInfo); ok {
		r0 = rf(opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*nats.ObjectInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(...nats.ListObjectsOpt) error); ok {
		r1 = rf(opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NatsObjectStoreMock_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type NatsObjectStoreMock_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - opts ...nats.ListObjectsOpt
func (_e *NatsObjectStoreMock_Expecter) List(opts ...interface{}) *NatsObjectStoreMock_List_Call {
	return &NatsObjectStoreMock_List_Call{Call: _e.mock.On("List",
		append([]interface{}{}, opts...)...)}
}

func (_c *NatsObjectStoreMock_List_Call) Run(run func(opts ...nats.ListObjectsOpt)) *NatsObjectStoreMock_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]nats.ListObjectsOpt, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(nats.ListObjectsOpt)
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *NatsObjectStoreMock_List_Call) Return(_a0 []*nats.ObjectInfo, _a1 error) *NatsObjectStoreMock_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *NatsObjectStoreMock_List_Call) RunAndReturn(run func(...nats.ListObjectsOpt) ([]*nats.ObjectInfo, error)) *NatsObjectStoreMock_List_Call {
	_c.Call.Return(run)
	return _c
}

// Put provides a mock function with given fields: obj, reader, opts
func (_m *NatsObjectStoreMock) Put(obj *nats.ObjectMeta, reader io.Reader, opts ...nats.ObjectOpt) (*nats.ObjectInfo, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, obj, reader)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *nats.ObjectInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(*nats.ObjectMeta, io.Reader, ...nats.ObjectOpt) (*nats.ObjectInfo, error)); ok {
		return rf(obj, reader, opts...)
	}
	if rf, ok := ret.Get(0).(func(*nats.ObjectMeta, io.Reader, ...nats.ObjectOpt) *nats.ObjectInfo); ok {
		r0 = rf(obj, reader, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*nats.ObjectInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(*nats.ObjectMeta, io.Reader, ...nats.ObjectOpt) error); ok {
		r1 = rf(obj, reader, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NatsObjectStoreMock_Put_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Put'
type NatsObjectStoreMock_Put_Call struct {
	*mock.Call
}

// Put is a helper method to define mock.On call
//   - obj *nats.ObjectMeta
//   - reader io.Reader
//   - opts ...nats.ObjectOpt
func (_e *NatsObjectStoreMock_Expecter) Put(obj interface{}, reader interface{}, opts ...interface{}) *NatsObjectStoreMock_Put_Call {
	return &NatsObjectStoreMock_Put_Call{Call: _e.mock.On("Put",
		append([]interface{}{obj, reader}, opts...)...)}
}

func (_c *NatsObjectStoreMock_Put_Call) Run(run func(obj *nats.ObjectMeta, reader io.Reader, opts ...nats.ObjectOpt)) *NatsObjectStoreMock_Put_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]nats.ObjectOpt, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(nats.ObjectOpt)
			}
		}
		run(args[0].(*nats.ObjectMeta), args[1].(io.Reader), variadicArgs...)
	})
	return _c
}

func (_c *NatsObjectStoreMock_Put_Call) Return(_a0 *nats.ObjectInfo, _a1 error) *NatsObjectStoreMock_Put_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *NatsObjectStoreMock_Put_Call) RunAndReturn(run func(*nats.ObjectMeta, io.Reader, ...nats.ObjectOpt) (*nats.ObjectInfo, error)) *NatsObjectStoreMock_Put_Call {
	_c.Call.Return(run)
	return _c
}

// PutBytes provides a mock function with given fields: name, data, opts
func (_m *NatsObjectStoreMock) PutBytes(name string, data []byte, opts ...nats.ObjectOpt) (*nats.ObjectInfo, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, name, data)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *nats.ObjectInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(string, []byte, ...nats.ObjectOpt) (*nats.ObjectInfo, error)); ok {
		return rf(name, data, opts...)
	}
	if rf, ok := ret.Get(0).(func(string, []byte, ...nats.ObjectOpt) *nats.ObjectInfo); ok {
		r0 = rf(name, data, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*nats.ObjectInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(string, []byte, ...nats.ObjectOpt) error); ok {
		r1 = rf(name, data, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NatsObjectStoreMock_PutBytes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PutBytes'
type NatsObjectStoreMock_PutBytes_Call struct {
	*mock.Call
}

// PutBytes is a helper method to define mock.On call
//   - name string
//   - data []byte
//   - opts ...nats.ObjectOpt
func (_e *NatsObjectStoreMock_Expecter) PutBytes(name interface{}, data interface{}, opts ...interface{}) *NatsObjectStoreMock_PutBytes_Call {
	return &NatsObjectStoreMock_PutBytes_Call{Call: _e.mock.On("PutBytes",
		append([]interface{}{name, data}, opts...)...)}
}

func (_c *NatsObjectStoreMock_PutBytes_Call) Run(run func(name string, data []byte, opts ...nats.ObjectOpt)) *NatsObjectStoreMock_PutBytes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]nats.ObjectOpt, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(nats.ObjectOpt)
			}
		}
		run(args[0].(string), args[1].([]byte), variadicArgs...)
	})
	return _c
}

func (_c *NatsObjectStoreMock_PutBytes_Call) Return(_a0 *nats.ObjectInfo, _a1 error) *NatsObjectStoreMock_PutBytes_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *NatsObjectStoreMock_PutBytes_Call) RunAndReturn(run func(string, []byte, ...nats.ObjectOpt) (*nats.ObjectInfo, error)) *NatsObjectStoreMock_PutBytes_Call {
	_c.Call.Return(run)
	return _c
}

// PutFile provides a mock function with given fields: file, opts
func (_m *NatsObjectStoreMock) PutFile(file string, opts ...nats.ObjectOpt) (*nats.ObjectInfo, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, file)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *nats.ObjectInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(string, ...nats.ObjectOpt) (*nats.ObjectInfo, error)); ok {
		return rf(file, opts...)
	}
	if rf, ok := ret.Get(0).(func(string, ...nats.ObjectOpt) *nats.ObjectInfo); ok {
		r0 = rf(file, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*nats.ObjectInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(string, ...nats.ObjectOpt) error); ok {
		r1 = rf(file, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NatsObjectStoreMock_PutFile_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PutFile'
type NatsObjectStoreMock_PutFile_Call struct {
	*mock.Call
}

// PutFile is a helper method to define mock.On call
//   - file string
//   - opts ...nats.ObjectOpt
func (_e *NatsObjectStoreMock_Expecter) PutFile(file interface{}, opts ...interface{}) *NatsObjectStoreMock_PutFile_Call {
	return &NatsObjectStoreMock_PutFile_Call{Call: _e.mock.On("PutFile",
		append([]interface{}{file}, opts...)...)}
}

func (_c *NatsObjectStoreMock_PutFile_Call) Run(run func(file string, opts ...nats.ObjectOpt)) *NatsObjectStoreMock_PutFile_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]nats.ObjectOpt, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(nats.ObjectOpt)
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *NatsObjectStoreMock_PutFile_Call) Return(_a0 *nats.ObjectInfo, _a1 error) *NatsObjectStoreMock_PutFile_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *NatsObjectStoreMock_PutFile_Call) RunAndReturn(run func(string, ...nats.ObjectOpt) (*nats.ObjectInfo, error)) *NatsObjectStoreMock_PutFile_Call {
	_c.Call.Return(run)
	return _c
}

// PutString provides a mock function with given fields: name, data, opts
func (_m *NatsObjectStoreMock) PutString(name string, data string, opts ...nats.ObjectOpt) (*nats.ObjectInfo, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, name, data)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *nats.ObjectInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string, ...nats.ObjectOpt) (*nats.ObjectInfo, error)); ok {
		return rf(name, data, opts...)
	}
	if rf, ok := ret.Get(0).(func(string, string, ...nats.ObjectOpt) *nats.ObjectInfo); ok {
		r0 = rf(name, data, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*nats.ObjectInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string, ...nats.ObjectOpt) error); ok {
		r1 = rf(name, data, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NatsObjectStoreMock_PutString_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PutString'
type NatsObjectStoreMock_PutString_Call struct {
	*mock.Call
}

// PutString is a helper method to define mock.On call
//   - name string
//   - data string
//   - opts ...nats.ObjectOpt
func (_e *NatsObjectStoreMock_Expecter) PutString(name interface{}, data interface{}, opts ...interface{}) *NatsObjectStoreMock_PutString_Call {
	return &NatsObjectStoreMock_PutString_Call{Call: _e.mock.On("PutString",
		append([]interface{}{name, data}, opts...)...)}
}

func (_c *NatsObjectStoreMock_PutString_Call) Run(run func(name string, data string, opts ...nats.ObjectOpt)) *NatsObjectStoreMock_PutString_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]nats.ObjectOpt, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(nats.ObjectOpt)
			}
		}
		run(args[0].(string), args[1].(string), variadicArgs...)
	})
	return _c
}

func (_c *NatsObjectStoreMock_PutString_Call) Return(_a0 *nats.ObjectInfo, _a1 error) *NatsObjectStoreMock_PutString_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *NatsObjectStoreMock_PutString_Call) RunAndReturn(run func(string, string, ...nats.ObjectOpt) (*nats.ObjectInfo, error)) *NatsObjectStoreMock_PutString_Call {
	_c.Call.Return(run)
	return _c
}

// Seal provides a mock function with given fields:
func (_m *NatsObjectStoreMock) Seal() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NatsObjectStoreMock_Seal_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Seal'
type NatsObjectStoreMock_Seal_Call struct {
	*mock.Call
}

// Seal is a helper method to define mock.On call
func (_e *NatsObjectStoreMock_Expecter) Seal() *NatsObjectStoreMock_Seal_Call {
	return &NatsObjectStoreMock_Seal_Call{Call: _e.mock.On("Seal")}
}

func (_c *NatsObjectStoreMock_Seal_Call) Run(run func()) *NatsObjectStoreMock_Seal_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *NatsObjectStoreMock_Seal_Call) Return(_a0 error) *NatsObjectStoreMock_Seal_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *NatsObjectStoreMock_Seal_Call) RunAndReturn(run func() error) *NatsObjectStoreMock_Seal_Call {
	_c.Call.Return(run)
	return _c
}

// Status provides a mock function with given fields:
func (_m *NatsObjectStoreMock) Status() (nats.ObjectStoreStatus, error) {
	ret := _m.Called()

	var r0 nats.ObjectStoreStatus
	var r1 error
	if rf, ok := ret.Get(0).(func() (nats.ObjectStoreStatus, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() nats.ObjectStoreStatus); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(nats.ObjectStoreStatus)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NatsObjectStoreMock_Status_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Status'
type NatsObjectStoreMock_Status_Call struct {
	*mock.Call
}

// Status is a helper method to define mock.On call
func (_e *NatsObjectStoreMock_Expecter) Status() *NatsObjectStoreMock_Status_Call {
	return &NatsObjectStoreMock_Status_Call{Call: _e.mock.On("Status")}
}

func (_c *NatsObjectStoreMock_Status_Call) Run(run func()) *NatsObjectStoreMock_Status_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *NatsObjectStoreMock_Status_Call) Return(_a0 nats.ObjectStoreStatus, _a1 error) *NatsObjectStoreMock_Status_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *NatsObjectStoreMock_Status_Call) RunAndReturn(run func() (nats.ObjectStoreStatus, error)) *NatsObjectStoreMock_Status_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateMeta provides a mock function with given fields: name, meta
func (_m *NatsObjectStoreMock) UpdateMeta(name string, meta *nats.ObjectMeta) error {
	ret := _m.Called(name, meta)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, *nats.ObjectMeta) error); ok {
		r0 = rf(name, meta)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NatsObjectStoreMock_UpdateMeta_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateMeta'
type NatsObjectStoreMock_UpdateMeta_Call struct {
	*mock.Call
}

// UpdateMeta is a helper method to define mock.On call
//   - name string
//   - meta *nats.ObjectMeta
func (_e *NatsObjectStoreMock_Expecter) UpdateMeta(name interface{}, meta interface{}) *NatsObjectStoreMock_UpdateMeta_Call {
	return &NatsObjectStoreMock_UpdateMeta_Call{Call: _e.mock.On("UpdateMeta", name, meta)}
}

func (_c *NatsObjectStoreMock_UpdateMeta_Call) Run(run func(name string, meta *nats.ObjectMeta)) *NatsObjectStoreMock_UpdateMeta_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(*nats.ObjectMeta))
	})
	return _c
}

func (_c *NatsObjectStoreMock_UpdateMeta_Call) Return(_a0 error) *NatsObjectStoreMock_UpdateMeta_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *NatsObjectStoreMock_UpdateMeta_Call) RunAndReturn(run func(string, *nats.ObjectMeta) error) *NatsObjectStoreMock_UpdateMeta_Call {
	_c.Call.Return(run)
	return _c
}

// Watch provides a mock function with given fields: opts
func (_m *NatsObjectStoreMock) Watch(opts ...nats.WatchOpt) (nats.ObjectWatcher, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 nats.ObjectWatcher
	var r1 error
	if rf, ok := ret.Get(0).(func(...nats.WatchOpt) (nats.ObjectWatcher, error)); ok {
		return rf(opts...)
	}
	if rf, ok := ret.Get(0).(func(...nats.WatchOpt) nats.ObjectWatcher); ok {
		r0 = rf(opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(nats.ObjectWatcher)
		}
	}

	if rf, ok := ret.Get(1).(func(...nats.WatchOpt) error); ok {
		r1 = rf(opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NatsObjectStoreMock_Watch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Watch'
type NatsObjectStoreMock_Watch_Call struct {
	*mock.Call
}

// Watch is a helper method to define mock.On call
//   - opts ...nats.WatchOpt
func (_e *NatsObjectStoreMock_Expecter) Watch(opts ...interface{}) *NatsObjectStoreMock_Watch_Call {
	return &NatsObjectStoreMock_Watch_Call{Call: _e.mock.On("Watch",
		append([]interface{}{}, opts...)...)}
}

func (_c *NatsObjectStoreMock_Watch_Call) Run(run func(opts ...nats.WatchOpt)) *NatsObjectStoreMock_Watch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]nats.WatchOpt, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(nats.WatchOpt)
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *NatsObjectStoreMock_Watch_Call) Return(_a0 nats.ObjectWatcher, _a1 error) *NatsObjectStoreMock_Watch_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *NatsObjectStoreMock_Watch_Call) RunAndReturn(run func(...nats.WatchOpt) (nats.ObjectWatcher, error)) *NatsObjectStoreMock_Watch_Call {
	_c.Call.Return(run)
	return _c
}

// NewNatsObjectStoreMock creates a new instance of NatsObjectStoreMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewNatsObjectStoreMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *NatsObjectStoreMock {
	mock := &NatsObjectStoreMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
