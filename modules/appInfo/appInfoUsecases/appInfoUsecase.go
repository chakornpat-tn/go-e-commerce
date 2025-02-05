package appinfoUsecases

import (
	"github.com/chakornpat-tn/go-rest-api/config"
	"github.com/chakornpat-tn/go-rest-api/modules/appinfo/appinfoRepositories"
)

type IAppinfoUsecases interface {
}

type appinfoUsecase struct {
	cfg                 config.IConfig
	appinfoRepositories appinfoRepositories.IAppinfoRepositories
}

func NewappinfoUsecase(cfg config.IConfig, appinfoRepositories appinfoRepositories.IAppinfoRepositories) IAppinfoUsecases {
	return &appinfoUsecase{
		cfg:                 cfg,
		appinfoRepositories: appinfoRepositories,
	}
}
