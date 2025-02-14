package servers

import (
	"github.com/chakornpat-tn/go-rest-api/modules/products/productsHandlers"
	"github.com/chakornpat-tn/go-rest-api/modules/products/productsRepositories"
	"github.com/chakornpat-tn/go-rest-api/modules/products/productsUsecases"
)

type IProductModule interface {
	Init()
	Repository() productsRepositories.IProductsRepository
	Usecase() productsUsecases.IProductsUsecase
	Handler() productsHandlers.IProductsHandler
}

type productsModule struct {
	*moduleFactory
	repository productsRepositories.IProductsRepository
	usecase    productsUsecases.IProductsUsecase
	handler    productsHandlers.IProductsHandler
}

func (m *moduleFactory) ProductsModule() IProductModule {
	productRepository := productsRepositories.NewProductsRepository(m.server.cfg, m.server.db, m.FilesModule().Usecase())
	productsUsecase := productsUsecases.NewProductsUsecase(productRepository)
	productsHandler := productsHandlers.NewProductsHandler(m.server.cfg, productsUsecase, m.FilesModule().Usecase())

	return &productsModule{
		moduleFactory: m,
		repository:    productRepository,
		usecase:       productsUsecase,
		handler:       productsHandler,
	}
}

func (p *productsModule) Init() {
	router := p.router.Group("/products")

	router.Get("/", p.mid.ApiKeyAuth(), p.handler.FindProducts)
	router.Get("/:product_id", p.mid.ApiKeyAuth(), p.handler.FindProductById)
	router.Delete("/:product_id", p.mid.JwtAuth(), p.mid.Authorize(2), p.handler.DeleteProduct)
	router.Patch("/:product_id", p.mid.JwtAuth(), p.mid.Authorize(2), p.handler.UpdateProduct)
	router.Post("/", p.mid.JwtAuth(), p.mid.Authorize(2), p.handler.AddProduct)
}

func (p *productsModule) Repository() productsRepositories.IProductsRepository {
	return p.repository
}

func (p *productsModule) Usecase() productsUsecases.IProductsUsecase {
	return p.usecase
}
func (p *productsModule) Handler() productsHandlers.IProductsHandler {
	return p.handler
}
