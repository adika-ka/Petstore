package repository

import (
	"context"
	"database/sql"
	"fmt"
	"petstore/internal/model"

	"github.com/jmoiron/sqlx"
)

type OrderRepository interface {
	Create(ctx context.Context, order model.Order) (model.Order, error)
	FindByID(ctx context.Context, orderID int) (model.Order, error)
	Delete(ctx context.Context, orderID int) error
	GetInventory(ctx context.Context) (map[string]int, error)
}

type orderRepo struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) OrderRepository {
	return &orderRepo{db: db}
}

func (r *orderRepo) Create(ctx context.Context, order model.Order) (model.Order, error) {
	query := `
		INSERT INTO orders (pet_id, quantity, ship_date, status, complete)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`

	var newID int
	err := r.db.QueryRowContext(ctx, query,
		order.PetID,
		order.Quantity,
		order.ShipDate,
		order.Status,
		order.Complete,
	).Scan(&newID)

	if err != nil {
		return order, fmt.Errorf("failed to insert order: %w", err)
	}

	order.ID = newID
	return order, nil
}

func (r *orderRepo) FindByID(ctx context.Context, orderID int) (model.Order, error) {
	query := `
		SELECT id, pet_id, quantity, ship_date, status, complete
		FROM orders WHERE id = $1
	`
	var order model.Order

	err := r.db.GetContext(ctx, &order, query, orderID)
	if err != nil {
		return order, fmt.Errorf("failed to find order by id: %w", err)
	}

	return order, nil
}

func (r *orderRepo) Delete(ctx context.Context, orderID int) error {
	query := `DELETE FROM orders WHERE id = $1`

	status, err := r.GetStatusByID(ctx, orderID)
	if err != nil {
		return err
	}

	if status == "delivered" {
		return fmt.Errorf("cannot delete a completed order")
	}

	_, err = r.db.ExecContext(ctx, query, orderID)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}
	return nil
}

func (r *orderRepo) GetInventory(ctx context.Context) (map[string]int, error) {
	query := `
		SELECT status, COUNT(*) as count
		FROM orders
		GROUP BY status
	`

	type statusCount struct {
		Status string `db:"status"`
		Count  int    `db:"count"`
	}

	var res []statusCount
	err := r.db.SelectContext(ctx, &res, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}

	inventory := make(map[string]int)
	for _, row := range res {
		inventory[row.Status] = row.Count
	}
	return inventory, nil
}

func (r *orderRepo) GetStatusByID(ctx context.Context, orderID int) (string, error) {
	var status string
	err := r.db.QueryRowContext(ctx, `SELECT status FROM orders WHERE id = $1`, orderID).Scan(&status)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("order not found")
	}
	return status, nil
}
