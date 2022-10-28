// Code generated by MockGen. DO NOT EDIT.
// Source: internal/client/redshift/client.go

// Package mock_redshift is a generated GoMock package.
package mock_redshift

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

// DailyBuyerDemandTableRefresh mocks base method.
func (m *MockClient) DailyBuyerDemandTableRefresh(ctx context.Context, bds model.BuyerDemands) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DailyBuyerDemandTableRefresh", ctx, bds)
	ret0, _ := ret[0].(error)
	return ret0
}

// DailyBuyerDemandTableRefresh indicates an expected call of DailyBuyerDemandTableRefresh.
func (mr *MockClientMockRecorder) DailyBuyerDemandTableRefresh(ctx, bds interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DailyBuyerDemandTableRefresh", reflect.TypeOf((*MockClient)(nil).DailyBuyerDemandTableRefresh), ctx, bds)
}
