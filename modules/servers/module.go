package servers

import (
	"github.com/chakornpat-tn/go-rest-api/modules/middlewares/middlewaresHandlers"
	"github.com/chakornpat-tn/go-rest-api/modules/middlewares/middlewaresRepositories"
	"github.com/chakornpat-tn/go-rest-api/modules/middlewares/middlewaresUsecases"
	monitorHandlers "github.com/chakornpat-tn/go-rest-api/modules/monitor/monitorHandler"

	"github.com/chakornpat-tn/go-rest-api/modules/users/usersHandlers"
	"github.com/chakornpat-tn/go-rest-api/modules/users/usersRepositories"
	"github.com/chakornpat-tn/go-rest-api/modules/users/usersUsecases"

	"github.com/chakornpat-tn/go-rest-api/modules/appinfo/appinfoHandlers"
	"github.com/chakornpat-tn/go-rest-api/modules/appinfo/appinfoRepositories"
	"github.com/chakornpat-tn/go-rest-api/modules/appinfo/appinfoUsecases"

	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
	AppinfoModule()
}

type moduleFactory struct {
	router fiber.Router
	server *server
	mid    middlewaresHandlers.IMiddlewaresHandler
}

func InitModule(r fiber.Router, s *server, mid middlewaresHandlers.IMiddlewaresHandler) IModuleFactory {
	return &moduleFactory{
		router: r,
		server: s,
		mid:    mid,
	}
}

func InitMiddlewares(s *server) middlewaresHandlers.IMiddlewaresHandler {
	repository := middlewaresRepositories.MiddlewareRepository(s.db)
	useCase := middlewaresUsecases.MiddlewareUsecase(repository)
	handler := middlewaresHandlers.MiddlewareHandler(useCase, s.cfg)
	return handler
}
func (m *moduleFactory) MonitorModule() {
	handler := monitorHandlers.MonitorHandler(m.server.cfg)

	m.router.Get("/", handler.HealthCheck)
}

func (m *moduleFactory) UsersModule() {
	repository := usersRepositories.NewUsersRepository(m.server.db)
	useCase := usersUsecases.NewUsersUsecase(repository, m.server.cfg)
	handler := usersHandlers.NewUsersHandler(m.server.cfg, useCase)

	router := m.router.Group("/users")

	router.Post("/signup", m.mid.ApiKeyAuth(), handler.SignUpCustomer)
	router.Post("/signin", m.mid.ApiKeyAuth(), handler.SignIn)
	router.Post("/signout", m.mid.ApiKeyAuth(), handler.SignOut)
	router.Post("/signup-admin", m.mid.JwtAuth(), m.mid.Authorize(2), handler.SignOut)
	router.Post("/refresh", m.mid.ApiKeyAuth(), handler.RefreshPassport)

	router.Get("/admin/secret", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateAdminToken)
	router.Get("/:user_Id", m.mid.JwtAuth(), m.mid.ParamsCheck(), handler.GetUserProfile)
}

func (m *moduleFactory) AppinfoModule() {
	repository := appinfoRepositories.NewappinfoRepository(m.server.db)
	useCase := appinfoUsecases.NewappinfoUsecase(m.server.cfg, repository)
	handler := appinfoHandlers.NewappinfoHandler(m.server.cfg, useCase)

	router := m.router.Group("/appinfo")

	router.Get("/apikey", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateApiKey)

}
