package service

import (
	"github.com/Puena/auction-house-product/internal/core/domain"
	"github.com/Puena/auction-house-product/internal/core/dto"
)

func mapProductDtoToDomain(product dto.Product) domain.Product {
	return domain.Product{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Media:       product.Media,
		CreatedAt:   product.CreatedAt,
		CreatedBy:   product.CreatedBy,
	}
}

func mapDomainProductToDto(product domain.Product) dto.Product {
	return dto.Product{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Media:       product.Media,
		CreatedAt:   product.CreatedAt,
		CreatedBy:   product.CreatedBy,
	}
}

func mapDtoInternalEventErrorToDomain(error dto.ProductEventError) domain.ProductEventError {
	return domain.ProductEventError{
		StreamName:          error.StreamName,
		ConsumerName:        error.ConsumerName,
		Subject:             error.Subject,
		Reference_event_key: error.Reference_event_key,
		Message:             error.Message,
		Code:                error.Code,
		Data:                error.Data,
		Headers:             error.Headers,
		Time:                error.Time,
	}
}
