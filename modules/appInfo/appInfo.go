package appinfo

type CategoryFilter struct {
	Title string `json:"title" query:"title"`
}

type Category struct {
	Id    int    `json:"id" db:"id"`
	Title string `json:"title" db:"title"`
}
