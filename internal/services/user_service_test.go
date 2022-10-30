package services

import (
	"context"
	"github.com/DmitriyV003/bonus/internal/applicationerrors"
	"github.com/DmitriyV003/bonus/internal/mocks"
	"github.com/DmitriyV003/bonus/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestUserService_Withdraw(t *testing.T) {
	type args struct {
		user        *models.User
		orderNumber int64
		sum         float64
	}
	tests := []struct {
		name               string
		args               args
		wantErr            bool
		isOrderNumberValid bool
		isLowFunds         bool
	}{
		{
			name:               "Withdraw for order with invalid number",
			wantErr:            true,
			isOrderNumberValid: false,
			isLowFunds:         false,
			args: args{
				user: &models.User{
					ID:        1,
					Login:     "test",
					Password:  "secret",
					Balance:   20000,
					CreatedAt: time.Time{},
					UpdatedAt: nil,
				},
				orderNumber: 546998047156148,
				sum:         100.00,
			},
		},
		{
			name:               "User withdraw with low funds",
			wantErr:            true,
			isOrderNumberValid: true,
			isLowFunds:         true,
			args: args{
				user: &models.User{
					ID:        1,
					Login:     "test",
					Password:  "secret",
					Balance:   10,
					CreatedAt: time.Time{},
					UpdatedAt: nil,
				},
				orderNumber: 546998047156148,
				sum:         100.00,
			},
		},
		{
			name:               "User withdraw successfully",
			wantErr:            false,
			isOrderNumberValid: true,
			isLowFunds:         false,
			args: args{
				user: &models.User{
					ID:        1,
					Login:     "test",
					Password:  "secret",
					Balance:   10000000,
					CreatedAt: time.Time{},
					UpdatedAt: nil,
				},
				orderNumber: 5469980471561486,
				sum:         100.00,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockValidator := mocks.NewMockOrderValidator(ctrl)
			mockUsers := mocks.NewMockUserRepository(ctrl)
			mockPayments := mocks.NewMockPaymentRepository(ctrl)
			mockPaymentService := mocks.NewMockPaymentService(ctrl)
			mockBalanceService := mocks.NewMockBalanceService(ctrl)
			u := &UserService{
				validator:      mockValidator,
				paymentService: mockPaymentService,
				users:          mockUsers,
				payments:       mockPayments,
				balanceService: mockBalanceService,
			}
			if !tt.isOrderNumberValid {
				mockValidator.EXPECT().Validate(tt.args.orderNumber).Times(1).Return(false)
			} else if tt.isLowFunds {
				mockValidator.EXPECT().Validate(tt.args.orderNumber).Times(1).Return(true)
			} else {
				mockValidator.EXPECT().Validate(tt.args.orderNumber).Times(1).Return(true)
				payment := models.Payment{
					ID:              1,
					Type:            models.WithdrawType,
					TransactionType: models.CREDIT,
					OrderNumber:     strconv.FormatInt(tt.args.orderNumber, 10),
					Amount:          int64(tt.args.sum * 10000),
					User:            nil,
					CreatedAt:       time.Now(),
					UpdatedAt:       nil,
				}
				mockPaymentService.EXPECT().
					CreateWithdrawPayment(context.Background(), tt.args.user, int64(tt.args.sum*10000), strconv.FormatInt(tt.args.orderNumber, 10)).
					Times(1).
					Return(&payment, nil)
				mockBalanceService.EXPECT().Withdraw(context.Background(), &payment, tt.args.user).
					Times(1).
					Return(nil)
			}

			if err := u.Withdraw(context.Background(), tt.args.user, strconv.FormatInt(tt.args.orderNumber, 10), tt.args.sum); (err != nil) != tt.wantErr {
				t.Errorf("Withdraw() error = %v, wantErr: %v", err, tt.wantErr)
				if !tt.isOrderNumberValid {
					assert.ErrorAs(t, err, applicationerrors.ErrInvalidOrderNumber)
				} else if tt.isLowFunds {
					assert.ErrorAs(t, err, applicationerrors.ErrLowUserABalance)
				}
			}
		})
	}
}
