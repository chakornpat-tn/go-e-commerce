package productsHandlers

import (
	"github.com/chakornpat-tn/go-rest-api/config"
	"github.com/chakornpat-tn/go-rest-api/modules/files/filesUsecases"
	"github.com/chakornpat-tn/go-rest-api/modules/products/productsUsecases"
)

type IProductsHandler interface {
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
