// Code generated by mockery v1.0.0
package mocks

import http "net/http"
import mock "github.com/stretchr/testify/mock"

// HandlerInterface is an autogenerated mock type for the HandlerInterface type
type HandlerInterface struct {
	mock.Mock
}

// Get provides a mock function with given fields: _a0
func (_m *HandlerInterface) Get(_a0 string) (*http.Response, error) {
	ret := _m.Called(_a0)

	var r0 *http.Response
	if rf, ok := ret.Get(0).(func(string) *http.Response); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Response)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetById provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *HandlerInterface) GetById(_a0 string, _a1 string, _a2 string, _a3 interface{}) error {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string, interface{}) error); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Getenv provides a mock function with given fields: _a0
func (_m *HandlerInterface) Getenv(_a0 string) string {
	ret := _m.Called(_a0)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}