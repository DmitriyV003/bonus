package services

import (
	"context"
	"fmt"
	"github.com/DmitriyV003/bonus/internal/applicationerrors"
	"github.com/DmitriyV003/bonus/internal/models"
	"github.com/DmitriyV003/bonus/internal/repository"
	"github.com/DmitriyV003/bonus/internal/requests"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

type UserService struct {
	validator      OrderValidator
	paymentService *PaymentService
	users          *repository.UserRepository
	payments       *repository.PaymentRepository
	authService    *AuthService
}

func NewUserService(
	validator OrderValidator,
	paymentService *PaymentService,
	users *repository.UserRepository,
	payments *repository.PaymentRepository,
	authService *AuthService,
) *UserService {
	return &UserService{
		validator:      validator,
		paymentService: paymentService,
		users:          users,
		payments:       payments,
		authService:    authService,
	}
}

func (u *UserService) Create(ctx context.Context, request *requests.RegistrationRequest) (*Token, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(request.Password), 14)
	if err != nil {
		return nil, fmt.Errorf("error to generate hash from password: %w", err)
	}

	user := models.User{
		Login:     request.Login,
		Password:  string(bytes),
		CreatedAt: time.Now(),
	}
	err = u.users.Create(ctx, &user)
	if err != nil {
		return nil, fmt.Errorf("unable to create user in db: %w", err)
	}

	dbUser, err := u.users.GetByLogin(ctx, request.Login)
	if err != nil {
		return nil, fmt.Errorf("error to get user by login: %w", err)
	}

	token, err := u.authService.LoginByUser(dbUser)
	if err != nil {
		return nil, fmt.Errorf("login user programmly: %w", err)
	}

	return token, nil
}

func (u *UserService) Withdraw(ctx context.Context, user *models.User, orderNumber string, sum float64) error {
	parsedOderNumber, err := strconv.ParseInt(orderNumber, 10, 64)
	if err != nil {
		return err
	}

	isValid := u.validator.Validate(parsedOderNumber)
	if !isValid {
		return fmt.Errorf("unable to validate order number %w", applicationerrors.ErrInvalidOrderNumber)
	}

	sumToWithdraw := int64(sum * 10000)
	if user.Balance < sumToWithdraw {
		return fmt.Errorf("unable to validate order number %w", applicationerrors.ErrLowUserABalance)
	}

	err = u.paymentService.CreateWithdrawPayment(ctx, user, sumToWithdraw, orderNumber)
	if err != nil {
		return fmt.Errorf("unable ro create payment for withdraw %w", err)
	}

	return nil
}

func (u *UserService) AllWithdrawsByUser(ctx context.Context, user *models.User) ([]*models.Payment, error) {
	payments, err := u.payments.GetWithdrawsByUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("error to get all withdraws: %w", err)
	}

	return payments, nil
}
