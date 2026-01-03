package orders

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	repo "github.com/joshbedo/golang-ecom/internal/adapters/postgres/sqlc"
)

type OrderItemStatus string

const (
	OrderItemStatusPending     OrderItemStatus = "pending"
	OrderItemStatusFulfilled   OrderItemStatus = "fulfilled"
	OrderItemStatusBackordered OrderItemStatus = "backordered"
)

// Text converts OrderItemStatus to pgtype.Text for database operations
func (s OrderItemStatus) Text() pgtype.Text {
	return pgtype.Text{String: string(s), Valid: true}
}

type orderItem struct {
	ProductID int64 `json:"productId"`
	Quantity  int32 `json:"quantity"`
}

type createOrderParams struct {
	CustomerID int64       `json:"customerId"`
	Items      []orderItem `json:"items"`
}

type Service interface {
	PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repo.Order, error)
}
