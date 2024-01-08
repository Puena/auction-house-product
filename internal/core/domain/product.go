package domain

import "time"

// Product represent product struct.
type Product struct {
	ID          string
	Name        string
	Description string
	Media       []string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	CreatedBy   string
}

// UpdateProduct represent available fields to update product.
type UpdateProduct struct {
	Name        string    // optional
	Description string    // optional
	Media       []string  // optional
	UpdatedAt   time.Time // required
}

type ProductEventError struct {
	StreamName          string
	ConsumerName        string
	Subject             string
	Reference_event_key string
	Message             string
	Code                int
	Data                []byte
	Headers             string
	Time                time.Time
}
