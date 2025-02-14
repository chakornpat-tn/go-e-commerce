package ordersUsecases

import (
	"fmt"
	"math"

	"github.com/chakornpat-tn/go-rest-api/modules/entities"
	"github.com/chakornpat-tn/go-rest-api/modules/orders"
	"github.com/chakornpat-tn/go-rest-api/modules/orders/ordersRepositories"
	"github.com/chakornpat-tn/go-rest-api/modules/products/productsRepositories"
)

type IOrdersUsecase interface {
	FindOrder(orderId string) (*orders.Order, error)
	FindOrders(req *orders.OrderFilter) *entities.PaginateRes
	InsertOrder(req *orders.Order) (*orders.Order, error)
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

func (u *ordersUsecase) FindOrders(req *orders.OrderFilter) *entities.PaginateRes {

	orders, count := u.ordersRepository.FindOrders(req)

	return &entities.PaginateRes{
		Data:      orders,
		Page:      req.Page,
		Limit:     req.Limit,
		TotalData: count,
		TotalPage: int(math.Ceil(float64(count) / float64(req.Limit))),
	}
}

func (u *ordersUsecase) InsertOrder(req *orders.Order) (*orders.Order, error) {

	for i := range req.Products {
		if req.Products[i].Product == nil {
			return nil, fmt.Errorf("product is nil")
		}

		prod, err := u.productRepository.FindProductById(req.Products[i].Product.Id)
		if err != nil {
			return nil, err
		}

		req.TotalPaid += req.Products[i].Product.Price * float64(req.Products[i].Qty)
		req.Products[i].Product = prod

	}

	orderId, err := u.ordersRepository.InsertOrder(req)
	if err != nil {
		return nil, err
	}

	order, err := u.ordersRepository.FindOrder(orderId)
	if err != nil {
		return nil, err
	}

	return order, nil
}
