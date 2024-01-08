package dto

// Command represent a command with key and generic value.
type Command[T any] struct {
	Key   string
	Value T
}

// CommandCreateProduct represent product created command.
type CommandCreateProduct struct {
	Command[CreateProduct]
}

// NewCommandCreateProduct create a new product created command.
func NewCommandCreateProduct(key string, value CreateProduct) CommandCreateProduct {
	return CommandCreateProduct{
		Command: Command[CreateProduct]{
			Key:   key,
			Value: value,
		},
	}
}

// CommandUpdateProduct represent product updated command.
type CommandUpdateProduct struct {
	Command[UpdateProduct]
}

// NewCommandUpdateProduct create a new product updated command.
func NewCommandUpdateProduct(key string, value UpdateProduct) CommandUpdateProduct {
	return CommandUpdateProduct{
		Command: Command[UpdateProduct]{
			Key:   key,
			Value: value,
		},
	}
}

// CommandDeleteProduct represent product deleted command.
type CommandDeleteProduct struct {
	Command[DeleteProduct]
}

// NewCommandDeleteProduct create a new product deleted command.
func NewCommandDeleteProduct(key string, value DeleteProduct) CommandDeleteProduct {
	return CommandDeleteProduct{
		Command: Command[DeleteProduct]{
			Key:   key,
			Value: value,
		},
	}
}
