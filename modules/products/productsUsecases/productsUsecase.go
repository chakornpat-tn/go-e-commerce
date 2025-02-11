package productsUsecases

import (
	"math"

	"github.com/chakornpat-tn/go-rest-api/modules/entities"
	"github.com/chakornpat-tn/go-rest-api/modules/products"
	"github.com/chakornpat-tn/go-rest-api/modules/products/productsRepositories"
)

type IProductsUsecase interface {
	FindProductById(productId string) (*products.Product, error)
	FindProducts(req *products.ProductFilter) *entities.PaginateRes
	AddProduct(req *products.Product) (*products.Product, error)
	UpdateProduct(req *products.Product) (*products.Product, error)
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

func (u *productUsecase) FindProducts(req *products.ProductFilter) *entities.PaginateRes {
	products, count := u.productsRepositories.FindProducts(req)

	return &entities.PaginateRes{
		Data:      products,
		Page:      req.Page,
		Limit:     req.Limit,
		TotalData: count,
		TotalPage: int(math.Ceil(float64(count) / float64(req.Limit))),
	}
}

func (u *productUsecase) AddProduct(req *products.Product) (*products.Product, error) {
	product, err := u.productsRepositories.InsertProduct(req)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (u *productUsecase) UpdateProduct(req *products.Product) (*products.Product, error) {
	product, err := u.productsRepositories.UpdateProduct(req)
	if err != nil {
		return nil, err
	}
	return product, nil
}
