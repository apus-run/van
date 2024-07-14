// Package cache is a generated GoMock package.
package cache

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockKeyGetter is a mock of KeyGetter interface.
type MockKeyGetter struct {
	ctrl     *gomock.Controller
	recorder *MockKeyGetterMockRecorder
}

// MockKeyGetterMockRecorder is the mock recorder for MockKeyGetter.
type MockKeyGetterMockRecorder struct {
	mock *MockKeyGetter
}

// NewMockKeyGetter creates a new mock instance.
func NewMockKeyGetter(ctrl *gomock.Controller) *MockKeyGetter {
	mock := &MockKeyGetter{ctrl: ctrl}
	mock.recorder = &MockKeyGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKeyGetter) EXPECT() *MockKeyGetterMockRecorder {
	return m.recorder
}

// CacheKey mocks base method.
func (m *MockKeyGetter) CacheKey() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CacheKey")
	ret0, _ := ret[0].(string)
	return ret0
}

// CacheKey indicates an expected call of CacheKey.
func (mr *MockKeyGetterMockRecorder) CacheKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CacheKey", reflect.TypeOf((*MockKeyGetter)(nil).CacheKey))
}
