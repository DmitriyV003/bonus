package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/DmitriyV003/bonus/internal/applicationerrors"
	"github.com/DmitriyV003/bonus/internal/clients/clientinterfaces"
	"github.com/DmitriyV003/bonus/internal/mocks"
	"github.com/DmitriyV003/bonus/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestOrderService_Create(t *testing.T) {
	type fields struct {
		validator      *mocks.MockOrderValidator
		bonusClient    *mocks.MockBonusClient
		orders         *mocks.MockOrderRepository
		users          *mocks.MockUserRepository
		paymentService *mocks.MockPaymentService
	}
	type args struct {
		user        *models.User
		orderNumber string
	}
	user := models.User{
		ID:        1,
		Login:     "test_user",
		Password:  "secret",
		Balance:   10000,
		CreatedAt: time.Now(),
		UpdatedAt: nil,
	}
	orderCreationTime := time.Now()
	models.NewOrder = func(number string, status string, amount int64, user *models.User) *models.Order {
		return &models.Order{
			ID:        0,
			Number:    "5469980471561486",
			Status:    models.NewStatus,
			Amount:    0,
			User:      user,
			CreatedAt: orderCreationTime,
			UpdatedAt: nil,
		}
	}

	tests := []struct {
		name      string
		args      args
		want      *models.Order
		prepare   func(f *fields)
		wantErr   assert.ErrorAssertionFunc
		targetErr error
	}{
		{
			name: "Create order with invalid number",
			args: args{
				user:        &user,
				orderNumber: "546998047156148",
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, applicationerrors.ErrInvalidOrderNumber)
			},
			targetErr: applicationerrors.ErrInvalidOrderNumber,
			prepare: func(f *fields) {
				f.validator.EXPECT().Validate(int64(546998047156148)).
					Times(1).
					Return(false)
			},
		},
		{
			name: "Create order with another user",
			args: args{
				user:        &user,
				orderNumber: "5469980471561486",
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, applicationerrors.ErrConflict)
			},
			targetErr: applicationerrors.ErrConflict,
			prepare: func(f *fields) {
				gomock.InOrder(
					f.validator.EXPECT().
						Validate(int64(5469980471561486)).
						Times(1).
						Return(true),
					f.orders.EXPECT().
						GetByNumber(context.Background(), "5469980471561486").
						Times(1).
						Return(&models.Order{
							ID:     1,
							Number: "5469980471561486",
							Status: models.ProcessingStatus,
							Amount: 20000,
							User: &models.User{
								ID:        2,
								Login:     "test_user_2",
								Password:  "secret",
								Balance:   10000,
								CreatedAt: time.Now(),
								UpdatedAt: nil,
							},
							CreatedAt: time.Now(),
							UpdatedAt: nil,
						}, nil),
				)

			},
		},
		{
			name: "Create already created order",
			args: args{
				user:        &user,
				orderNumber: "5469980471561486",
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, applicationerrors.ErrModelAlreadyCreated)
			},
			targetErr: applicationerrors.ErrModelAlreadyCreated,
			prepare: func(f *fields) {
				f.validator.EXPECT().
					Validate(int64(5469980471561486)).
					Times(1).
					Return(true)
				f.orders.EXPECT().
					GetByNumber(context.Background(), "5469980471561486").
					Times(1).
					Return(&models.Order{
						ID:        1,
						Number:    "5469980471561486",
						Status:    models.ProcessingStatus,
						Amount:    20000,
						User:      &user,
						CreatedAt: orderCreationTime,
						UpdatedAt: nil,
					}, nil)
			},
		},
		{
			name: "Create new order successfully",
			args: args{
				user:        &user,
				orderNumber: "5469980471561486",
			},
			want: &models.Order{
				ID:        1,
				Number:    "5469980471561486",
				Status:    models.ProcessedStatus,
				Amount:    20000,
				User:      &user,
				CreatedAt: orderCreationTime,
				UpdatedAt: nil,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return false
			},
			targetErr: nil,
			prepare: func(f *fields) {
				order := models.Order{
					ID:        1,
					Number:    "5469980471561486",
					Status:    models.NewStatus,
					Amount:    0,
					User:      &user,
					CreatedAt: orderCreationTime,
					UpdatedAt: nil,
				}
				f.validator.EXPECT().
					Validate(int64(5469980471561486)).
					Times(1).
					Return(true)
				f.orders.EXPECT().
					GetByNumber(context.Background(), "5469980471561486").
					Times(1).
					Return(nil, applicationerrors.ErrNotFound)
				f.orders.EXPECT().
					Create(context.Background(), &models.Order{
						ID:        0,
						Number:    "5469980471561486",
						Status:    models.NewStatus,
						Amount:    0,
						User:      &user,
						CreatedAt: orderCreationTime,
						UpdatedAt: nil,
					}).
					Times(1).
					Return(&order, nil)
				f.bonusClient.EXPECT().
					GetOrderDetails(order.Number).
					Times(1).
					Return(&clientinterfaces.OrderDetailsResponse{
						Order:  order.Number,
						Status: models.ProcessedStatus,
						Amount: 2,
					}, nil)
				f.orders.EXPECT().
					UpdateByID(context.Background(), &models.Order{
						ID:        1,
						Number:    "5469980471561486",
						Status:    models.ProcessedStatus,
						Amount:    20000,
						User:      &user,
						CreatedAt: orderCreationTime,
						UpdatedAt: nil,
					}).
					Times(1).
					Return(nil)
				f.paymentService.EXPECT().
					CreateAccrualPayment(context.Background(), &user, int64(2*10000), order.Number).
					Times(1).
					Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockValidator := mocks.NewMockOrderValidator(ctrl)
			mockUsers := mocks.NewMockUserRepository(ctrl)
			mockOrders := mocks.NewMockOrderRepository(ctrl)
			mockPaymentService := mocks.NewMockPaymentService(ctrl)
			mockBonusClient := mocks.NewMockBonusClient(ctrl)
			mockFields := fields{
				validator:      mockValidator,
				bonusClient:    mockBonusClient,
				orders:         mockOrders,
				users:          mockUsers,
				paymentService: mockPaymentService,
			}
			myself := &OrderService{
				validator:      mockValidator,
				bonusClient:    mockBonusClient,
				orders:         mockOrders,
				users:          mockUsers,
				paymentService: mockPaymentService,
			}
			if tt.prepare != nil {
				tt.prepare(&mockFields)
			}
			got, err := myself.Create(context.Background(), tt.args.user, tt.args.orderNumber)
			if tt.wantErr(t, err, fmt.Sprintf("Create(%v, %v, %v)", context.Background(), tt.args.user, tt.args.orderNumber)) {
				assert.ErrorIs(t, err, tt.targetErr)
			}
			assert.Equalf(t, tt.want, got, "Create(%v, %v, %v)", context.Background(), tt.args.user, tt.args.orderNumber)
		})
	}
}
