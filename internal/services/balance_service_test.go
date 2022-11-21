package services

import (
	"context"
	"testing"
	"time"

	"github.com/DmitriyV003/bonus/internal/mocks"
	"github.com/DmitriyV003/bonus/internal/models"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestBalanceService_Accrual(t *testing.T) {
	type args struct {
		payment *models.Payment
		user    *models.User
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		wantBalance int64
	}{
		{
			name:        "Success accrual balance",
			wantErr:     false,
			wantBalance: 30000,
			args: args{
				payment: &models.Payment{
					ID:              1,
					Type:            models.AccrualType,
					TransactionType: models.DEBIT,
					OrderNumber:     "5469980471561486",
					Amount:          10000,
					User:            nil,
					CreatedAt:       time.Time{},
					UpdatedAt:       nil,
				},
				user: &models.User{
					ID:        1,
					Login:     "test",
					Password:  "secret",
					Balance:   20000,
					CreatedAt: time.Time{},
					UpdatedAt: nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPayments := mocks.NewMockPaymentRepository(ctrl)
			mockUsers := mocks.NewMockUserRepository(ctrl)
			bs := &BalanceService{
				payments: mockPayments,
				users:    mockUsers,
			}

			mockUsers.EXPECT().UpdateBalance(context.Background(), tt.args.user).Times(1).Return(nil)

			if err := bs.Accrual(context.Background(), tt.args.payment, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("Accrual() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.wantBalance, tt.args.user.Balance)
		})
	}
}

func TestBalanceService_Withdraw(t *testing.T) {
	type args struct {
		payment *models.Payment
		user    *models.User
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		wantBalance int64
	}{
		{
			name:        "Success withdraw from balance",
			wantErr:     false,
			wantBalance: 10000,
			args: args{
				payment: &models.Payment{
					ID:              1,
					Type:            models.WithdrawType,
					TransactionType: models.CREDIT,
					OrderNumber:     "5469980471561486",
					Amount:          10000,
					User:            nil,
					CreatedAt:       time.Time{},
					UpdatedAt:       nil,
				},
				user: &models.User{
					ID:        1,
					Login:     "test",
					Password:  "secret",
					Balance:   20000,
					CreatedAt: time.Time{},
					UpdatedAt: nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPayments := mocks.NewMockPaymentRepository(ctrl)
			mockUsers := mocks.NewMockUserRepository(ctrl)
			bs := &BalanceService{
				payments: mockPayments,
				users:    mockUsers,
			}

			mockUsers.EXPECT().UpdateBalance(context.Background(), tt.args.user).Times(1).Return(nil)

			if err := bs.Withdraw(context.Background(), tt.args.payment, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("Withdraw() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.wantBalance, tt.args.user.Balance)
		})
	}
}
