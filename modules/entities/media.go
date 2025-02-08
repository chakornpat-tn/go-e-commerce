package entities

type Image struct {
	id       string `db:"id" json:"id"`
	FileName string `db:"filename" json:"filename"`
	Url      string `db:"url" json:"url"`
}
