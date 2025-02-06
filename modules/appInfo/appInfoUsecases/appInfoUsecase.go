package appinfoUsecases

import (
	"github.com/chakornpat-tn/go-rest-api/config"
	"github.com/chakornpat-tn/go-rest-api/modules/appinfo"
	"github.com/chakornpat-tn/go-rest-api/modules/appinfo/appinfoRepositories"
)

type IAppinfoUsecases interface {
	FindCategory(req *appinfo.CategoryFilter) ([]*appinfo.Category, error)
	InsertCategory(req []*appinfo.Category) error
	DeleteCategory(categoryId int) error
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

func (u *appinfoUsecase) FindCategory(req *appinfo.CategoryFilter) ([]*appinfo.Category, error) {
	return u.appinfoRepositories.FindCategory(req)
}

func (u *appinfoUsecase) InsertCategory(req []*appinfo.Category) error {
	if err := u.appinfoRepositories.InsertCategory(req); err != nil {
		return err
	}

	return nil
}

func (u *appinfoUsecase) DeleteCategory(categoryId int) error {
	if err := u.appinfoRepositories.DeleteCategory(categoryId); err != nil {
		return err
	}
	return nil

}
