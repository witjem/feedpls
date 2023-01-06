// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package feed_test

import (
	"context"
	"github.com/witjem/feedpls/internal/pkg/feed"
	"io"
	"sync"
)

// Ensure, that HTTPClientMock does implement feed.HTTPClient.
// If this is not the case, regenerate this file with moq.
var _ feed.HTTPClient = &HTTPClientMock{}

// HTTPClientMock is a mock implementation of feed.HTTPClient.
//
//	func TestSomethingThatUsesHTTPClient(t *testing.T) {
//
//		// make and configure a mocked feed.HTTPClient
//		mockedHTTPClient := &HTTPClientMock{
//			GetFunc: func(ctx context.Context, url string) (io.ReadCloser, error) {
//				panic("mock out the Get method")
//			},
//		}
//
//		// use mockedHTTPClient in code that requires feed.HTTPClient
//		// and then make assertions.
//
//	}
type HTTPClientMock struct {
	// GetFunc mocks the Get method.
	GetFunc func(ctx context.Context, url string) (io.ReadCloser, error)

	// calls tracks calls to the methods.
	calls struct {
		// Get holds details about calls to the Get method.
		Get []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// URL is the url argument value.
			URL string
		}
	}
	lockGet sync.RWMutex
}

// Get calls GetFunc.
func (mock *HTTPClientMock) Get(ctx context.Context, url string) (io.ReadCloser, error) {
	if mock.GetFunc == nil {
		panic("HTTPClientMock.GetFunc: method is nil but HTTPClient.Get was just called")
	}
	callInfo := struct {
		Ctx context.Context
		URL string
	}{
		Ctx: ctx,
		URL: url,
	}
	mock.lockGet.Lock()
	mock.calls.Get = append(mock.calls.Get, callInfo)
	mock.lockGet.Unlock()
	return mock.GetFunc(ctx, url)
}

// GetCalls gets all the calls that were made to Get.
// Check the length with:
//
//	len(mockedHTTPClient.GetCalls())
func (mock *HTTPClientMock) GetCalls() []struct {
	Ctx context.Context
	URL string
} {
	var calls []struct {
		Ctx context.Context
		URL string
	}
	mock.lockGet.RLock()
	calls = mock.calls.Get
	mock.lockGet.RUnlock()
	return calls
}