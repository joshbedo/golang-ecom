// @todo(medium)(API) - maybe add gorm and repo.go with ORM connection to db

package products

import (
	"context"

	repo "github.com/joshbedo/golang-ecom/internal/adapters/postgres/sqlc"
)

type Service interface {
	ListProducts(ctx context.Context) ([]repo.Product, error)
	GetProductByID(ctx context.Context, id int64) (repo.Product, error)
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

func (s *svc) GetProductByID(ctx context.Context, id int64) (repo.Product, error) {
	return s.repo.GetProductByID(ctx, id)
}
