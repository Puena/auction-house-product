package dto

// Query represent a query with key and generic value.
type Query[T any] struct {
	Key   string
	Value T
}

// QueryFindProduct represent product found query.
type QueryFindProduct struct {
	Query[QueryProduct]
}

// NewQueryFindProduct create a new product found query.
func NewQueryFindProduct(key string, value QueryProduct) QueryFindProduct {
	return QueryFindProduct{
		Query: Query[QueryProduct]{
			Key:   key,
			Value: value,
		},
	}
}

// QueryFindProducts represent products found query.
type QueryFindProducts struct {
	Query[QueryProducts]
}

// NewQueryFindProducts create a new products found query.
func NewQueryFindProducts(key string, value QueryProducts) QueryFindProducts {
	return QueryFindProducts{
		Query: Query[QueryProducts]{
			Key:   key,
			Value: value,
		},
	}
}
