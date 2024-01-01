package ports

import (
	"context"

	"github.com/Puena/auction-house-product/internal/core/domain"
)

// PublishEventsRepository represent publish events repository.
type PublishEventsRepository interface {
	// Publish event product created.
	ProductCreated(ctx context.Context, event domain.EventProductCreated) error
	// Publish event product updated.
	ProductUpdated(ctx context.Context, event domain.EventProductUpdated) error
	// Publish event product deleted.
	ProductDeleted(ctx context.Context, event domain.EventProductDeleted) error
	// Publish event product found.
	ProductFound(ctx context.Context, event domain.EventProductFound) error
	// Publish event products found.
	ProductsFound(ctx context.Context, event domain.EventProductsFound) error
}
