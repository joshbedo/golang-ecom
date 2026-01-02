// @todo(medium)(API) - maybe add gorm and repo.go with ORM connection to db

package products

import (
	"context"

	repo "github.com/joshbedo/golang-ecom/internal/adapters/postgres/sqlc"
)

// @todo(medium)(Testing) - could create in-memory layer for testing
type Service interface {
	ListProducts(ctx context.Context) ([]repo.Product, error)
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{repo: repo}
}

func (s *svc) ListProducts(ctx context.Context) ([]repo.Product, error) {
	return s.repo.ListProducts(ctx)
}
