package users

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       string `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	Email    string `db:"email" json:"email"`
	RoleId   int    `db:"role_id" json:"role_id"`
}

type UserRegisterReq struct {
	Username string `db:"username" json:"username" form:"username"`
	Password string `db:"password" json:"password" form:"password"`
	Email    string `db:"email" json:"email" form:"email"`
}

type UserCredential struct {
	Email    string `db:"email" json:"email" form:"email"`
	Password string `db:"password" json:"password" form:"password"`
}

type UserCredentialCheck struct {
	Id       string `db:"id" json:"id"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
	Username string `db:"username" json:"username"`
	RoleId   int    `db:"role_id" json:"role_id"`
}

type UserClaims struct {
	Id     string `db:"id" json:"id"`
	RoleId int    `db:"role_id" json:"role"`
}

func (obj *UserRegisterReq) BcryptHashing() error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(obj.Password), 10)
	if err != nil {
		return fmt.Errorf("bcrypt hashing failed: %v", err)
	}
	obj.Password = string(hashPassword)
	return nil
}

func (obj *UserRegisterReq) CheckEmailPattern() bool {
	pattern := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
	match, err := regexp.MatchString(pattern, obj.Email)
	if err != nil {
		return false
	}
	return match
}

type UsersPassport struct {
	User  *User      `json:"user"`
	Token *UserToken `json:"token"`
}

type UserToken struct {
	Id           string `db:"id" json:"id"`
	AccessToken  string `db:"access_token" json:"access_token"`
	RefreshToken string `db:"refresh_token" json:"refresh_token"`
}

type UserRefreshCredential struct {
	RefreshToken string `db:"refresh_token" json:"refresh_token" form:"refresh_token"`
}

type OAuth struct {
	Id     string `db:"id" json:"id"`
	UserId string `db:"user_id" json:"user_id"`
}

type UserRemoveCredential struct {
	OAuthId string `db:"oauth_id" json:"oauth_id" form:"oauth_id"`
}
