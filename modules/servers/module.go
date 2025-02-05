package servers

import (
	"github.com/chakornpat-tn/go-rest-api/modules/middlewares/middlewaresHandlers"
	"github.com/chakornpat-tn/go-rest-api/modules/middlewares/middlewaresRepositories"
	"github.com/chakornpat-tn/go-rest-api/modules/middlewares/middlewaresUsecases"
	monitorHandlers "github.com/chakornpat-tn/go-rest-api/modules/monitor/monitorHandler"

	"github.com/chakornpat-tn/go-rest-api/modules/users/usersHandlers"
	"github.com/chakornpat-tn/go-rest-api/modules/users/usersRepositories"
	"github.com/chakornpat-tn/go-rest-api/modules/users/usersUsecases"

	"github.com/chakornpat-tn/go-rest-api/modules/appInfo/appInfoHandlers"
	"github.com/chakornpat-tn/go-rest-api/modules/appInfo/appInfoRepositories"
	"github.com/chakornpat-tn/go-rest-api/modules/appInfo/appInfoUsecases"

	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
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

	router.Post("/signup", handler.SignUpCustomer)
	router.Post("/signin", handler.SignIn)
	router.Post("/signout", handler.SignOut)
	router.Post("/refresh", handler.RefreshPassport)

	router.Get("/admin/secret", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateAdminToken)
	router.Get("/:user_Id", m.mid.JwtAuth(), m.mid.ParamsCheck(), handler.GetUserProfile)
}

func (m *moduleFactory) AppInfoModule() {
	repository := appInfoRepositories.NewAppInfoRepository(m.server.db)
	useCase := appInfoUsecases.NewAppInfoUsecase(m.server.cfg, repository)
	handler := appInfoHandlers.NewAppInfoHandler(m.server.cfg, useCase)

	router := m.router.Group("/appInfo")

	_ = router
	_ = handler

}
