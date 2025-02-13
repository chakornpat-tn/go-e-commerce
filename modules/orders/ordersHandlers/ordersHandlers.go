package ordersHandlers

import (
	"strings"
	"time"

	"github.com/chakornpat-tn/go-rest-api/config"
	"github.com/chakornpat-tn/go-rest-api/modules/entities"
	"github.com/chakornpat-tn/go-rest-api/modules/orders"
	"github.com/chakornpat-tn/go-rest-api/modules/orders/ordersUsecases"
	"github.com/gofiber/fiber/v2"
)

type ordersHandlersErrCode string

const (
	findOrderErr  ordersHandlersErrCode = "orders-001"
	findOrdersErr ordersHandlersErrCode = "orders-002"
)

type IOrdersHandler interface {
	FindOrder(c *fiber.Ctx) error
	FindOrders(c *fiber.Ctx) error
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

func (h *ordersHandler) FindOrders(c *fiber.Ctx) error {
	req := &orders.OrderFilter{
		SortReq:       &entities.SortReq{},
		PaginationReq: &entities.PaginationReq{},
	}

	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(fiber.StatusBadRequest, string(findOrdersErr), err.Error()).Res()
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 5 {
		req.Limit = 5
	}

	orderByMap := map[string]string{
		"id":         `"o"."id"`,
		"created_at": `"o"."created_at"`,
	}
	if orderByMap[req.OrderBy] == "" {
		req.OrderBy = orderByMap["id"]
	}

	sortMap := map[string]string{
		"ASC":  "ASC",
		"DESC": "DESC",
	}
	req.Sort = strings.ToUpper(req.Sort)
	if sortMap[req.Sort] == "" {
		req.Sort = sortMap["DESC"]
	}
	if req.StartDate != "" {
		start, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return entities.NewResponse(c).Error(fiber.StatusBadRequest, string(findOrdersErr), err.Error()).Res()
		}
		req.StartDate = start.Format("2006-01-02")
	}

	if req.EndDate != "" {
		end, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return entities.NewResponse(c).Error(fiber.StatusBadRequest, string(findOrdersErr), err.Error()).Res()
		}
		req.EndDate = end.Format("2006-01-02")
	}

	orders := h.ordersUsecase.FindOrders(req)

	return entities.NewResponse(c).Success(fiber.StatusOK, orders).Res()
}
