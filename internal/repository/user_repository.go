package repository

import (
	"context"
	"errors"
	"time"

	"github.com/DmitriyV003/bonus/internal/applicationerrors"
	"github.com/DmitriyV003/bonus/internal/models"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
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
	if err != nil && !errors.Is(err, applicationerrors.ErrNotFound) {
		return applicationerrors.ErrInternalServer
	}

	if dbUser != nil {
		return applicationerrors.ErrConflict
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return applicationerrors.ErrNotFound
	}

	_, err = users.db.Exec(ctx, sql, user.Login, user.Password, user.CreatedAt)

	if err != nil {
		return applicationerrors.ErrInternalServer
	}

	return nil
}

func (users *UserRepository) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	sql := `SELECT id, login, password, created_at FROM users WHERE login = $1`
	var user models.User

	row := users.db.QueryRow(ctx, sql, login)

	err := row.Scan(&user.ID, &user.Login, &user.Password, &user.CreatedAt)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, applicationerrors.ErrNotFound
	}

	return &user, nil
}

func (users *UserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	sql := `SELECT id, login, balance, created_at FROM users WHERE id = $1`
	var user models.User

	row := users.db.QueryRow(ctx, sql, id)

	err := row.Scan(&user.ID, &user.Login, &user.Balance, &user.CreatedAt)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, applicationerrors.ErrNotFound
	}

	return &user, nil
}

func (users *UserRepository) UpdateBalance(ctx context.Context, user *models.User) error {
	sql := `UPDATE users SET balance = $1, updated_at = $2 WHERE id = $3`

	_, err := users.db.Exec(ctx, sql, user.Balance, time.Now(), user.ID)
	if err != nil {
		return err
	}

	return nil
}
