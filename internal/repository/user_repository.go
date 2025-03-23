package repository

import (
	"context"
	"fmt"
	"petstore/internal/model"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(ctx context.Context, user model.User) (model.User, error)
	CreateBatch(ctx context.Context, users []model.User) ([]model.User, error)
	FindByUsername(ctx context.Context, username string) (model.User, error)
	Update(ctx context.Context, username string, user model.User) (model.User, error)
	Delete(ctx context.Context, username string) error
}

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepo{db: db}
}

func (u *userRepo) Create(ctx context.Context, user model.User) (model.User, error) {
	query := `
		INSERT INTO users (username, first_name, last_name, email, password, phone, user_status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id;
	`

	var newID int
	err := u.db.QueryRowContext(ctx, query,
		user.Username,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.Phone,
		user.UserStatus,
	).Scan(&newID)

	if err != nil {
		return user, fmt.Errorf("failed to insert user: %w", err)
	}

	user.ID = int64(newID)
	return user, nil
}

func (u *userRepo) CreateBatch(ctx context.Context, users []model.User) ([]model.User, error) {
	createdUsers := make([]model.User, 0, len(users))

	for _, user := range users {
		createdUser, err := u.Create(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("failed creating users: %w", err)
		}
		createdUsers = append(createdUsers, createdUser)
	}

	return createdUsers, nil
}

func (u *userRepo) FindByUsername(ctx context.Context, username string) (model.User, error) {
	query := `
		SELECT id, username, first_name, last_name, email, password, phone, user_status
		FROM users WHERE username = $1
	`

	var user model.User

	err := u.db.GetContext(ctx, &user, query, username)
	if err != nil {
		return user, fmt.Errorf("failed to find user by username %s: %w", username, err)
	}

	return user, nil
}

func (u *userRepo) Update(ctx context.Context, username string, user model.User) (model.User, error) {
	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, email = $3, password = $4, phone = $5, user_status = $6
		WHERE username = $7
	`

	_, err := u.db.ExecContext(ctx, query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.Phone,
		user.UserStatus,
		username,
	)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to update user: %w", err)
	}

	updatedUser, err := u.FindByUsername(ctx, username)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to retrieve updated user: %w", err)
	}

	return updatedUser, nil
}

func (u *userRepo) Delete(ctx context.Context, username string) error {
	query := `DELETE FROM users WHERE username = $1`

	_, err := u.db.ExecContext(ctx, query, username)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
