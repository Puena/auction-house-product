package service

import (
	"fmt"

	"github.com/Puena/auction-house-product/internal/core/dto"
)

type serviceValidationError struct {
	message string
}

func (e *serviceValidationError) Error() string {
	return fmt.Sprintf("service validation error: %s", e.message)
}

func newServiceValidationError(message string) error {
	return &serviceValidationError{message: message}
}

func validateAuthUserID(authUserID string) error {
	if authUserID == "" {
		return newServiceValidationError("auth user id is empty")
	}
	return nil
}

func validateCreateProductValue(value dto.CreateProduct) error {
	if value.Name == "" {
		return newServiceValidationError("name is empty")
	}
	if value.Description == "" {
		return newServiceValidationError("description is empty")
	}

	return nil
}

func validateUpdateProductValue(value dto.UpdateProduct) error {
	someOneIsNotEmpty := false
	if value.Name != "" {
		someOneIsNotEmpty = true
	}
	if value.Description != "" {
		someOneIsNotEmpty = true
	}
	if len(value.Media) > 0 {
		someOneIsNotEmpty = true
	}

	if !someOneIsNotEmpty {
		return newServiceValidationError("all values are empty, nothing to update")
	}

	return nil
}

func validateDeleteProductValue(value dto.DeleteProduct) error {
	if value.ID == "" {
		return newServiceValidationError("id is empty")
	}

	if value.CreatedBy == "" {
		return newServiceValidationError("created by is empty")
	}

	return nil
}

func validateThatUserIDAndCreatedByAreTheSameID(userID, createdBy string) error {
	if userID != createdBy {
		return newServiceValidationError("you are not owner of this product")
	}
	return nil
}

func validateProductID(productID string) error {
	if productID == "" {
		return newServiceValidationError("product id is empty")
	}
	return nil
}
