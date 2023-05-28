// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// NotificationHandler is an autogenerated mock type for the NotificationHandler type
type NotificationHandler struct {
	mock.Mock
}

// OnMessage provides a mock function with given fields: _a0
func (_m *NotificationHandler) OnMessage(_a0 string) {
	_m.Called(_a0)
}

type mockConstructorTestingTNewNotificationHandler interface {
	mock.TestingT
	Cleanup(func())
}

// NewNotificationHandler creates a new instance of NotificationHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewNotificationHandler(t mockConstructorTestingTNewNotificationHandler) *NotificationHandler {
	mock := &NotificationHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
