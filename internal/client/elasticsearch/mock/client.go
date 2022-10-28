// Code generated by MockGen. DO NOT EDIT.
// Source: internal/client/elasticsearch/client.go

// Package mock_elasticsearch is a generated GoMock package.
package mock_elasticsearch

import (
	context "context"
	reflect "reflect"

	model "github.com/HomesNZ/buyer-demand/internal/model"
	gomock "github.com/golang/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// QueryAll mocks base method.
func (m *MockClient) QueryAll(ctx context.Context) (model.MapItemESs, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryAll", ctx)
	ret0, _ := ret[0].(model.MapItemESs)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryAll indicates an expected call of QueryAll.
func (mr *MockClientMockRecorder) QueryAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryAll", reflect.TypeOf((*MockClient)(nil).QueryAll), ctx)
}
