package ordersHandlers

import (
	"strings"

	"github.com/chakornpat-tn/go-rest-api/config"
	"github.com/chakornpat-tn/go-rest-api/modules/entities"
	"github.com/chakornpat-tn/go-rest-api/modules/orders/ordersUsecases"
	"github.com/gofiber/fiber/v2"
)

type ordersHandlersErrCode string

const (
	findOrderErr ordersHandlersErrCode = "orders-001"
)

type IOrdersHandler interface {
	FindOrder(c *fiber.Ctx) error
}

type ordersHandler struct {
	cfg           config.IConfig
	ordersUsecase ordersUsecases.IOrdersUsecase
}

func NewOrdersHandler(cfg config.IConfig, ordersUsecase ordersUsecases.IOrdersUsecase) IOrdersHandler {
	return &ordersHandler{
		cfg:           cfg,
		ordersUsecase: ordersUsecase,
	}
}

func (h *ordersHandler) FindOrder(c *fiber.Ctx) error {
	orderId := strings.Trim(c.Params("order_id"), " ")

	order, err := h.ordersUsecase.FindOrder(orderId)
	if err != nil {
		return entities.NewResponse(c).Error(fiber.StatusInternalServerError, string(findOrderErr), err.Error()).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, order).Res()
}
