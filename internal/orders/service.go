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

	// Process each order item
	for _, item := range tempOrder.Items {
		// Verify product exists
		product, err := qtx.FindProductByID(ctx, item.ProductID)
		if err != nil {
			return repo.Order{}, ErrProductNotFound
		}

		// Create order item with 'pending' status
		orderItem, err := qtx.CreateOrderItem(ctx, repo.CreateOrderItemParams{
			OrderID:    order.ID,
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			PriceCents: product.PriceInCents,
			Status:     OrderItemStatusPending.Text(),
		})
		if err != nil {
			return repo.Order{}, err
		}

		// Lock and check stock availability
		product, err = qtx.FindProductByIDForUpdate(ctx, item.ProductID)
		if err != nil {
			return repo.Order{}, ErrProductNotFound
		}

		// Check if we have enough stock
		hasStock := product.Quantity >= item.Quantity

		if hasStock {
			// Stock available: decrement product quantity and update status to fulfilled
			newQuantity := product.Quantity - item.Quantity
			_, err = qtx.UpdateProductQuantity(ctx, repo.UpdateProductQuantityParams{
				ID:       item.ProductID,
				Quantity: newQuantity,
			})
			if err != nil {
				return repo.Order{}, err
			}
			_, err = qtx.UpdateOrderItemStatus(ctx, repo.UpdateOrderItemStatusParams{
				ID:     orderItem.ID,
				Status: OrderItemStatusFulfilled.Text(),
			})
			if err != nil {
				return repo.Order{}, err
			}
		} else {
			// Stock not available: update status to backordered (don't decrement stock)
			_, err = qtx.UpdateOrderItemStatus(ctx, repo.UpdateOrderItemStatusParams{
				ID:     orderItem.ID,
				Status: OrderItemStatusBackordered.Text(),
			})
			if err != nil {
				return repo.Order{}, err
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return repo.Order{}, err
	}

	return order, nil
}
