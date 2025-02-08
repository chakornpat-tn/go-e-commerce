package productsUsecases

import (
	"github.com/chakornpat-tn/go-rest-api/modules/products"
	"github.com/chakornpat-tn/go-rest-api/modules/products/productsRepositories"
)

type IProductsUsecase interface {
	FindProductById(productId string) (*products.Product, error)
}

type productUsecase struct {
	productsRepositories productsRepositories.IProductsRepository
}

func NewProductsUsecase(productsRepositories productsRepositories.IProductsRepository) IProductsUsecase {
	return &productUsecase{
		productsRepositories: productsRepositories,
	}
}

func (u *productUsecase) FindProductById(productId string) (*products.Product, error) {
	product, err := u.productsRepositories.FindProductById(productId)
	if err != nil {
		return nil, err
	}
	return product, nil
}
