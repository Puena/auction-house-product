package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	logger "github.com/Puena/auction-house-logger"
	"github.com/Puena/auction-house-product/internal/core/domain"
	"github.com/jackc/pgx/v5/pgconn"
)

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *productRepository {
	return &productRepository{db: db}
}

// Actions
//
// CreateProduct create a new product.
func (p *productRepository) CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error) {
	const query = "INSERT INTO products (id, name, description, media, created_at, created_by) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, name, description, media, created_at, created_by"

	logger.Info().Str("product", fmt.Sprintf("%+v", product)).Msg("insert into products")
	res := p.db.QueryRowContext(ctx, query, &product.ID, &product.Name, &product.Description, &product.Media, &product.CreatedAt, &product.CreatedBy)
	if res.Err() != nil {
		return domain.Product{}, res.Err()
	}

	var createdProduct domain.Product
	if err := res.Scan(&createdProduct.ID, &createdProduct.Name, &createdProduct.Description, &createdProduct.Media, &createdProduct.CreatedAt, &createdProduct.CreatedBy); err != nil {
		return domain.Product{}, err
	}

	return createdProduct, nil
}

// UpdateProduct update a product.
func (p *productRepository) UpdateProduct(ctx context.Context, productID string, toUpdate domain.UpdateProduct) (domain.Product, error) {
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString("UPDATE products SET ")
	queryValues := make([]interface{}, 0)
	countValues := func() int {
		return len(queryValues) + 1
	}
	insertCommaIfValueExists := func() {
		if len(queryValues) > 1 {
			queryBuilder.WriteString(", ")
		}
	}

	if toUpdate.Name != "" {
		queryBuilder.WriteString(fmt.Sprintf("name = $%d", countValues()))
		queryValues = append(queryValues, toUpdate.Name)
	}
	if toUpdate.Description != "" {
		insertCommaIfValueExists()
		queryBuilder.WriteString(fmt.Sprintf("description = $%d", countValues()))
		queryValues = append(queryValues, toUpdate.Description)
	}
	if len(toUpdate.Media) > 0 {
		insertCommaIfValueExists()
		queryBuilder.WriteString(fmt.Sprintf("media = array(SELECT DISTINCT unnest(media || $%d))", countValues()))
		queryValues = append(queryValues, toUpdate.Media)
	}

	queryValues = append(queryValues, productID)
	queryBuilder.WriteString(fmt.Sprintf(" WHERE id = $%d RETURNING id, name, description, media, created_at, created_by", countValues()))

	logger.Info().Str("productID", productID).Str("update_data", fmt.Sprintf("%+v", toUpdate)).Msg("update products")
	row := p.db.QueryRowContext(ctx, queryBuilder.String(), queryValues...)
	if row.Err() != nil {
		return domain.Product{}, row.Err()
	}

	var updatedProduct domain.Product
	err := row.Scan(&updatedProduct.ID, &updatedProduct.Name, &updatedProduct.Description, &updatedProduct.Media, &updatedProduct.CreatedAt, &updatedProduct.CreatedBy)
	if err != nil {
		return domain.Product{}, err
	}

	return updatedProduct, nil
}

// DeleteProduct delete a product.
func (p *productRepository) DeleteProduct(ctx context.Context, productID string) (domain.Product, error) {
	const query = "DELETE FROM products WHERE id = $1 RETURNING id, name, description, media, created_at, created_by"

	logger.Info().Str("product_id", productID).Msg("delete from products")
	row := p.db.QueryRowContext(ctx, query, productID)
	if row.Err() != nil {
		return domain.Product{}, row.Err()
	}

	var deletedProduct domain.Product
	err := row.Scan(&deletedProduct.ID, &deletedProduct.Name, &deletedProduct.Description, &deletedProduct.Media, &deletedProduct.CreatedAt, &deletedProduct.CreatedBy)
	if err != nil {
		return domain.Product{}, err
	}

	return deletedProduct, nil
}

// FindOne find a product by id.
func (p *productRepository) FindOne(ctx context.Context, productID string) (domain.Product, error) {
	const query = "SELECT id, name, description, media, created_at, created_by FROM products WHERE id = $1"

	logger.Info().Str("product_id", productID).Msg("select from products")
	row := p.db.QueryRowContext(ctx, query, productID)
	if row.Err() != nil {
		return domain.Product{}, row.Err()
	}

	var product domain.Product
	err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Media, &product.CreatedAt, &product.CreatedBy)
	if err != nil {
		return domain.Product{}, err
	}

	return product, nil
}

// FindAll find all products.
func (p *productRepository) FindAll(ctx context.Context) (products []domain.Product, err error) {
	// TODO: specify default limit
	const query = "SELECT id, name, description, media, created_at, created_by FROM products"

	logger.Info().Msg("select from products")

	var rows *sql.Rows
	rows, err = p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cErr := rows.Close(); cErr != nil {
			err = errors.Join(err, cErr)
		}
	}()

	products = make([]domain.Product, 0)
	for rows.Next() {
		var product domain.Product
		err = rows.Scan(&product.ID, &product.Name, &product.Description, &product.Media, &product.CreatedAt, &product.CreatedBy)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

// Errors
//
// ConflictError check if error is a conflict error (unique constrain).
func (p *productRepository) ConflictError(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

// NotFoundError check if error is a not found error.
func (p *productRepository) NotFoundError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
