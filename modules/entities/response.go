package entities

import (
	"github.com/chakornpat-tn/go-rest-api/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type IResponse interface {
	Success(code int, data any) IResponse
	Error(code int, tractId, msg string) IResponse
	Res() error
}

type Response struct {
	StatusCode int
	Data       any
	ErrorRes   *ErrorResponse
	Context    *fiber.Ctx
	IsError    bool
}

type ErrorResponse struct {
	TraceId string `json:"traceId"`
	Msg     string `json:"msg"`
}

func NewResponse(c *fiber.Ctx) IResponse {
	return &Response{
		Context: c,
	}
}

func (r *Response) Success(code int, data any) IResponse {
	r.StatusCode = code
	r.Data = data
	logger.InitLogger(r.Context, &r.Data).Print().Save()
	return r
}
func (r *Response) Error(code int, tractId, msg string) IResponse {
	r.StatusCode = code
	r.IsError = true
	r.ErrorRes = &ErrorResponse{
		TraceId: tractId,
		Msg:     msg,
	}
	logger.InitLogger(r.Context, &r.ErrorRes).Print().Save()
	return r
}
func (r *Response) Res() error {
	if r.IsError {
		return r.Context.Status(r.StatusCode).JSON(r.ErrorRes)
	}
	return r.Context.Status(fiber.StatusOK).JSON(r.Data)
}
