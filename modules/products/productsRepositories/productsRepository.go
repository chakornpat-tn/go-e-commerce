package productsRepositories

import (
	"encoding/json"
	"fmt"

	"github.com/chakornpat-tn/go-rest-api/config"
	"github.com/chakornpat-tn/go-rest-api/modules/entities"
	"github.com/chakornpat-tn/go-rest-api/modules/files/filesUsecases"
	"github.com/chakornpat-tn/go-rest-api/modules/products"
	productsPatterns "github.com/chakornpat-tn/go-rest-api/modules/products/productPatterns"
	"github.com/jmoiron/sqlx"
)

type IProductsRepository interface {
	FindProductById(productId string) (*products.Product, error)
	FindProducts(filter *products.ProductFilter) ([]*products.Product, int)
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

func (r *productRepository) FindProductById(productId string) (*products.Product, error) {

	query := `
		SELECT 
		to_jsonb("t")
	FROM (
		SELECT
			"p"."id",
			"p"."title",
			"p"."description",
			"p"."price",
			"p"."created_at",
			"p"."updated_at",
			(
				SELECT 
					to_jsonb("ct")
				FROM (
					SELECT 
						"c"."id",
						"c"."title"
					FROM "categories" "c"
					LEFT JOIN "products_categories" "pc" ON "pc"."category_id" = "c"."id"
					WHERE "pc"."product_id" = "p"."id"
				) as "ct" 
			) as "category",
			(
				SELECT 
					COALESCE(
						array_to_json(
							array_agg("it")
						), '[]'::json 
					)
				FROM (
					SELECT 
						"i"."id",
						"i"."filename",
						"i"."url"
					FROM "images" "i"
					WHERE "i"."product_id" = "p"."id"
				) as "it"
			) as images
		FROM "products" "p"
		WHERE "p"."id" = $1
		LIMIT 1
	) as "t";
	`

	productBytes := make([]byte, 0)
	product := &products.Product{
		Images: make([]*entities.Image, 0),
	}

	if err := r.db.Get(&productBytes, query, productId); err != nil {
		return nil, fmt.Errorf("error get product: %w", err)
	}

	if err := json.Unmarshal(productBytes, &product); err != nil {
		return nil, fmt.Errorf("error unmarshal product: %w", err)
	}

	return product, nil
}

func (r *productRepository) FindProducts(req *products.ProductFilter) ([]*products.Product, int) {
	builder := productsPatterns.FindProductBuilder(r.db, req)
	engineer := productsPatterns.FindProductEngineer(builder)

	result := engineer.FindProduct().Result()
	count := engineer.CountProduct().Count()
	return result, count
}
