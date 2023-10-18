// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// MetadataMock is an autogenerated mock type for the metadata type
type MetadataMock struct {
	mock.Mock
}

type MetadataMock_Expecter struct {
	mock *mock.Mock
}

func (_m *MetadataMock) EXPECT() *MetadataMock_Expecter {
	return &MetadataMock_Expecter{mock: &_m.Mock}
}

// GetKeyValueStoreProcessName provides a mock function with given fields:
func (_m *MetadataMock) GetKeyValueStoreProcessName() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MetadataMock_GetKeyValueStoreProcessName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetKeyValueStoreProcessName'
type MetadataMock_GetKeyValueStoreProcessName_Call struct {
	*mock.Call
}

// GetKeyValueStoreProcessName is a helper method to define mock.On call
func (_e *MetadataMock_Expecter) GetKeyValueStoreProcessName() *MetadataMock_GetKeyValueStoreProcessName_Call {
	return &MetadataMock_GetKeyValueStoreProcessName_Call{Call: _e.mock.On("GetKeyValueStoreProcessName")}
}

func (_c *MetadataMock_GetKeyValueStoreProcessName_Call) Run(run func()) *MetadataMock_GetKeyValueStoreProcessName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MetadataMock_GetKeyValueStoreProcessName_Call) Return(_a0 string) *MetadataMock_GetKeyValueStoreProcessName_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MetadataMock_GetKeyValueStoreProcessName_Call) RunAndReturn(run func() string) *MetadataMock_GetKeyValueStoreProcessName_Call {
	_c.Call.Return(run)
	return _c
}

// GetKeyValueStoreProductName provides a mock function with given fields:
func (_m *MetadataMock) GetKeyValueStoreProductName() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MetadataMock_GetKeyValueStoreProductName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetKeyValueStoreProductName'
type MetadataMock_GetKeyValueStoreProductName_Call struct {
	*mock.Call
}

// GetKeyValueStoreProductName is a helper method to define mock.On call
func (_e *MetadataMock_Expecter) GetKeyValueStoreProductName() *MetadataMock_GetKeyValueStoreProductName_Call {
	return &MetadataMock_GetKeyValueStoreProductName_Call{Call: _e.mock.On("GetKeyValueStoreProductName")}
}

func (_c *MetadataMock_GetKeyValueStoreProductName_Call) Run(run func()) *MetadataMock_GetKeyValueStoreProductName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MetadataMock_GetKeyValueStoreProductName_Call) Return(_a0 string) *MetadataMock_GetKeyValueStoreProductName_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MetadataMock_GetKeyValueStoreProductName_Call) RunAndReturn(run func() string) *MetadataMock_GetKeyValueStoreProductName_Call {
	_c.Call.Return(run)
	return _c
}

// GetKeyValueStoreWorkflowName provides a mock function with given fields:
func (_m *MetadataMock) GetKeyValueStoreWorkflowName() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MetadataMock_GetKeyValueStoreWorkflowName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetKeyValueStoreWorkflowName'
type MetadataMock_GetKeyValueStoreWorkflowName_Call struct {
	*mock.Call
}

// GetKeyValueStoreWorkflowName is a helper method to define mock.On call
func (_e *MetadataMock_Expecter) GetKeyValueStoreWorkflowName() *MetadataMock_GetKeyValueStoreWorkflowName_Call {
	return &MetadataMock_GetKeyValueStoreWorkflowName_Call{Call: _e.mock.On("GetKeyValueStoreWorkflowName")}
}

func (_c *MetadataMock_GetKeyValueStoreWorkflowName_Call) Run(run func()) *MetadataMock_GetKeyValueStoreWorkflowName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MetadataMock_GetKeyValueStoreWorkflowName_Call) Return(_a0 string) *MetadataMock_GetKeyValueStoreWorkflowName_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MetadataMock_GetKeyValueStoreWorkflowName_Call) RunAndReturn(run func() string) *MetadataMock_GetKeyValueStoreWorkflowName_Call {
	_c.Call.Return(run)
	return _c
}

// GetObjectStoreName provides a mock function with given fields:
func (_m *MetadataMock) GetObjectStoreName() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MetadataMock_GetObjectStoreName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetObjectStoreName'
type MetadataMock_GetObjectStoreName_Call struct {
	*mock.Call
}

// GetObjectStoreName is a helper method to define mock.On call
func (_e *MetadataMock_Expecter) GetObjectStoreName() *MetadataMock_GetObjectStoreName_Call {
	return &MetadataMock_GetObjectStoreName_Call{Call: _e.mock.On("GetObjectStoreName")}
}

func (_c *MetadataMock_GetObjectStoreName_Call) Run(run func()) *MetadataMock_GetObjectStoreName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MetadataMock_GetObjectStoreName_Call) Return(_a0 string) *MetadataMock_GetObjectStoreName_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MetadataMock_GetObjectStoreName_Call) RunAndReturn(run func() string) *MetadataMock_GetObjectStoreName_Call {
	_c.Call.Return(run)
	return _c
}

// GetProcess provides a mock function with given fields:
func (_m *MetadataMock) GetProcess() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MetadataMock_GetProcess_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetProcess'
type MetadataMock_GetProcess_Call struct {
	*mock.Call
}

// GetProcess is a helper method to define mock.On call
func (_e *MetadataMock_Expecter) GetProcess() *MetadataMock_GetProcess_Call {
	return &MetadataMock_GetProcess_Call{Call: _e.mock.On("GetProcess")}
}

func (_c *MetadataMock_GetProcess_Call) Run(run func()) *MetadataMock_GetProcess_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MetadataMock_GetProcess_Call) Return(_a0 string) *MetadataMock_GetProcess_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MetadataMock_GetProcess_Call) RunAndReturn(run func() string) *MetadataMock_GetProcess_Call {
	_c.Call.Return(run)
	return _c
}

// GetProduct provides a mock function with given fields:
func (_m *MetadataMock) GetProduct() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MetadataMock_GetProduct_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetProduct'
type MetadataMock_GetProduct_Call struct {
	*mock.Call
}

// GetProduct is a helper method to define mock.On call
func (_e *MetadataMock_Expecter) GetProduct() *MetadataMock_GetProduct_Call {
	return &MetadataMock_GetProduct_Call{Call: _e.mock.On("GetProduct")}
}

func (_c *MetadataMock_GetProduct_Call) Run(run func()) *MetadataMock_GetProduct_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MetadataMock_GetProduct_Call) Return(_a0 string) *MetadataMock_GetProduct_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MetadataMock_GetProduct_Call) RunAndReturn(run func() string) *MetadataMock_GetProduct_Call {
	_c.Call.Return(run)
	return _c
}

// GetVersion provides a mock function with given fields:
func (_m *MetadataMock) GetVersion() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MetadataMock_GetVersion_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetVersion'
type MetadataMock_GetVersion_Call struct {
	*mock.Call
}

// GetVersion is a helper method to define mock.On call
func (_e *MetadataMock_Expecter) GetVersion() *MetadataMock_GetVersion_Call {
	return &MetadataMock_GetVersion_Call{Call: _e.mock.On("GetVersion")}
}

func (_c *MetadataMock_GetVersion_Call) Run(run func()) *MetadataMock_GetVersion_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MetadataMock_GetVersion_Call) Return(_a0 string) *MetadataMock_GetVersion_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MetadataMock_GetVersion_Call) RunAndReturn(run func() string) *MetadataMock_GetVersion_Call {
	_c.Call.Return(run)
	return _c
}

// GetWorkflow provides a mock function with given fields:
func (_m *MetadataMock) GetWorkflow() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MetadataMock_GetWorkflow_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetWorkflow'
type MetadataMock_GetWorkflow_Call struct {
	*mock.Call
}

// GetWorkflow is a helper method to define mock.On call
func (_e *MetadataMock_Expecter) GetWorkflow() *MetadataMock_GetWorkflow_Call {
	return &MetadataMock_GetWorkflow_Call{Call: _e.mock.On("GetWorkflow")}
}

func (_c *MetadataMock_GetWorkflow_Call) Run(run func()) *MetadataMock_GetWorkflow_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MetadataMock_GetWorkflow_Call) Return(_a0 string) *MetadataMock_GetWorkflow_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MetadataMock_GetWorkflow_Call) RunAndReturn(run func() string) *MetadataMock_GetWorkflow_Call {
	_c.Call.Return(run)
	return _c
}

// NewMetadataMock creates a new instance of MetadataMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMetadataMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *MetadataMock {
	mock := &MetadataMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
