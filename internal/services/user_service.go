package services

import (
	"context"
	"fmt"
	"github.com/DmitriyV003/bonus/internal/application_errors"
	"github.com/DmitriyV003/bonus/internal/container"
	"github.com/DmitriyV003/bonus/internal/models"
	"github.com/DmitriyV003/bonus/internal/requests"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

type UserService struct {
	container      *container.Container
	user           *models.User
	validator      OrderValidator
	paymentService *PaymentService
}

func NewUserService(container *container.Container, user *models.User, validator OrderValidator, paymentService *PaymentService) *UserService {
	return &UserService{
		container:      container,
		user:           user,
		validator:      validator,
		paymentService: paymentService,
	}
}

func (u *UserService) Create(request *requests.RegistrationRequest, jwtSecret string) (*Token, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(request.Password), 14)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Login:     request.Login,
		Password:  string(bytes),
		CreatedAt: time.Now(),
	}
	err = u.container.Users.Create(context.Background(), &user)
	if err != nil {
		return nil, fmt.Errorf("unable to create user in db: %w", err)
	}

	dbUser, err := u.container.Users.GetByLogin(context.Background(), request.Login)
	if err != nil {
		return nil, fmt.Errorf("error to get user by login: %w", err)
	}

	authService := NewAuthService(u.container, jwtSecret)
	token, err := authService.LoginByUser(dbUser)
	if err != nil {
		return nil, fmt.Errorf("login user programmly: %w", err)
	}

	return token, nil
}

func (u *UserService) Withdraw(orderNumber string, sum float64) error {
	parsedOderNumber, err := strconv.ParseInt(orderNumber, 10, 64)
	if err != nil {
		return err
	}

	isValid := u.validator.Validate(parsedOderNumber)
	if !isValid {
		return fmt.Errorf("unable to validate order number %w", application_errors.ErrInvalidOrderNumber)
	}

	sumToWithdraw := int64(sum * 10000)
	if u.user.Balance < sumToWithdraw {
		return fmt.Errorf("unable to validate order number %w", application_errors.ErrLowUserABalance)
	}

	err = u.paymentService.CreateWithdrawPayment(u.user, sumToWithdraw, orderNumber)
	if err != nil {
		return fmt.Errorf("unable ro create payment for withdraw %w", err)
	}

	return nil
}

func (u *UserService) AllWithdrawsByUser(user *models.User) ([]*models.Payment, error) {
	payments, err := u.container.Payments.GetWithdrawsByUser(context.Background(), user)
	if err != nil {
		return nil, fmt.Errorf("error to get all withdraws: %w", err)
	}

	return payments, nil
}
