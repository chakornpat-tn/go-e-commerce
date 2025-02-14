package servers

import (
	"github.com/chakornpat-tn/go-rest-api/modules/files/filesHandlers"
	"github.com/chakornpat-tn/go-rest-api/modules/files/filesUsecases"
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

	"github.com/chakornpat-tn/go-rest-api/modules/products/productsHandlers"
	"github.com/chakornpat-tn/go-rest-api/modules/products/productsRepositories"
	"github.com/chakornpat-tn/go-rest-api/modules/products/productsUsecases"

	"github.com/chakornpat-tn/go-rest-api/modules/orders/ordersHandlers"
	"github.com/chakornpat-tn/go-rest-api/modules/orders/ordersRepositories"
	"github.com/chakornpat-tn/go-rest-api/modules/orders/ordersUsecases"

	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
	AppinfoModule()
	FilesModule()
	ProductsModule()
	OrdersModule()
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
	router.Post("/categories", m.mid.JwtAuth(), m.mid.Authorize(2), handler.AddCategory)
	router.Delete("/:category_id/categories/", m.mid.JwtAuth(), m.mid.Authorize(2), handler.RemoveCategory)
	router.Get("/categories", m.mid.ApiKeyAuth(), handler.FindCategory)
	router.Delete("/categories", m.mid.ApiKeyAuth(), handler.FindCategory)

}

func (m *moduleFactory) FilesModule() {
	handler := filesHandlers.NewFilesHandler(
		m.server.cfg,
		filesUsecases.NewFilesUsecase(m.server.cfg),
	)

	router := m.router.Group("/files")

	router.Post("/upload", m.mid.JwtAuth(), m.mid.Authorize(2), handler.UploadFiles)
	router.Delete("/delete", m.mid.JwtAuth(), m.mid.Authorize(2), handler.DeleteFiles)
}

func (m *moduleFactory) ProductsModule() {
	filesUsecase := filesUsecases.NewFilesUsecase(m.server.cfg)

	repository := productsRepositories.NewProductsRepository(m.server.cfg, m.server.db, filesUsecase)
	useCase := productsUsecases.NewProductsUsecase(repository)
	handler := productsHandlers.NewProductsHandler(m.server.cfg, useCase, filesUsecase)

	router := m.router.Group("/products")

	router.Get("/", m.mid.ApiKeyAuth(), handler.FindProducts)
	router.Get("/:product_id", m.mid.ApiKeyAuth(), handler.FindProductById)
	router.Delete("/:product_id", m.mid.JwtAuth(), m.mid.Authorize(2), handler.DeleteProduct)
	router.Patch("/:product_id", m.mid.JwtAuth(), m.mid.Authorize(2), handler.UpdateProduct)
	router.Post("/", m.mid.JwtAuth(), m.mid.Authorize(2), handler.AddProduct)
}

func (m *moduleFactory) OrdersModule() {

	filesUsecase := filesUsecases.NewFilesUsecase(m.server.cfg)
	productsRepo := productsRepositories.NewProductsRepository(m.server.cfg, m.server.db, filesUsecase)

	repository := ordersRepositories.NewOrdersRepository(m.server.db)
	useCase := ordersUsecases.NewOrdersUsecase(repository, productsRepo)
	handler := ordersHandlers.NewOrdersHandler(m.server.cfg, useCase)

	router := m.router.Group("/orders")

	router.Get("/", m.mid.JwtAuth(), m.mid.Authorize(2), handler.FindOrders)
	router.Get("/:user_id/:order_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), handler.FindOrder)

	router.Patch("/:user_id/:order_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), handler.UpdateOrder)

	router.Post("/", m.mid.JwtAuth(), handler.InsertOrder)

}
