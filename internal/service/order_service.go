package service

import (
	"context"
	"petstore/internal/model"
	"petstore/internal/repository"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order model.Order) (model.Order, error)
	FindOrderByID(ctx context.Context, orderID int) (model.Order, error)
	DeleteOrder(ctx context.Context, orderID int) error
	GetInventory(ctx context.Context) (map[string]int, error)
}

type orderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) OrderService {
	return &orderService{repo: repo}
}

func (o *orderService) CreateOrder(ctx context.Context, order model.Order) (model.Order, error) {
	return o.repo.Create(ctx, order)
}

func (o *orderService) FindOrderByID(ctx context.Context, orderID int) (model.Order, error) {
	return o.repo.FindByID(ctx, orderID)
}

func (o *orderService) DeleteOrder(ctx context.Context, orderID int) error {
	return o.repo.Delete(ctx, orderID)
}

func (o *orderService) GetInventory(ctx context.Context) (map[string]int, error) {
	return o.repo.GetInventory(ctx)
}
