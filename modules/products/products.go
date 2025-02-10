package products

import (
	"github.com/chakornpat-tn/go-rest-api/modules/appinfo"
	"github.com/chakornpat-tn/go-rest-api/modules/entities"
)

type Product struct {
	Id          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Category    *appinfo.Category `json:"category"`
	Price       float64           `json:"price"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
	Images      []*entities.Image `json:"images"`
}

type ProductFilter struct {
	Id     string `query:"id"`
	Search string `query:"search"`
	*entities.PaginationReq
	*entities.SortReq
}
