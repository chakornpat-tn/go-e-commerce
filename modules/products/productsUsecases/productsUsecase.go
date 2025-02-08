package productsUsecases

import (
	"github.com/chakornpat-tn/go-rest-api/modules/products/productsRepositories"
)

type IProductsUsecase interface {
}

type productUsecase struct {
	productsRepositories productsRepositories.IProductsRepository
}

func NewProductsUsecase(productsRepositories productsRepositories.IProductsRepository) IProductsUsecase {
	return &productUsecase{
		productsRepositories: productsRepositories,
	}
}
