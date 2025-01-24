package servers

import (
	"github.com/chakornpat-tn/go-rest-api/modules/middlewares/middlewaresHandlers"
	"github.com/chakornpat-tn/go-rest-api/modules/middlewares/middlewaresRepositories"
	"github.com/chakornpat-tn/go-rest-api/modules/middlewares/middlewaresUsecases"
	monitorHandlers "github.com/chakornpat-tn/go-rest-api/modules/monitor/monitorHandler"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
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
