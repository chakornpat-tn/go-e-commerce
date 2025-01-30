package middlewaresUsecases

import (
	"github.com/chakornpat-tn/go-rest-api/modules/middlewares"
	"github.com/chakornpat-tn/go-rest-api/modules/middlewares/middlewaresRepositories"
)

type IMiddlewaresUsecase interface {
	FindAccessToken(userId, accessToken string) bool
	FindRoles() ([]*middlewares.Role, error)
}

type middlewaresUsecase struct {
	middlewaresRepository middlewaresRepositories.IMiddlewaresRepository
}

func MiddlewareUsecase(middlewaresRepository middlewaresRepositories.IMiddlewaresRepository) IMiddlewaresUsecase {
	return &middlewaresUsecase{
		middlewaresRepository: middlewaresRepository,
	}
}

func (u *middlewaresUsecase) FindAccessToken(userId, accessToken string) bool {
	return u.middlewaresRepository.FindAccessToken(userId, accessToken)
}

func (u *middlewaresUsecase) FindRoles() ([]*middlewares.Role, error) {
	roles, err := u.middlewaresRepository.FindRole()
	if err != nil {
		return nil, err
	}

	return roles, nil
}
