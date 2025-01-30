package middlewaresRepositories

import (
	"fmt"

	"github.com/chakornpat-tn/go-rest-api/modules/middlewares"
	"github.com/jmoiron/sqlx"
)

type IMiddlewaresRepository interface {
	FindAccessToken(userId, accessToken string) bool
	FindRole() ([]*middlewares.Role, error)
}

type middlewaresRepository struct {
	db *sqlx.DB
}

func MiddlewareRepository(db *sqlx.DB) IMiddlewaresRepository {
	return &middlewaresRepository{
		db: db,
	}
}

func (r *middlewaresRepository) FindAccessToken(userId, accessToken string) bool {

	query := `SELECT (CASE WHEN COUNT(*) > 0 THEN true ELSE false END)
	FROM
		"oauth"
	WHERE
	"user_id" = $1
	AND "access_token" = $2;`

	var foundToken bool
	if err := r.db.Get(&foundToken, query, userId, accessToken); err != nil {
		return false
	}
	return foundToken
}

func (r *middlewaresRepository) FindRole() ([]*middlewares.Role, error) {
	query := `
	SELECT
		"id",
		"title"
	FROM
		"roles"
	ORDER BY "id" DESC;`

	var roles = make([]*middlewares.Role, 0)
	if err := r.db.Select(&roles, query); err != nil {
		return nil, fmt.Errorf("roles are empty")
	}

	return roles, nil
}
