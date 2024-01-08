package postgres

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

// NewProductRepository create a new product repository.
func NewProductRepository(db *sql.DB) *productRepository {
	return &productRepository{db: db}
}

// Actions
//
// CreateProduct create a new product.
func (p *productRepository) CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error) {
	const query = "INSERT INTO products (id, name, description, media, created_at, created_by) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, name, description, media, created_at, created_by"

	logger.Info().Str("product", fmt.Sprintf("%+v", product)).Msg("insert into products")
	logger.Debug().Str("query", query).Msg("insert into product")

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

// UpdateProduct update a product if at least one of toUpdate value provided.
func (p *productRepository) UpdateProduct(ctx context.Context, productID string, userID string, toUpdate domain.UpdateProduct) (domain.Product, error) {
	// arrange
	queryBuilder := strings.Builder{}
	queryValues := make([]interface{}, 0)
	getLastValuePosition := func() int {
		return len(queryValues) + 1
	}
	insertComma := func() {
		if len(queryValues) > 1 {
			queryBuilder.WriteString(",")
		}
	}

	// building query
	queryBuilder.WriteString("UPDATE products SET")
	if toUpdate.Name != "" {
		queryBuilder.WriteString(fmt.Sprintf(" name = $%d", getLastValuePosition()))
		queryValues = append(queryValues, toUpdate.Name)
	}
	if toUpdate.Description != "" {
		insertComma()
		queryBuilder.WriteString(fmt.Sprintf(" description = $%d", getLastValuePosition()))
		queryValues = append(queryValues, toUpdate.Description)
	}
	if len(toUpdate.Media) > 0 {
		insertComma()
		queryBuilder.WriteString(fmt.Sprintf(" media = array(SELECT DISTINCT unnest(media || $%d))", getLastValuePosition()))
		queryValues = append(queryValues, toUpdate.Media)
	}
	if len(queryValues) > 0 {
		insertComma()
		queryBuilder.WriteString(fmt.Sprintf(" updated_at = $%d", getLastValuePosition()))
		queryValues = append(queryValues, toUpdate.UpdatedAt)
	}
	queryValues = append(queryValues, productID)
	queryBuilder.WriteString(fmt.Sprintf(" WHERE id = $%d", getLastValuePosition()))
	queryValues = append(queryValues, userID)
	queryBuilder.WriteString(fmt.Sprintf(" AND created_by = $%d", getLastValuePosition()))
	queryBuilder.WriteString(" RETURNING id, name, description, media, created_at, updated_at, created_by")
	// Done building query.
	query := queryBuilder.String()

	logger.Info().Str("productID", productID).Str("update_data", fmt.Sprintf("%+v", toUpdate)).Msg("update products")
	logger.Debug().Str("query", query).Msg("update products")

	row := p.db.QueryRowContext(ctx, query, queryValues...)
	if row.Err() != nil {
		return domain.Product{}, row.Err()
	}

	var updatedProduct domain.Product
	err := row.Scan(&updatedProduct.ID, &updatedProduct.Name, &updatedProduct.Description, &updatedProduct.Media, &updatedProduct.CreatedAt, &updatedProduct.UpdatedAt, &updatedProduct.CreatedBy)
	if err != nil {
		return domain.Product{}, err
	}

	return updatedProduct, nil
}

// DeleteProduct delete a product.
func (p *productRepository) DeleteProduct(ctx context.Context, productID string, userID string) (domain.Product, error) {
	const query = "DELETE FROM products WHERE id = $1 AND created_by = $2 RETURNING id, name, description, media, created_at, updated_at, created_by"

	logger.Info().Str("product_id", productID).Msg("delete from products")
	logger.Debug().Str("query", query).Msg("delete from products")
	row := p.db.QueryRowContext(ctx, query, productID, userID)
	if row.Err() != nil {
		return domain.Product{}, row.Err()
	}

	var deletedProduct domain.Product
	err := row.Scan(&deletedProduct.ID, &deletedProduct.Name, &deletedProduct.Description, &deletedProduct.Media, &deletedProduct.CreatedAt, &deletedProduct.UpdatedAt, &deletedProduct.CreatedBy)
	if err != nil {
		return domain.Product{}, err
	}

	return deletedProduct, nil
}

// FindOne find a product by id.
func (p *productRepository) FindOne(ctx context.Context, productID string) (domain.Product, error) {
	const query = "SELECT id, name, description, media, created_at, updated_at, created_by FROM products WHERE id = $1"

	logger.Info().Str("product_id", productID).Msg("select from products")
	logger.Debug().Str("query", query).Msg("select from products")

	row := p.db.QueryRowContext(ctx, query, productID)
	if row.Err() != nil {
		return domain.Product{}, row.Err()
	}

	var product domain.Product
	err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Media, &product.CreatedAt, &product.UpdatedAt, &product.CreatedBy)
	if err != nil {
		return domain.Product{}, err
	}

	return product, nil
}

// FindAll find all products.
func (p *productRepository) FindAll(ctx context.Context) (products []domain.Product, err error) {
	// TODO: handle hardcoded limit
	const query = "SELECT id, name, description, media, created_at, updated_at, created_by FROM products LIMIT 50"

	logger.Info().Msg("select all from products")
	logger.Debug().Str("query", query).Msg("select all from products")

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

	for rows.Next() {
		var product domain.Product
		err = rows.Scan(&product.ID, &product.Name, &product.Description, &product.Media, &product.CreatedAt, &product.UpdatedAt, &product.CreatedBy)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	// if products is empty, return error
	if len(products) == 0 {
		return nil, sql.ErrNoRows
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
