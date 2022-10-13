package repository

import (
	"context"
	"errors"
	"github.com/DmitriyV003/bonus/cmd/gophermart/application_errors"
	"github.com/DmitriyV003/bonus/cmd/gophermart/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: pool}
}

func (users *UserRepository) Create(ctx context.Context, user *models.User) error {
	sql := `INSERT INTO users (login, password, created_at) VALUES ($1, $2, $3)`

	dbUser, err := users.GetByLogin(ctx, user.Login)
	if err != nil && !errors.Is(err, application_errors.ErrNotFound) {
		return application_errors.ErrInternalServer
	}

	if dbUser != nil {
		return application_errors.ErrConflict
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return application_errors.ErrNotFound
	}

	_, err = users.db.Exec(ctx, sql, user.Login, user.Password, user.CreatedAt)

	if err != nil {
		return application_errors.ErrInternalServer
	}

	return nil
}

func (users *UserRepository) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	sql := `SELECT id, login, password, created_at FROM users WHERE login = $1`
	var user models.User

	row := users.db.QueryRow(ctx, sql, login)

	err := row.Scan(&user.Id, &user.Login, &user.Password, &user.CreatedAt)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, application_errors.ErrNotFound
	}

	return &user, nil
}
