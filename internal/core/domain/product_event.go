package domain

// Event represent event.
type Event[T any] struct {
	// Key used as idempotence key.
	Key   string
	Value T
}

// NewEvent create a new event.
func NewEvent[T any](key string, value T) Event[T] {
	return Event[T]{key, value}
}

// EventProductCreated represent product created event.
type EventProductCreated struct {
	Event[Product]
}

// NewEventProductCreated create a new product created event.
func NewEventProductCreated(key string, product Product) EventProductCreated {
	return EventProductCreated{NewEvent[Product](key, product)}
}

// EventProductUpdated represent product updated event.
type EventProductUpdated struct {
	Event[Product]
}

// NewEventProductUpdated create a new product updated event.
func NewEventProductUpdated(key string, product Product) EventProductUpdated {
	return EventProductUpdated{NewEvent[Product](key, product)}
}

// EventProductDeleted represent product deleted event.
type EventProductDeleted struct {
	Event[Product]
}

// NewEventProductDeleted create a new product deleted event.
func NewEventProductDeleted(key string, product Product) EventProductDeleted {
	return EventProductDeleted{NewEvent[Product](key, product)}
}

// EventProductFound represent product found event.
type EventProductFound struct {
	Event[Product]
}

// NewEventProductFound create a new product found event.
func NewEventProductFound(key string, product Product) EventProductFound {
	return EventProductFound{NewEvent[Product](key, product)}
}

// EventProductsFound represent products found event.
type EventProductsFound struct {
	Event[[]Product]
}

// NewEventProductsFound create a new products found event.
func NewEventProductsFound(key string, products []Product) EventProductsFound {
	return EventProductsFound{NewEvent[[]Product](key, products)}
}

// EventProductError represent internal error event.
type EventProductError struct {
	Event[ProductEventError]
}

// NewEventProductError create a new internal error event.
func NewEventProductError(key string, internalError ProductEventError) EventProductError {
	return EventProductError{NewEvent[ProductEventError](key, internalError)}
}
