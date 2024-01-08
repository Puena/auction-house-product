package ports

import (
	"context"

	"github.com/Puena/auction-house-product/internal/core/domain"
)

// ProductRepository represent product repository port.
type ProductRepository interface {
	//
	// Actions
	//
	// CreateProduct create a new product.
	CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error)
	// UpdateProduct update a product.
	UpdateProduct(ctx context.Context, productID string, userID string, toUpdate domain.UpdateProduct) (domain.Product, error)
	// DeleteProduct delete a product.
	DeleteProduct(ctx context.Context, productID string, userID string) (domain.Product, error)
	// FindOne find a product by id.
	FindOne(ctx context.Context, productID string) (domain.Product, error)
	// FindAll find all products.
	FindAll(ctx context.Context) ([]domain.Product, error)
	//
	// Errors
	//
	// ConflictError check if error is a conflict error (unique constrain).
	ConflictError(err error) bool
	// NotFoundError check if error is a not found error.
	NotFoundError(err error) bool
}
