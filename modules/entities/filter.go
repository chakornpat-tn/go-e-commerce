package entities

type PaginationReq struct {
	Page      int `query:"page"`
	Limit     int `query:"limit"`
	TotalPage int `query:"total_page"`
	TotalData int `query:"total_data"`
}

type SortReq struct {
	OrderBy string `quey:"order_by"`
	Sort    string `query:"sort"` // DESC ASC
}
