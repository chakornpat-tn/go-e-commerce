package appInfoHandlers

import (
	"github.com/chakornpat-tn/go-rest-api/config"
	"github.com/chakornpat-tn/go-rest-api/modules/appInfo/appInfoUsecases"
)

type IAppInfoHandler interface {
}

type appInfoHandler struct {
	cfg            config.IConfig
	appInfoUsecase appInfoUsecases.IAppInfoUsecases
}

func NewAppInfoHandler(cfg config.IConfig, appInfoUsecase appInfoUsecases.IAppInfoUsecases) IAppInfoHandler {
	return &appInfoHandler{
		cfg:            cfg,
		appInfoUsecase: appInfoUsecase,
	}
}
