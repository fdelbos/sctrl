// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	sctrl "github.com/fdelbos/sctrl"
	mock "github.com/stretchr/testify/mock"
)

// ResponseHandler is an autogenerated mock type for the ResponseHandler type
type ResponseHandler struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0
func (_m *ResponseHandler) Execute(_a0 sctrl.ResponseReader) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(sctrl.ResponseReader) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewResponseHandler interface {
	mock.TestingT
	Cleanup(func())
}

// NewResponseHandler creates a new instance of ResponseHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewResponseHandler(t mockConstructorTestingTNewResponseHandler) *ResponseHandler {
	mock := &ResponseHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
