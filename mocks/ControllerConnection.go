// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ControllerConnection is an autogenerated mock type for the ControllerConnection type
type ControllerConnection struct {
	mock.Mock
}

// Read provides a mock function with given fields: _a0
func (_m *ControllerConnection) Read(_a0 []byte) (int, error) {
	ret := _m.Called(_a0)

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte) (int, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func([]byte) int); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Write provides a mock function with given fields: _a0
func (_m *ControllerConnection) Write(_a0 []byte) (int, error) {
	ret := _m.Called(_a0)

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte) (int, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func([]byte) int); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewControllerConnection interface {
	mock.TestingT
	Cleanup(func())
}

// NewControllerConnection creates a new instance of ControllerConnection. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewControllerConnection(t mockConstructorTestingTNewControllerConnection) *ControllerConnection {
	mock := &ControllerConnection{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}