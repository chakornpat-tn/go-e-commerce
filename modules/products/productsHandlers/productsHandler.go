package productsHandlers

import (
	"strings"

	"github.com/chakornpat-tn/go-rest-api/config"
	"github.com/chakornpat-tn/go-rest-api/modules/appinfo"
	"github.com/chakornpat-tn/go-rest-api/modules/entities"
	"github.com/chakornpat-tn/go-rest-api/modules/files/filesUsecases"
	"github.com/chakornpat-tn/go-rest-api/modules/products"
	"github.com/chakornpat-tn/go-rest-api/modules/products/productsUsecases"
	"github.com/gofiber/fiber/v2"
)

type productsHandlerErrCode string

const (
	FindProductByIdErr productsHandlerErrCode = "products-001"
	FindProductsErr    productsHandlerErrCode = "products-002"
	InsertProductsErr  productsHandlerErrCode = "products-003"
)

type IProductsHandler interface {
	FindProductById(c *fiber.Ctx) error
	FindProducts(c *fiber.Ctx) error
	AddProduct(c *fiber.Ctx) error
}

type productsHandler struct {
	cfg             config.IConfig
	productsUsecase productsUsecases.IProductsUsecase
	filesUsecase    filesUsecases.IFilesUsecase
}

func NewProductsHandler(cfg config.IConfig, productsUsecase productsUsecases.IProductsUsecase, filesUsecase filesUsecases.IFilesUsecase) IProductsHandler {
	return &productsHandler{
		cfg:             cfg,
		productsUsecase: productsUsecase,
		filesUsecase:    filesUsecase,
	}
}

func (h *productsHandler) FindProductById(c *fiber.Ctx) error {
	productId := strings.Trim(c.Params("product_id"), " ")

	product, err := h.productsUsecase.FindProductById(productId)
	if err != nil {
		return entities.NewResponse(c).Error(fiber.ErrInternalServerError.Code, string(FindProductByIdErr), err.Error()).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, product).Res()
}

func (h *productsHandler) FindProducts(c *fiber.Ctx) error {
	req := new(products.ProductFilter)
	req.PaginationReq = new(entities.PaginationReq)
	req.SortReq = new(entities.SortReq)

	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code, string(FindProductsErr), err.Error()).Res()
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.Limit < 5 {
		req.Limit = 5
	}

	if req.OrderBy == "" {
		req.OrderBy = "title"
	}

	if req.Sort == "" {
		req.Sort = "ASC"
	}

	res := h.productsUsecase.FindProducts(req)

	return entities.NewResponse(c).Success(fiber.StatusOK, res).Res()
}

func (h *productsHandler) AddProduct(c *fiber.Ctx) error {
	req := &products.Product{
		Category: new(appinfo.Category),
		Images:   make([]*entities.Image, 0),
	}

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code, string(InsertProductsErr), err.Error()).Res()
	}

	if req.Category.Id <= 0 {
		return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code, string(InsertProductsErr), "category id is invalid").Res()
	}

	product, err := h.productsUsecase.AddProduct(req)
	if err != nil {
		return entities.NewResponse(c).Error(fiber.ErrInternalServerError.Code, string(InsertProductsErr), err.Error()).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, product).Res()
}
