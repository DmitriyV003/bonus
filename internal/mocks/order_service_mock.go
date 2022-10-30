// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/services/interfaces/order_service.go

// Package mock_interfaces is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	models "github.com/DmitriyV003/bonus/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockOrderService is a mock of OrderService interface.
type MockOrderService struct {
	ctrl     *gomock.Controller
	recorder *MockOrderServiceMockRecorder
}

// MockOrderServiceMockRecorder is the mock recorder for MockOrderService.
type MockOrderServiceMockRecorder struct {
	mock *MockOrderService
}

// NewMockOrderService creates a new mock instance.
func NewMockOrderService(ctrl *gomock.Controller) *MockOrderService {
	mock := &MockOrderService{ctrl: ctrl}
	mock.recorder = &MockOrderServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderService) EXPECT() *MockOrderServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockOrderService) Create(ctx context.Context, user *models.User, orderNumber string) (*models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, user, orderNumber)
	ret0, _ := ret[0].(*models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockOrderServiceMockRecorder) Create(ctx, user, orderNumber interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockOrderService)(nil).Create), ctx, user, orderNumber)
}

// OrdersByUser mocks base method.
func (m *MockOrderService) OrdersByUser(ctx context.Context, user *models.User) ([]*models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OrdersByUser", ctx, user)
	ret0, _ := ret[0].([]*models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OrdersByUser indicates an expected call of OrdersByUser.
func (mr *MockOrderServiceMockRecorder) OrdersByUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OrdersByUser", reflect.TypeOf((*MockOrderService)(nil).OrdersByUser), ctx, user)
}

// PollPendingOrders mocks base method.
func (m *MockOrderService) PollPendingOrders(ctx context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PollPendingOrders", ctx)
}

// PollPendingOrders indicates an expected call of PollPendingOrders.
func (mr *MockOrderServiceMockRecorder) PollPendingOrders(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PollPendingOrders", reflect.TypeOf((*MockOrderService)(nil).PollPendingOrders), ctx)
}
