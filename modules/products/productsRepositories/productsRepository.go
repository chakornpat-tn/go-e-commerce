package productsRepositories

import (
	"github.com/chakornpat-tn/go-rest-api/config"
	"github.com/chakornpat-tn/go-rest-api/modules/files/filesUsecases"
	"github.com/jmoiron/sqlx"
)

type IProductsRepository interface {
}

type productRepository struct {
	cfg          config.IConfig
	db           *sqlx.DB
	filesUsecase filesUsecases.IFilesUsecase
}

func NewProductsRepository(cfg config.IConfig, db *sqlx.DB, filesUsecase filesUsecases.IFilesUsecase) IProductsRepository {
	return &productRepository{
		cfg:          cfg,
		db:           db,
		filesUsecase: filesUsecase,
	}
}
