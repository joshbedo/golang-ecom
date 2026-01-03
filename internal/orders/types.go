package orders

import (
	"context"

	repo "github.com/joshbedo/golang-ecom/internal/adapters/postgres/sqlc"
)

type OrderItemStatus string

const (
	OrderItemStatusPending     OrderItemStatus = "pending"
	OrderItemStatusFulfilled   OrderItemStatus = "fulfilled"
	OrderItemStatusBackordered OrderItemStatus = "backordered"
)

type orderItem struct {
	ProductID int64 `json:"productId"`
	Quantity  int32 `json:"quantity"`
	Status    OrderItemStatus
}

type createOrderParams struct {
	CustomerID int64       `json:"customerId"`
	Items      []orderItem `json:"items"`
}

type Service interface {
	PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repo.Order, error)
}
