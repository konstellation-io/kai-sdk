// Code generated by mockery v2.32.0. DO NOT EDIT.

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

// GetEphemeralStorageName provides a mock function with given fields:
func (_m *MetadataMock) GetEphemeralStorageName() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MetadataMock_GetEphemeralStorageName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetEphemeralStorageName'
type MetadataMock_GetEphemeralStorageName_Call struct {
	*mock.Call
}

// GetEphemeralStorageName is a helper method to define mock.On call
func (_e *MetadataMock_Expecter) GetEphemeralStorageName() *MetadataMock_GetEphemeralStorageName_Call {
	return &MetadataMock_GetEphemeralStorageName_Call{Call: _e.mock.On("GetEphemeralStorageName")}
}

func (_c *MetadataMock_GetEphemeralStorageName_Call) Run(run func()) *MetadataMock_GetEphemeralStorageName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MetadataMock_GetEphemeralStorageName_Call) Return(_a0 string) *MetadataMock_GetEphemeralStorageName_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MetadataMock_GetEphemeralStorageName_Call) RunAndReturn(run func() string) *MetadataMock_GetEphemeralStorageName_Call {
	_c.Call.Return(run)
	return _c
}

// GetGlobalCentralizedConfigurationName provides a mock function with given fields:
func (_m *MetadataMock) GetGlobalCentralizedConfigurationName() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MetadataMock_GetGlobalCentralizedConfigurationName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetGlobalCentralizedConfigurationName'
type MetadataMock_GetGlobalCentralizedConfigurationName_Call struct {
	*mock.Call
}

// GetGlobalCentralizedConfigurationName is a helper method to define mock.On call
func (_e *MetadataMock_Expecter) GetGlobalCentralizedConfigurationName() *MetadataMock_GetGlobalCentralizedConfigurationName_Call {
	return &MetadataMock_GetGlobalCentralizedConfigurationName_Call{Call: _e.mock.On("GetGlobalCentralizedConfigurationName")}
}

func (_c *MetadataMock_GetGlobalCentralizedConfigurationName_Call) Run(run func()) *MetadataMock_GetGlobalCentralizedConfigurationName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MetadataMock_GetGlobalCentralizedConfigurationName_Call) Return(_a0 string) *MetadataMock_GetGlobalCentralizedConfigurationName_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MetadataMock_GetGlobalCentralizedConfigurationName_Call) RunAndReturn(run func() string) *MetadataMock_GetGlobalCentralizedConfigurationName_Call {
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

// GetProcessCentralizedConfigurationName provides a mock function with given fields:
func (_m *MetadataMock) GetProcessCentralizedConfigurationName() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MetadataMock_GetProcessCentralizedConfigurationName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetProcessCentralizedConfigurationName'
type MetadataMock_GetProcessCentralizedConfigurationName_Call struct {
	*mock.Call
}

// GetProcessCentralizedConfigurationName is a helper method to define mock.On call
func (_e *MetadataMock_Expecter) GetProcessCentralizedConfigurationName() *MetadataMock_GetProcessCentralizedConfigurationName_Call {
	return &MetadataMock_GetProcessCentralizedConfigurationName_Call{Call: _e.mock.On("GetProcessCentralizedConfigurationName")}
}

func (_c *MetadataMock_GetProcessCentralizedConfigurationName_Call) Run(run func()) *MetadataMock_GetProcessCentralizedConfigurationName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MetadataMock_GetProcessCentralizedConfigurationName_Call) Return(_a0 string) *MetadataMock_GetProcessCentralizedConfigurationName_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MetadataMock_GetProcessCentralizedConfigurationName_Call) RunAndReturn(run func() string) *MetadataMock_GetProcessCentralizedConfigurationName_Call {
	_c.Call.Return(run)
	return _c
}

// GetProcessType provides a mock function with given fields:
func (_m *MetadataMock) GetProcessType() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MetadataMock_GetProcessType_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetProcessType'
type MetadataMock_GetProcessType_Call struct {
	*mock.Call
}

// GetProcessType is a helper method to define mock.On call
func (_e *MetadataMock_Expecter) GetProcessType() *MetadataMock_GetProcessType_Call {
	return &MetadataMock_GetProcessType_Call{Call: _e.mock.On("GetProcessType")}
}

func (_c *MetadataMock_GetProcessType_Call) Run(run func()) *MetadataMock_GetProcessType_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MetadataMock_GetProcessType_Call) Return(_a0 string) *MetadataMock_GetProcessType_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MetadataMock_GetProcessType_Call) RunAndReturn(run func() string) *MetadataMock_GetProcessType_Call {
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

// GetProductCentralizedConfigurationName provides a mock function with given fields:
func (_m *MetadataMock) GetProductCentralizedConfigurationName() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MetadataMock_GetProductCentralizedConfigurationName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetProductCentralizedConfigurationName'
type MetadataMock_GetProductCentralizedConfigurationName_Call struct {
	*mock.Call
}

// GetProductCentralizedConfigurationName is a helper method to define mock.On call
func (_e *MetadataMock_Expecter) GetProductCentralizedConfigurationName() *MetadataMock_GetProductCentralizedConfigurationName_Call {
	return &MetadataMock_GetProductCentralizedConfigurationName_Call{Call: _e.mock.On("GetProductCentralizedConfigurationName")}
}

func (_c *MetadataMock_GetProductCentralizedConfigurationName_Call) Run(run func()) *MetadataMock_GetProductCentralizedConfigurationName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MetadataMock_GetProductCentralizedConfigurationName_Call) Return(_a0 string) *MetadataMock_GetProductCentralizedConfigurationName_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MetadataMock_GetProductCentralizedConfigurationName_Call) RunAndReturn(run func() string) *MetadataMock_GetProductCentralizedConfigurationName_Call {
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

// GetWorkflowCentralizedConfigurationName provides a mock function with given fields:
func (_m *MetadataMock) GetWorkflowCentralizedConfigurationName() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MetadataMock_GetWorkflowCentralizedConfigurationName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetWorkflowCentralizedConfigurationName'
type MetadataMock_GetWorkflowCentralizedConfigurationName_Call struct {
	*mock.Call
}

// GetWorkflowCentralizedConfigurationName is a helper method to define mock.On call
func (_e *MetadataMock_Expecter) GetWorkflowCentralizedConfigurationName() *MetadataMock_GetWorkflowCentralizedConfigurationName_Call {
	return &MetadataMock_GetWorkflowCentralizedConfigurationName_Call{Call: _e.mock.On("GetWorkflowCentralizedConfigurationName")}
}

func (_c *MetadataMock_GetWorkflowCentralizedConfigurationName_Call) Run(run func()) *MetadataMock_GetWorkflowCentralizedConfigurationName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MetadataMock_GetWorkflowCentralizedConfigurationName_Call) Return(_a0 string) *MetadataMock_GetWorkflowCentralizedConfigurationName_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MetadataMock_GetWorkflowCentralizedConfigurationName_Call) RunAndReturn(run func() string) *MetadataMock_GetWorkflowCentralizedConfigurationName_Call {
	_c.Call.Return(run)
	return _c
}

// GetWorkflowType provides a mock function with given fields:
func (_m *MetadataMock) GetWorkflowType() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MetadataMock_GetWorkflowType_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetWorkflowType'
type MetadataMock_GetWorkflowType_Call struct {
	*mock.Call
}

// GetWorkflowType is a helper method to define mock.On call
func (_e *MetadataMock_Expecter) GetWorkflowType() *MetadataMock_GetWorkflowType_Call {
	return &MetadataMock_GetWorkflowType_Call{Call: _e.mock.On("GetWorkflowType")}
}

func (_c *MetadataMock_GetWorkflowType_Call) Run(run func()) *MetadataMock_GetWorkflowType_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MetadataMock_GetWorkflowType_Call) Return(_a0 string) *MetadataMock_GetWorkflowType_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MetadataMock_GetWorkflowType_Call) RunAndReturn(run func() string) *MetadataMock_GetWorkflowType_Call {
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
