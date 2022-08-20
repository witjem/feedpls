// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"
	io "io"

	mock "github.com/stretchr/testify/mock"
)

// WebClient is an autogenerated mock type for the WebClient type
type WebClient struct {
	mock.Mock
}

// Get provides a mock function with given fields: ctx, url
func (_m *WebClient) Get(ctx context.Context, url string) (io.ReadCloser, error) {
	ret := _m.Called(ctx, url)

	var r0 io.ReadCloser
	if rf, ok := ret.Get(0).(func(context.Context, string) io.ReadCloser); ok {
		r0 = rf(ctx, url)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(io.ReadCloser)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, url)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewWebClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewWebClient creates a new instance of WebClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewWebClient(t mockConstructorTestingTNewWebClient) *WebClient {
	mock := &WebClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}