package ordersUsecases

import (
	"github.com/chakornpat-tn/go-rest-api/modules/orders"
	"github.com/chakornpat-tn/go-rest-api/modules/orders/ordersRepositories"
	"github.com/chakornpat-tn/go-rest-api/modules/products/productsRepositories"
)

type IOrdersUsecase interface {
	FindOrder(orderId string) (*orders.Order, error)
}

type ordersUsecase struct {
	ordersRepository  ordersRepositories.IOrdersRepository
	productRepository productsRepositories.IProductsRepository
}

func NewOrdersUsecase(ordersRepository ordersRepositories.IOrdersRepository, productRepository productsRepositories.IProductsRepository) IOrdersUsecase {
	return &ordersUsecase{
		ordersRepository:  ordersRepository,
		productRepository: productRepository,
	}
}

func (u *ordersUsecase) FindOrder(orderId string) (*orders.Order, error) {
	order, err := u.ordersRepository.FindOrder(orderId)
	if err != nil {
		return nil, err
	}
	return order, nil
}
