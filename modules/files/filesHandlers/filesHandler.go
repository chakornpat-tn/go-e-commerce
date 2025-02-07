package filesHandlers

import (
	"fmt"
	"math"
	"path/filepath"
	"strings"

	"github.com/chakornpat-tn/go-rest-api/config"
	"github.com/chakornpat-tn/go-rest-api/modules/entities"
	"github.com/chakornpat-tn/go-rest-api/modules/files"
	"github.com/chakornpat-tn/go-rest-api/modules/files/filesUsecases"
	"github.com/chakornpat-tn/go-rest-api/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type filesHandlerErrCode string

const (
	uploadErr     filesHandlerErrCode = "files-001"
	deleteFileErr filesHandlerErrCode = "files-002"
)

type IFilesHandler interface {
	UploadFiles(c *fiber.Ctx) error
	DeleteFiles(c *fiber.Ctx) error
}

type filesHandler struct {
	cfg          config.IConfig
	filesUsecase filesUsecases.IFilesUsecase
}

func NewFilesHandler(cfg config.IConfig, filesUsecase filesUsecases.IFilesUsecase) IFilesHandler {
	return &filesHandler{
		cfg:          cfg,
		filesUsecase: filesUsecase,
	}
}

func (h *filesHandler) UploadFiles(c *fiber.Ctx) error {
	req := make([]*files.FileReq, 0)
	form, err := c.MultipartForm()
	if err != nil {
		return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code, string(uploadErr), err.Error()).Res()
	}

	filesReq := form.File["files"]
	destination := c.FormValue("destination")

	// File ext validation
	extMap := map[string]string{
		"png":  "png",
		"jpg":  "jpb",
		"jpeg": "jpeg",
	}

	for _, file := range filesReq {
		ext := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
		if extMap[ext] != ext || extMap[ext] == "" {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(uploadErr),
				"extension is not acceptable",
			).Res()
		}

		if file.Size > int64(h.cfg.APP().Filelimit()) {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(uploadErr),
				fmt.Sprintf("file size must less than %f MiB", float64(h.cfg.APP().Filelimit())/math.Pow(1024, 2)),
			).Res()
		}

		fileName := utils.RandFileName(ext)
		req = append(req, &files.FileReq{
			File:        file,
			Destination: destination + "/" + fileName,
			FileName:    fileName,
			Extension:   ext,
		})
	}

	res, err := h.filesUsecase.UpLoadToGCP(req)
	if err != nil {
		return entities.NewResponse(c).Error(fiber.ErrInternalServerError.Code, string(uploadErr), err.Error()).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, res).Res()
}

func (h *filesHandler) DeleteFiles(c *fiber.Ctx) error {
	req := make([]*files.DeleteFileReq, 0)
	if err := c.BodyParser(&req); err != nil {
		return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code, string(deleteFileErr), err.Error()).Res()
	}

	if err := h.filesUsecase.DeleteFileOnGCP(req); err != nil {
		return entities.NewResponse(c).Error(fiber.ErrInternalServerError.Code, string(deleteFileErr), err.Error()).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, nil).Res()
}
