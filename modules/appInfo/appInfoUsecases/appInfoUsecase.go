package appInfoUsecases

import (
	"github.com/chakornpat-tn/go-rest-api/config"
	"github.com/chakornpat-tn/go-rest-api/modules/appInfo/appInfoRepositories"
)

type IAppInfoUsecases interface {
}

type appInfoUsecase struct {
	cfg                 config.IConfig
	appInfoRepositories appInfoRepositories.IAppInfoRepositories
}

func NewAppInfoUsecase(cfg config.IConfig, appInfoRepositories appInfoRepositories.IAppInfoRepositories) IAppInfoUsecases {
	return &appInfoUsecase{
		cfg:                 cfg,
		appInfoRepositories: appInfoRepositories,
	}
}
