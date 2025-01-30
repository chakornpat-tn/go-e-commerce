package usersUsecases

import (
	"github.com/chakornpat-tn/go-rest-api/config"
	"github.com/chakornpat-tn/go-rest-api/modules/users"
	"github.com/chakornpat-tn/go-rest-api/modules/users/usersRepositories"
	"github.com/chakornpat-tn/go-rest-api/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UsersPassport, error)
	InsertAdmin(req *users.UserRegisterReq) (*users.UsersPassport, error)
	GetPassPort(req *users.UserCredential) (*users.UsersPassport, error)
	RefreshToken(req *users.UserRefreshCredential) (*users.UsersPassport, error)
	DeleteOAuth(oauthId string) error
	GetUserProfile(userId string) (*users.User, error)
}

type usersUsecase struct {
	cfg               config.IConfig
	usersRepositories usersRepositories.IUserRepository
}

func NewUsersUsecase(usersRepositories usersRepositories.IUserRepository, cfg config.IConfig) IUserUsecase {
	return &usersUsecase{
		usersRepositories: usersRepositories,
		cfg:               cfg,
	}
}

func (u *usersUsecase) InsertCustomer(req *users.UserRegisterReq) (*users.UsersPassport, error) {
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}

	result, err := u.usersRepositories.InsertUser(req, false)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *usersUsecase) InsertAdmin(req *users.UserRegisterReq) (*users.UsersPassport, error) {
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}

	result, err := u.usersRepositories.InsertUser(req, true)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *usersUsecase) GetPassPort(req *users.UserCredential) (*users.UsersPassport, error) {
	// Find User
	user, err := u.usersRepositories.FindUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, err
	}

	accessToken, err := auth.NewAuth(auth.AccessToken, u.cfg.JWT(), &users.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	})
	if err != nil {
		return nil, err
	}

	refreshToken, err := auth.NewAuth(auth.RefreshToken, u.cfg.JWT(), &users.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	})
	if err != nil {
		return nil, err
	}

	passport := &users.UsersPassport{
		User: &users.User{
			Id:       user.Id,
			Email:    user.Email,
			Username: user.Username,
			RoleId:   user.RoleId,
		},
		Token: &users.UserToken{
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken.SignToken(),
		},
	}

	if err := u.usersRepositories.InsertOAuth(passport); err != nil {
		return nil, err
	}

	return passport, nil
}

func (u *usersUsecase) RefreshToken(req *users.UserRefreshCredential) (*users.UsersPassport, error) {
	claims, err := auth.ParseToken(u.cfg.JWT(), req.RefreshToken)
	if err != nil {
		return nil, err
	}

	oauth, err := u.usersRepositories.FindOAuth(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	profile, err := u.usersRepositories.GetProfile(oauth.UserId)
	if err != nil {
		return nil, err
	}

	newClaims := &users.UserClaims{
		Id:     profile.Id,
		RoleId: profile.RoleId,
	}

	accessToken, err := auth.NewAuth(auth.AccessToken, u.cfg.JWT(), newClaims)
	if err != nil {
		return nil, err
	}

	refreshToken := auth.RepeatToke(u.cfg.JWT(), newClaims, claims.ExpiresAt.Unix())

	passport := &users.UsersPassport{
		User: profile,
		Token: &users.UserToken{
			Id:           oauth.Id,
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken,
		},
	}

	if err := u.usersRepositories.UpdateOAuth(passport.Token); err != nil {
		return nil, err
	}

	return passport, nil

}

func (u *usersUsecase) DeleteOAuth(oauthId string) error {
	if err := u.usersRepositories.DeleteOAuth(oauthId); err != nil {
		return err
	}
	return nil
}

func (u *usersUsecase) GetUserProfile(userId string) (*users.User, error) {
	user, err := u.usersRepositories.GetProfile(userId)
	if err != nil {
		return nil, err
	}
	return user, err
}
