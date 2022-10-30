// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/clients/client_interfaces/bonus_client.go

// Package mock_client_interfaces is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	client_interfaces "github.com/DmitriyV003/bonus/internal/clients/clientinterfaces"
	gomock "github.com/golang/mock/gomock"
)

// MockBonusClient is a mock of BonusClient interface.
type MockBonusClient struct {
	ctrl     *gomock.Controller
	recorder *MockBonusClientMockRecorder
}

// MockBonusClientMockRecorder is the mock recorder for MockBonusClient.
type MockBonusClientMockRecorder struct {
	mock *MockBonusClient
}

// NewMockBonusClient creates a new mock instance.
func NewMockBonusClient(ctrl *gomock.Controller) *MockBonusClient {
	mock := &MockBonusClient{ctrl: ctrl}
	mock.recorder = &MockBonusClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBonusClient) EXPECT() *MockBonusClientMockRecorder {
	return m.recorder
}

// CreateOrder mocks base method.
func (m *MockBonusClient) CreateOrder(orderNumber string) (*client_interfaces.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrder", orderNumber)
	ret0, _ := ret[0].(*client_interfaces.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateOrder indicates an expected call of CreateOrder.
func (mr *MockBonusClientMockRecorder) CreateOrder(orderNumber interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrder", reflect.TypeOf((*MockBonusClient)(nil).CreateOrder), orderNumber)
}

// GetOrderDetails mocks base method.
func (m *MockBonusClient) GetOrderDetails(orderNumber string) (*client_interfaces.OrderDetailsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrderDetails", orderNumber)
	ret0, _ := ret[0].(*client_interfaces.OrderDetailsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrderDetails indicates an expected call of GetOrderDetails.
func (mr *MockBonusClientMockRecorder) GetOrderDetails(orderNumber interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrderDetails", reflect.TypeOf((*MockBonusClient)(nil).GetOrderDetails), orderNumber)
}
