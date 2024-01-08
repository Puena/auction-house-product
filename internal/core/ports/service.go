package ports

import (
	"context"

	"github.com/Puena/auction-house-product/internal/core/dto"
)

// Service is the interface that wraps the basic methods of the service.
type Service interface {
	// CreateProduct creates a product and return created product or error.
	CreateProduct(ctx context.Context, authUserID string, data dto.CommandCreateProduct) (dto.Product, error)
	// UpdateProduct updates a product and return updated product or error.
	UpdateProduct(ctx context.Context, authUserID string, data dto.CommandUpdateProduct) (dto.Product, error)
	// DeleteProduct deletes a product and return deleted product or error.
	DeleteProduct(ctx context.Context, authUserID string, data dto.CommandDeleteProduct) (dto.Product, error)
	// FindProduct finds a product and return found product or error.
	FindProduct(ctx context.Context, authUserID string, data dto.QueryFindProduct) (dto.Product, error)
	// FindProducts finds products and return found products or error.
	FindProducts(ctx context.Context, authUserID string, data dto.QueryFindProducts) ([]dto.Product, error)
	// PublishEventProductCreated publish event product created.
	PublishEventProductCreated(ctx context.Context, authUserID string, product dto.Product) error
	// PublishEventProductUpdated publish event product updated.
	PublishEventProductUpdated(ctx context.Context, authUserID string, product dto.Product) error
	// PublishEventProductDeleted publish event product deleted.
	PublishEventProductDeleted(ctx context.Context, authUserID string, product dto.Product) error
	// PublishEventProductFound publish event product found.
	PublishEventProductFound(ctx context.Context, authUserID string, product dto.Product) error
	// PublishEventProductsFound publish event products found, limit 50 products.
	PublishEventProductsFound(ctx context.Context, authUserID string, products []dto.Product) error
	// PublishEventProductError publish event product error.
	PublishEventProductError(ctx context.Context, authUserID string, productError dto.ProductEventError) error
	/*
		Errors
	*/
	// ValidationError returns true if the error is a validation error.
	ValidationError(err error) bool
	// UniqueConstrainError returns true if the error is a unique constrain error.
	UniqueConstrainError(err error) bool
	// NotFoundError returns true if the error is a not found error.
	NotFoundError(err error) bool
}
