package orders

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	repo "github.com/joshbedo/golang-ecom/internal/adapters/postgres/sqlc"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrProductNoStock  = errors.New("product doesn't have enough stock")
)

type svc struct {
	repo *repo.Queries
	db   *pgx.Conn // Needed for transactions
}

func NewService(repo *repo.Queries, db *pgx.Conn) Service {
	return &svc{
		repo: repo,
		db:   db,
	}
}

func (s *svc) PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repo.Order, error) {
	// Validate payload - could be in handler
	if tempOrder.CustomerID == 0 {
		return repo.Order{}, fmt.Errorf("customer ID is required")
	}

	if len(tempOrder.Items) == 0 {
		return repo.Order{}, fmt.Errorf("at least one item is required")
	}

	// Create an order transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Order{}, err
	}
	defer tx.Rollback(ctx)

	qtx := s.repo.WithTx(tx)

	// Create an order
	order, err := qtx.CreateOrder(ctx, tempOrder.CustomerID)
	if err != nil {
		return repo.Order{}, err
	}

	// Look if product exists
	for _, item := range tempOrder.Items {
		product, err := qtx.FindProductByID(ctx, item.ProductID)
		if err != nil {
			return repo.Order{}, ErrProductNotFound
		}

		// Check and update stock - probably should go before creating order items
		newQuantity := product.Quantity - item.Quantity
		if newQuantity < 0 {
			return repo.Order{}, ErrProductNoStock
		}

		// Update the product's quantity
		_, err = qtx.UpdateProductQuantity(ctx, repo.UpdateProductQuantityParams{
			ID:       item.ProductID,
			Quantity: newQuantity,
		})
		if err != nil {
			return repo.Order{}, err
		}

		// Create order item
		_, err = qtx.CreateOrderItem(ctx, repo.CreateOrderItemParams{
			OrderID:    order.ID,
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			PriceCents: product.PriceInCents,
		})
		if err != nil {
			return repo.Order{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return repo.Order{}, err
	}

	return order, nil
}
