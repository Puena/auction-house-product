package nats

import (
	"github.com/Puena/auction-house-product/internal/core/domain"
	"github.com/Puena/auction-house-product/internal/core/dto"
	"github.com/Puena/auction-messages-golang"
)

func mapCommandCreateProductToDto(msg *auction.CommandCreateProduct) dto.CommandCreateProduct {
	value := dto.CreateProduct{
		Name:        msg.Value.Name,
		Description: msg.Value.Description,
		Media:       msg.Value.Media,
		CreatedBy:   msg.Value.CreatedBy,
	}
	return dto.NewCommandCreateProduct(msg.Key, value)
}

func mapCommandUpdateProductToDto(msg *auction.CommandUpdateProduct) dto.CommandUpdateProduct {
	value := dto.UpdateProduct{
		Name:        msg.Value.Name,
		Description: msg.Value.Description,
		Media:       msg.Value.Media,
	}
	return dto.NewCommandUpdateProduct(msg.Key, value)
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

func mapCommandDeleteProductToDto(msg *auction.CommandDeleteProduct) dto.CommandDeleteProduct {
	value := dto.DeleteProduct{
		ID:        msg.Value.Id,
		CreatedBy: msg.Value.CreatedBy,
	}
	return dto.NewCommandDeleteProduct(msg.Key, value)
}

func mapQueryFindProductToDto(msg *auction.QueryFindProduct) dto.QueryFindProduct {
	value := dto.QueryProduct{
		ID: msg.Value.Id,
	}
	return dto.NewQueryFindProduct(msg.Key, value)
}

func mapQueryFindProductsToDto(msg *auction.QueryFindProducts) dto.QueryFindProducts {
	value := dto.QueryProducts{}
	return dto.NewQueryFindProducts(msg.Key, value)
}
