package domain

// Product represent product struct.
type Product struct {
	ID          string
	Name        string
	Description string
	Media       []string
	CreatedAt   string
	CreatedBy   string
}

// UpdateProduct represent available fields to update product.
type UpdateProduct struct {
	Name        string
	Description string
	Media       []string
}
