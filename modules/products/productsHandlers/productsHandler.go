package productsHandlers

import (
	"strings"

	"github.com/chakornpat-tn/go-rest-api/config"
	"github.com/chakornpat-tn/go-rest-api/modules/entities"
	"github.com/chakornpat-tn/go-rest-api/modules/files/filesUsecases"
	"github.com/chakornpat-tn/go-rest-api/modules/products/productsUsecases"
	"github.com/gofiber/fiber/v2"
)

type productsHandlerErrCode string

const (
	FindProductByIdErr productsHandlerErrCode = "products-001"
)

type IProductsHandler interface {
	FindProductById(c *fiber.Ctx) error
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
