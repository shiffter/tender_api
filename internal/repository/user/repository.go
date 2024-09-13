package user

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"tender/internal/repository"
	repoUserModel "tender/internal/repository/model"
	"tender/internal/repository/model/converter"
	serviceUserModel "tender/internal/service/model"
	"tender/pkg/stderrs"
)

type repo struct {
	db *pgxpool.Pool
}

func NewUserRepos(pool *pgxpool.Pool) repository.UsersRepos {
	return &repo{db: pool}
}

func (r *repo) Get(ctx context.Context, username string) (*serviceUserModel.User, error) {
	var (
		sql = `SELECT id, username, first_name, last_name, created_at, updated_at
				FROM employee
				WHERE username = $1`
		user repoUserModel.User
	)

	err := r.db.QueryRow(ctx, sql, username).
		Scan(&user.ID, &user.Username, &user.FirstName,
			&user.LastName, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, stderrs.ErrUserNotFound
		}

		return nil, err
	}

	return converter.ToServiceUserFromRepo(&user), nil
}

func (r *repo) GetByID(ctx context.Context, id uuid.UUID) (*serviceUserModel.User, error) {
	var (
		sql = `SELECT id, username, first_name, last_name, created_at, updated_at
				FROM employee
				WHERE id = $1`
		user repoUserModel.User
	)

	err := r.db.QueryRow(ctx, sql, id).
		Scan(&user.ID, &user.Username, &user.FirstName,
			&user.LastName, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, stderrs.ErrUserNotFound
		}

		return nil, err
	}

	return converter.ToServiceUserFromRepo(&user), nil
}
