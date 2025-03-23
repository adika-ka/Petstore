package service

import (
	"context"
	"fmt"
	"petstore/internal/config"
	"petstore/internal/model"
	"petstore/internal/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, user model.User) (model.User, error)
	CreateUserBatch(ctx context.Context, users []model.User) ([]model.User, error)
	FindUserByUsername(ctx context.Context, username string) (model.User, error)
	UpdateUser(ctx context.Context, username string, user model.User) (model.User, error)
	DeleteUser(ctx context.Context, username string) error
	Login(ctx context.Context, username, password string) (string, error)
	Logout(ctx context.Context) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (u *userService) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = hashedPassword

	return u.repo.Create(ctx, user)
}

func (u *userService) CreateUserBatch(ctx context.Context, users []model.User) ([]model.User, error) {
	return u.repo.CreateBatch(ctx, users)
}

func (u *userService) FindUserByUsername(ctx context.Context, username string) (model.User, error) {
	return u.repo.FindByUsername(ctx, username)
}

func (u *userService) UpdateUser(ctx context.Context, username string, user model.User) (model.User, error) {
	return u.repo.Update(ctx, username, user)
}

func (u *userService) DeleteUser(ctx context.Context, username string) error {
	return u.repo.Delete(ctx, username)
}

func (u *userService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := u.repo.FindByUsername(ctx, username)
	if err != nil {
		return "", fmt.Errorf("user with username %s not found: %w", username, err)
	}
	if err := checkPasswordHash(user.Password, password); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	_, token, err := config.TokenAuth.Encode(map[string]interface{}{
		"username": user.Username,
	})
	if err != nil {
		return "", fmt.Errorf("failed generating token: %w", err)
	}

	return token, nil
}

func(u *userService)Logout(ctx context.Context) error{
	return nil
}
