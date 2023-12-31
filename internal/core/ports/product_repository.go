package ports

import (
	"context"

	"github.com/Puena/autcion-house-product/internal/core/domain"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, product domain.Product) error
}
