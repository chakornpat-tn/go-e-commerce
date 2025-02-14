package middlewaresHandlers

import (
	"strings"

	"github.com/chakornpat-tn/go-rest-api/config"
	"github.com/chakornpat-tn/go-rest-api/modules/entities"
	"github.com/chakornpat-tn/go-rest-api/modules/middlewares/middlewaresUsecases"
	"github.com/chakornpat-tn/go-rest-api/pkg/auth"
	"github.com/chakornpat-tn/go-rest-api/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
)

type middlewareHandlerErrCode string

const (
	routerCheckErr middlewareHandlerErrCode = "router-001"
	jwtAuthErr     middlewareHandlerErrCode = "router-002"
	paramsErr      middlewareHandlerErrCode = "router-003"
	authorizeErr   middlewareHandlerErrCode = "router-004"
	apiKeyErr      middlewareHandlerErrCode = "router-005"
)

type IMiddlewaresHandler interface {
	Cors() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler
	JwtAuth() fiber.Handler
	ParamsCheck() fiber.Handler
	Authorize(expectRoleId ...int) fiber.Handler
	ApiKeyAuth() fiber.Handler
}

type middlewaresHandler struct {
	cfg                config.IConfig
	middlewaresUsecase middlewaresUsecases.IMiddlewaresUsecase
}

func MiddlewareHandler(middlewaresUseCase middlewaresUsecases.IMiddlewaresUsecase, cfg config.IConfig) IMiddlewaresHandler {
	return &middlewaresHandler{
		middlewaresUsecase: middlewaresUseCase,
		cfg:                cfg,
	}
}

func (h *middlewaresHandler) Cors() fiber.Handler {
	return cors.New(cors.Config{
		Next:             cors.ConfigDefault.Next,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,HEAD,PATCH",
		AllowHeaders:     "",
		AllowCredentials: false,
		MaxAge:           0,
	})
}

func (h *middlewaresHandler) RouterCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return entities.NewResponse(c).Error(
			fiber.StatusNotFound,
			string(routerCheckErr),
			"Router not found",
		).Res()

	}
}

func (h *middlewaresHandler) Logger() fiber.Handler {
	return fiberLogger.New(fiberLogger.Config{
		Format:     "${time} [${ip}] ${status} - ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Asia/Bangkok",
	})
}

func (h *middlewaresHandler) JwtAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		result, err := auth.ParseToken(h.cfg.JWT(), token)
		if err != nil {
			return entities.NewResponse(c).Error(fiber.StatusUnauthorized, string(jwtAuthErr), err.Error()).Res()
		}

		claims := result.Claims
		if !h.middlewaresUsecase.FindAccessToken(claims.Id, token) {
			return entities.NewResponse(c).Error(fiber.StatusUnauthorized, string(jwtAuthErr), "no permission to access").Res()
		}

		c.Locals("userId", claims.Id)
		c.Locals("userRoleId", claims.RoleId)

		return c.Next()
	}
}

func (h *middlewaresHandler) ParamsCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("userId")
		if c.Locals("userRoleId").(int) == 2 {
			return c.Next()
		}
		if c.Params("user_id") != userId {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(paramsErr),
				"param error",
			).Res()
		}
		return c.Next()
	}
}

func (h *middlewaresHandler) Authorize(expectRoleId ...int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRoleId, ok := c.Locals("userRoleId").(int)
		if !ok {
			return entities.NewResponse(c).Error(fiber.StatusUnauthorized, string(authorizeErr), "user_id type error").Res()
		}

		roles, err := h.middlewaresUsecase.FindRoles()
		if err != nil {
			return entities.NewResponse(c).Error(fiber.StatusInternalServerError, string(authorizeErr), err.Error()).Res()
		}

		sum := 0
		for _, v := range expectRoleId {
			sum += v
		}

		expectValueBinary := utils.BinaryConverter(sum, len(roles))
		userValueBinary := utils.BinaryConverter(userRoleId, len(roles))

		for i := range userValueBinary {
			if userValueBinary[i]&expectValueBinary[i] == 1 {
				return c.Next()
			}
		}

		return entities.NewResponse(c).Error(
			fiber.ErrUnauthorized.Code,
			string(authorizeErr),
			"no permission to access",
		).Res()
	}
}

func (h *middlewaresHandler) ApiKeyAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := c.Get("X-API-KEY")
		if _, err := auth.ParseApiKey(h.cfg.JWT(), key); err != nil {
			return entities.NewResponse(c).Error(fiber.ErrUnauthorized.Code, string(apiKeyErr), "invalid api key").Res()
		}

		return c.Next()
	}
}
