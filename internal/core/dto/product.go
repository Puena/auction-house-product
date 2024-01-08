package dto

import (
	"fmt"
	"time"
)

type Product struct {
	ID          string
	Name        string
	Description string
	Media       []string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	CreatedBy   string
}

type CreateProduct struct {
	Name        string
	Description string
	Media       []string
	CreatedBy   string
}

type UpdateProduct struct {
	ID          string
	Name        string
	Description string
	Media       []string
	CreatedBy   string
}

type DeleteProduct struct {
	ID        string
	CreatedBy string
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

func (p *ProductEventError) Error() string {
	return fmt.Sprintf("error occured at stream: %s, consumer: %s, subject: %s, msg_id: %s, code: %d, message: %s, at: %s", p.StreamName, p.ConsumerName, p.Subject, p.Reference_event_key, p.Code, p.Message, p.Time)
}

type QueryProduct struct {
	ID string
}

type QueryProducts struct {
}
