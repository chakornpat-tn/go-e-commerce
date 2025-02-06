package appinfoHandlers

import (
	"strconv"
	"strings"

	"github.com/chakornpat-tn/go-rest-api/config"
	"github.com/chakornpat-tn/go-rest-api/modules/appinfo"
	"github.com/chakornpat-tn/go-rest-api/modules/appinfo/appinfoUsecases"
	"github.com/chakornpat-tn/go-rest-api/modules/entities"
	"github.com/chakornpat-tn/go-rest-api/pkg/auth"
	"github.com/gofiber/fiber/v2"
)

type appinfoHandlersErrCode string

const (
	generateApiKeyErr appinfoHandlersErrCode = "apiinfo-001"
	findCategoryErr   appinfoHandlersErrCode = "apiinfo-002"
	addCategoryErr    appinfoHandlersErrCode = "apiinfo-003"
	removeCategoryErr appinfoHandlersErrCode = "apiinfo-004"
)

type IAppinfoHandler interface {
	GenerateApiKey(c *fiber.Ctx) error
	FindCategory(c *fiber.Ctx) error
	AddCategory(c *fiber.Ctx) error
	RemoveCategory(c *fiber.Ctx) error
}

type appinfoHandler struct {
	cfg            config.IConfig
	appinfoUsecase appinfoUsecases.IAppinfoUsecases
}

func NewappinfoHandler(cfg config.IConfig, appinfoUsecase appinfoUsecases.IAppinfoUsecases) IAppinfoHandler {
	return &appinfoHandler{
		cfg:            cfg,
		appinfoUsecase: appinfoUsecase,
	}
}

func (h *appinfoHandler) GenerateApiKey(c *fiber.Ctx) error {
	apiKey, err := auth.NewAuth(
		auth.ApiKey,
		h.cfg.JWT(),
		nil,
	)

	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(generateApiKeyErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, &struct {
		Key string `json:"key"`
	}{
		Key: apiKey.SignToken(),
	}).Res()

}

func (h *appinfoHandler) FindCategory(c *fiber.Ctx) error {
	req := new(appinfo.CategoryFilter)
	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findCategoryErr),
			err.Error(),
		).Res()
	}

	categories, err := h.appinfoUsecase.FindCategory(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findCategoryErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, categories).Res()

}

func (h *appinfoHandler) AddCategory(c *fiber.Ctx) error {
	req := make([]*appinfo.Category, 0)
	if err := c.BodyParser(&req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(addCategoryErr),
			err.Error(),
		).Res()
	}
	if len(req) == 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(addCategoryErr),
			"empty request",
		).Res()
	}

	if err := h.appinfoUsecase.InsertCategory(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(addCategoryErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, req).Res()
}

func (h *appinfoHandler) RemoveCategory(c *fiber.Ctx) error {
	categoryId := strings.Trim(c.Params("category_id"), " ")
	categoryIdInt, err := strconv.Atoi(categoryId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(removeCategoryErr),
			"category id type is invalid",
		).Res()
	}
	if categoryIdInt <= 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(removeCategoryErr),
			"invalid category id",
		).Res()
	}

	if err := h.appinfoUsecase.DeleteCategory(categoryIdInt); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(removeCategoryErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, &struct {
		CategoryId int `json:"category_id"`
	}{
		CategoryId: categoryIdInt,
	}).Res()
}
