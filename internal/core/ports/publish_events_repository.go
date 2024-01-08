package ports

import (
	"context"

	"github.com/Puena/auction-house-product/internal/core/domain"
)

// PublishEventsRepository represent publish events repository.
type PublishEventsRepository interface {
	// Publish event product created.
	ProductCreated(ctx context.Context, userID string, event domain.EventProductCreated) error
	// Publish event product updated.
	ProductUpdated(ctx context.Context, userID string, event domain.EventProductUpdated) error
	// Publish event product deleted.
	ProductDeleted(ctx context.Context, userID string, event domain.EventProductDeleted) error
	// Publish event product found.
	ProductFound(ctx context.Context, userID string, event domain.EventProductFound) error
	// Publish event products found.
	ProductsFound(ctx context.Context, userID string, event domain.EventProductsFound) error
	// ProductError publish product error.
	ProductError(ctx context.Context, userID string, event domain.EventProductError) error
}
