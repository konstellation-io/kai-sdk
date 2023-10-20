// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import (
	messaging "github.com/konstellation-io/kai-sdk/go-sdk/sdk/messaging"
	mock "github.com/stretchr/testify/mock"
)

// CentralizedConfigMock is an autogenerated mock type for the centralizedConfig type
type CentralizedConfigMock struct {
	mock.Mock
}

type CentralizedConfigMock_Expecter struct {
	mock *mock.Mock
}

func (_m *CentralizedConfigMock) EXPECT() *CentralizedConfigMock_Expecter {
	return &CentralizedConfigMock_Expecter{mock: &_m.Mock}
}

// DeleteConfig provides a mock function with given fields: key, scope
func (_m *CentralizedConfigMock) DeleteConfig(key string, scope messaging.Scope) error {
	ret := _m.Called(key, scope)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, messaging.Scope) error); ok {
		r0 = rf(key, scope)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CentralizedConfigMock_DeleteConfig_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteConfig'
type CentralizedConfigMock_DeleteConfig_Call struct {
	*mock.Call
}

// DeleteConfig is a helper method to define mock.On call
//   - key string
//   - scope messaging.Scope
func (_e *CentralizedConfigMock_Expecter) DeleteConfig(key interface{}, scope interface{}) *CentralizedConfigMock_DeleteConfig_Call {
	return &CentralizedConfigMock_DeleteConfig_Call{Call: _e.mock.On("DeleteConfig", key, scope)}
}

func (_c *CentralizedConfigMock_DeleteConfig_Call) Run(run func(key string, scope messaging.Scope)) *CentralizedConfigMock_DeleteConfig_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(messaging.Scope))
	})
	return _c
}

func (_c *CentralizedConfigMock_DeleteConfig_Call) Return(_a0 error) *CentralizedConfigMock_DeleteConfig_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CentralizedConfigMock_DeleteConfig_Call) RunAndReturn(run func(string, messaging.Scope) error) *CentralizedConfigMock_DeleteConfig_Call {
	_c.Call.Return(run)
	return _c
}

// GetConfig provides a mock function with given fields: key, scope
func (_m *CentralizedConfigMock) GetConfig(key string, scope ...messaging.Scope) (string, error) {
	_va := make([]interface{}, len(scope))
	for _i := range scope {
		_va[_i] = scope[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, key)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string, ...messaging.Scope) (string, error)); ok {
		return rf(key, scope...)
	}
	if rf, ok := ret.Get(0).(func(string, ...messaging.Scope) string); ok {
		r0 = rf(key, scope...)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string, ...messaging.Scope) error); ok {
		r1 = rf(key, scope...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CentralizedConfigMock_GetConfig_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetConfig'
type CentralizedConfigMock_GetConfig_Call struct {
	*mock.Call
}

// GetConfig is a helper method to define mock.On call
//   - key string
//   - scope ...messaging.Scope
func (_e *CentralizedConfigMock_Expecter) GetConfig(key interface{}, scope ...interface{}) *CentralizedConfigMock_GetConfig_Call {
	return &CentralizedConfigMock_GetConfig_Call{Call: _e.mock.On("GetConfig",
		append([]interface{}{key}, scope...)...)}
}

func (_c *CentralizedConfigMock_GetConfig_Call) Run(run func(key string, scope ...messaging.Scope)) *CentralizedConfigMock_GetConfig_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]messaging.Scope, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(messaging.Scope)
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *CentralizedConfigMock_GetConfig_Call) Return(_a0 string, _a1 error) *CentralizedConfigMock_GetConfig_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CentralizedConfigMock_GetConfig_Call) RunAndReturn(run func(string, ...messaging.Scope) (string, error)) *CentralizedConfigMock_GetConfig_Call {
	_c.Call.Return(run)
	return _c
}

// SetConfig provides a mock function with given fields: key, value, scope
func (_m *CentralizedConfigMock) SetConfig(key string, value string, scope ...messaging.Scope) error {
	_va := make([]interface{}, len(scope))
	for _i := range scope {
		_va[_i] = scope[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, key, value)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, ...messaging.Scope) error); ok {
		r0 = rf(key, value, scope...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CentralizedConfigMock_SetConfig_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetConfig'
type CentralizedConfigMock_SetConfig_Call struct {
	*mock.Call
}

// SetConfig is a helper method to define mock.On call
//   - key string
//   - value string
//   - scope ...messaging.Scope
func (_e *CentralizedConfigMock_Expecter) SetConfig(key interface{}, value interface{}, scope ...interface{}) *CentralizedConfigMock_SetConfig_Call {
	return &CentralizedConfigMock_SetConfig_Call{Call: _e.mock.On("SetConfig",
		append([]interface{}{key, value}, scope...)...)}
}

func (_c *CentralizedConfigMock_SetConfig_Call) Run(run func(key string, value string, scope ...messaging.Scope)) *CentralizedConfigMock_SetConfig_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]messaging.Scope, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(messaging.Scope)
			}
		}
		run(args[0].(string), args[1].(string), variadicArgs...)
	})
	return _c
}

func (_c *CentralizedConfigMock_SetConfig_Call) Return(_a0 error) *CentralizedConfigMock_SetConfig_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CentralizedConfigMock_SetConfig_Call) RunAndReturn(run func(string, string, ...messaging.Scope) error) *CentralizedConfigMock_SetConfig_Call {
	_c.Call.Return(run)
	return _c
}

// NewCentralizedConfigMock creates a new instance of CentralizedConfigMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCentralizedConfigMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *CentralizedConfigMock {
	mock := &CentralizedConfigMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
