package servers

import (
	"github.com/chakornpat-tn/go-rest-api/modules/files/filesHandlers"
	"github.com/chakornpat-tn/go-rest-api/modules/files/filesUsecases"
)

type IFileModule interface {
	Init()
	Usecase() filesUsecases.IFilesUsecase
	Handler() filesHandlers.IFilesHandler
}

type filesModule struct {
	*moduleFactory
	usecase filesUsecases.IFilesUsecase
	handler filesHandlers.IFilesHandler
}

func (m *moduleFactory) FilesModule() IFileModule {
	usecase := filesUsecases.NewFilesUsecase(m.server.cfg)
	handler := filesHandlers.NewFilesHandler(m.server.cfg, usecase)

	return &filesModule{
		moduleFactory: m,
		usecase:       usecase,
		handler:       handler,
	}
}

func (f *filesModule) Init() {
	router := f.router.Group("/files")
	router.Post("/upload", f.mid.JwtAuth(), f.mid.Authorize(2), f.handler.UploadFiles)
	router.Delete("/delete", f.mid.JwtAuth(), f.mid.Authorize(2), f.handler.DeleteFiles)
}

func (f *filesModule) Usecase() filesUsecases.IFilesUsecase {
	return f.usecase
}
func (f *filesModule) Handler() filesHandlers.IFilesHandler {
	return f.handler
}
