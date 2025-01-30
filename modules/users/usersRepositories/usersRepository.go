package usersRepositories

import (
	"context"
	"fmt"
	"time"

	"github.com/chakornpat-tn/go-rest-api/modules/users"
	"github.com/chakornpat-tn/go-rest-api/modules/users/usersPatterns"
	"github.com/jmoiron/sqlx"
)

type IUserRepository interface {
	InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UsersPassport, error)
	FindUserByEmail(email string) (*users.UserCredentialCheck, error)
	InsertOAuth(req *users.UsersPassport) error
	FindOAuth(refreshToken string) (*users.OAuth, error)
	UpdateOAuth(req *users.UserToken) error
	GetProfile(id string) (*users.User, error)
	DeleteOAuth(oauthId string) error
}

type usersRepository struct {
	db *sqlx.DB
}

func NewUsersRepository(db *sqlx.DB) IUserRepository {
	return &usersRepository{db: db}
}

func (r *usersRepository) InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UsersPassport, error) {
	result := usersPatterns.InsertUser(r.db, req, isAdmin)
	var err error
	if isAdmin {
		result, err = result.Admin()
		if err != nil {
			return nil, err
		}
	} else {
		result, err = result.Customer()
		if err != nil {
			return nil, err
		}
	}

	// Get Result from inserting
	user, err := result.Result()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *usersRepository) FindUserByEmail(email string) (*users.UserCredentialCheck, error) {
	query := `SELECT "id", "email", "password", "username", "role_id" FROM users WHERE email = $1`
	user := &users.UserCredentialCheck{}
	err := r.db.Get(user, query, email)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (r *usersRepository) InsertOAuth(req *users.UsersPassport) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	query := `INSERT INTO oauth ("user_id", "access_token", "refresh_token") VALUES ($1, $2, $3) RETURNING "id";`
	if err := r.db.QueryRowContext(ctx, query, req.User.Id, req.Token.AccessToken, req.Token.RefreshToken).Scan(&req.Token.Id); err != nil {
		return err
	}
	return nil
}

func (r *usersRepository) FindOAuth(refreshToken string) (*users.OAuth, error) {
	query := `SELECT 
		"id",
		"user_id"
	FROM "oauth"
	WHERE "refresh_token" = $1;`
	oauth := &users.OAuth{}
	if err := r.db.Get(oauth, query, refreshToken); err != nil {
		return nil, fmt.Errorf("oauth not found")
	}
	return oauth, nil
}

func (r *usersRepository) UpdateOAuth(req *users.UserToken) error {
	query := `
	UPDATE "oauth" SET
		"access_token" = :access_token,
		"refresh_token" = :refresh_token
	WHERE "id" = :id;`

	if _, err := r.db.NamedExecContext(context.Background(), query, req); err != nil {
		return fmt.Errorf("update oauth failed: %v", err)
	}
	return nil
}

func (r *usersRepository) GetProfile(userId string) (*users.User, error) {
	query := `
	SELECT
		"id",
		"email",
		"username",
		"role_id"
	FROM "users"
	WHERE "id" = $1;`

	profile := new(users.User)
	if err := r.db.Get(profile, query, userId); err != nil {
		return nil, fmt.Errorf("get user failed: %v", err)
	}
	return profile, nil
}

func (r *usersRepository) DeleteOAuth(oauthId string) error {
	query := `DELETE FROM "oauth" WHERE "id" = $1;`

	if _, err := r.db.ExecContext(context.Background(), query, oauthId); err != nil {
		return fmt.Errorf("delete oauth failed: %v", err)
	}

	return nil
}
