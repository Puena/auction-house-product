package domain

// Event represent event.
type Event[T any] struct {
	Value T
}

// NewEvent create a new event.
func NewEvent[T any](value T) Event[T] {
	return Event[T]{value}
}

// EventProductCreated represent product created event.
type EventProductCreated struct {
	Event[Product]
}

// NewEventProductCreated create a new product created event.
func NewEventProductCreated(product Product) EventProductCreated {
	return EventProductCreated{NewEvent[Product](product)}
}

// EventProductUpdated represent product updated event.
type EventProductUpdated struct {
	Event[Product]
}

// NewEventProductUpdated create a new product updated event.
func NewEventProductUpdated(product Product) EventProductUpdated {
	return EventProductUpdated{NewEvent[Product](product)}
}

// EventProductDeleted represent product deleted event.
type EventProductDeleted struct {
	Event[Product]
}

// NewEventProductDeleted create a new product deleted event.
func NewEventProductDeleted(product Product) EventProductDeleted {
	return EventProductDeleted{NewEvent[Product](product)}
}

// EventProductFound represent product found event.
type EventProductFound struct {
	Event[Product]
}

// NewEventProductFound create a new product found event.
func NewEventProductFound(product Product) EventProductFound {
	return EventProductFound{NewEvent[Product](product)}
}

// EventProductsFound represent products found event.
type EventProductsFound struct {
	Event[[]Product]
}

// NewEventProductsFound create a new products found event.
func NewEventProductsFound(products []Product) EventProductsFound {
	return EventProductsFound{NewEvent[[]Product](products)}
}
