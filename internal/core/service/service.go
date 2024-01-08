package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Puena/auction-house-product/internal/core/domain"
	"github.com/Puena/auction-house-product/internal/core/dto"
	"github.com/Puena/auction-house-product/internal/core/ports"
	"github.com/oklog/ulid/v2"
)

const (
	eventKeyPrefixProductCreated = "created_product"
	eventKeyPrefixProductUpdated = "updated_product"
	eventKeyPrefixProductDeleted = "deleted_product"
	eventKeyPrefixProductFound   = "found_product"
	eventKeyPrefixProductsFound  = "found_products"
	eventKeyPrefixEventError     = "product_event_error"
)

type productService struct {
	database ports.ProductRepository
	publish  ports.PublishEventsRepository
}

// NewProductService creates a new product service.
func NewProductService(database ports.ProductRepository, publish ports.PublishEventsRepository) *productService {
	return &productService{
		database: database,
		publish:  publish,
	}
}

// CreateProduct creates a new product and save it in the database.
func (s *productService) CreateProduct(ctx context.Context, authUserID string, data dto.CommandCreateProduct) (dto.Product, error) {
	// validate
	if err := validateAuthUserID(authUserID); err != nil {
		return dto.Product{}, err
	}
	if err := validateCreateProductValue(data.Value); err != nil {
		return dto.Product{}, err
	}

	// create new product entity
	newProduct := domain.Product{
		ID:          ulid.Make().String(),
		Name:        data.Value.Name,
		Description: data.Value.Description,
		Media:       data.Value.Media,
		CreatedAt:   time.Now(),
		CreatedBy:   authUserID,
	}

	// save it to database
	createdProduct, err := s.database.CreateProduct(ctx, newProduct)
	if err != nil {
		return dto.Product{}, err
	}

	return mapDomainProductToDto(createdProduct), nil
}

// UpdateProduct updates a product and return updated product or error.
func (s *productService) UpdateProduct(ctx context.Context, authUserID string, data dto.CommandUpdateProduct) (dto.Product, error) {
	// validate
	if err := validateAuthUserID(authUserID); err != nil {
		return dto.Product{}, err
	}
	if err := validateUpdateProductValue(data.Value); err != nil {
		return dto.Product{}, err
	}
	if err := validateThatUserIDAndCreatedByAreTheSameID(authUserID, data.Value.CreatedBy); err != nil {
		return dto.Product{}, err
	}

	// update product entity
	updateProduct := domain.UpdateProduct{
		Name:        data.Value.Name,
		Description: data.Value.Description,
		Media:       data.Value.Media,
	}

	// save updated product to the database
	updatedProduct, err := s.database.UpdateProduct(ctx, data.Value.ID, authUserID, updateProduct)
	if err != nil {
		return dto.Product{}, err
	}

	return mapDomainProductToDto(updatedProduct), nil
}

// DeleteProduct deletes a product and return deleted product or error.
func (s *productService) DeleteProduct(ctx context.Context, authUserID string, data dto.CommandDeleteProduct) (dto.Product, error) {
	// validate
	if err := validateAuthUserID(authUserID); err != nil {
		return dto.Product{}, err
	}
	if err := validateDeleteProductValue(data.Value); err != nil {
		return dto.Product{}, err
	}
	if err := validateThatUserIDAndCreatedByAreTheSameID(authUserID, data.Value.CreatedBy); err != nil {
		return dto.Product{}, err
	}

	// delete product from database
	deletedProduct, err := s.database.DeleteProduct(ctx, data.Value.ID, authUserID)
	if err != nil {
		return dto.Product{}, err
	}

	return mapDomainProductToDto(deletedProduct), nil
}

// FindProduct finds a product and return found product or error if not found.
func (s *productService) FindProduct(ctx context.Context, authUserID string, data dto.QueryFindProduct) (dto.Product, error) {
	// validate
	if err := validateAuthUserID(authUserID); err != nil {
		return dto.Product{}, err
	}
	if err := validateProductID(data.Value.ID); err != nil {
		return dto.Product{}, err
	}

	// find product in database
	foundProduct, err := s.database.FindOne(ctx, data.Value.ID)
	if err != nil {
		return dto.Product{}, err
	}

	return mapDomainProductToDto(foundProduct), nil
}

// FindProducts finds products and return found products or error if not found.
func (s *productService) FindProducts(ctx context.Context, authUserID string, data dto.QueryFindProducts) ([]dto.Product, error) {
	// validate
	if err := validateAuthUserID(authUserID); err != nil {
		return []dto.Product{}, err
	}

	// find products in database
	foundProducts, err := s.database.FindAll(ctx)
	if err != nil {
		return []dto.Product{}, err
	}

	var products []dto.Product
	for _, product := range foundProducts {
		products = append(products, mapDomainProductToDto(product))
	}

	return products, nil
}

// PublishEventProductCreated publish event product created.
func (s *productService) PublishEventProductCreated(ctx context.Context, authUserID string, product dto.Product) error {
	eventKey := fmt.Sprintf("%s_%s", eventKeyPrefixProductCreated, product.ID)
	event := domain.NewEventProductCreated(eventKey, mapProductDtoToDomain(product))

	return s.publish.ProductCreated(ctx, authUserID, event)
}

// PublishEventProductUpdated publish event product updated.
func (s *productService) PublishEventProductUpdated(ctx context.Context, authUserID string, product dto.Product) error {
	eventKey := fmt.Sprintf("%s_%s_%s", eventKeyPrefixProductUpdated, product.ID, product.UpdatedAt)
	event := domain.NewEventProductUpdated(eventKey, mapProductDtoToDomain(product))

	return s.publish.ProductUpdated(ctx, authUserID, event)
}

// PublishEventProductDeleted publish event product deleted.
func (s *productService) PublishEventProductDeleted(ctx context.Context, authUserID string, product dto.Product) error {
	eventKey := fmt.Sprintf("%s_%s", eventKeyPrefixProductDeleted, product.ID)
	event := domain.NewEventProductDeleted(eventKey, mapProductDtoToDomain(product))

	return s.publish.ProductDeleted(ctx, authUserID, event)
}

// PublishEventProductFound publish event product found.
func (s *productService) PublishEventProductFound(ctx context.Context, authUserID string, product dto.Product) error {
	eventKey := fmt.Sprintf("%s_%s", eventKeyPrefixProductFound, ulid.Make().String())
	event := domain.NewEventProductFound(eventKey, mapProductDtoToDomain(product))

	return s.publish.ProductFound(ctx, authUserID, event)
}

// PublishEventProductsFound publish event products found.
func (s *productService) PublishEventProductsFound(ctx context.Context, authUserID string, products []dto.Product) error {
	var domainProducts []domain.Product
	for _, product := range products {
		domainProducts = append(domainProducts, mapProductDtoToDomain(product))
	}

	eventKey := fmt.Sprintf("%s_%s", eventKeyPrefixProductsFound, ulid.Make().String())
	event := domain.NewEventProductsFound(eventKey, domainProducts)

	return s.publish.ProductsFound(ctx, authUserID, event)
}

// PublishProductError publish event product error.
func (s *productService) PublishProductError(ctx context.Context, authUserID string, productError dto.ProductEventError) error {
	eventKey := fmt.Sprintf("%s_%s_%d", eventKeyPrefixEventError, productError.Reference_event_key, productError.Time.UnixMicro())

	return s.publish.ProductError(ctx, authUserID, domain.NewEventProductError(eventKey, mapDtoInternalEventErrorToDomain(productError)))
}

// ValidationError returns true if the error is a service validation error.
func (s *productService) ValidationError(err error) bool {
	var validationErr serviceValidationError
	return err != nil && errors.Is(err, &validationErr)
}

// UniqueConstrainError returns true if the error is a service unique constrain error.
func (s *productService) UniqueConstrainError(err error) bool {
	return s.database.ConflictError(err)
}
