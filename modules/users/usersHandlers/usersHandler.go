package usersHandlers

import (
	"strings"

	"github.com/chakornpat-tn/go-rest-api/config"
	"github.com/chakornpat-tn/go-rest-api/modules/entities"
	"github.com/chakornpat-tn/go-rest-api/modules/users"
	"github.com/chakornpat-tn/go-rest-api/modules/users/usersUsecases"
	"github.com/chakornpat-tn/go-rest-api/pkg/auth"
	"github.com/gofiber/fiber/v2"
)

type userHandlersErrCode string

const (
	signUpCustomerErr  userHandlersErrCode = "users-001"
	signInErr          userHandlersErrCode = "users-002"
	refreshErr         userHandlersErrCode = "users-003"
	signOutErr         userHandlersErrCode = "users-004"
	signUpAdminErr     userHandlersErrCode = "users-005"
	generateAdminToken userHandlersErrCode = "users-006"
	getUserProfileErr  userHandlersErrCode = "users-007"
)

type IUserHandler interface {
	SignUpCustomer(c *fiber.Ctx) error
	SignUpAdmin(c *fiber.Ctx) error
	SignIn(c *fiber.Ctx) error
	SignOut(c *fiber.Ctx) error
	RefreshPassport(c *fiber.Ctx) error
	GenerateAdminToken(c *fiber.Ctx) error
	GetUserProfile(c *fiber.Ctx) error
}

type usersHandler struct {
	cfg         config.IConfig
	userUsecase usersUsecases.IUserUsecase
}

func NewUsersHandler(cfg config.IConfig, userUsecase usersUsecases.IUserUsecase) IUserHandler {
	return &usersHandler{
		cfg:         cfg,
		userUsecase: userUsecase,
	}
}

func (h *usersHandler) SignUpCustomer(c *fiber.Ctx) error {
	req := new(users.UserRegisterReq)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code, string(signUpCustomerErr), err.Error()).Res()
	}

	if !req.CheckEmailPattern() {
		return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code, string(signUpCustomerErr), "Invalid email pattern").Res()
	}

	res, err := h.userUsecase.InsertCustomer(req)
	if err != nil {
		switch err.Error() {
		case "username has been used":
			return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code, string(signUpCustomerErr), err.Error()).Res()
		case "email has been used":
			return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code, string(signUpCustomerErr), err.Error()).Res()
		default:
			return entities.NewResponse(c).Error(fiber.StatusInternalServerError, string(signUpCustomerErr), err.Error()).Res()

		}
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, res).Res()
}

func (h *usersHandler) SignUpAdmin(c *fiber.Ctx) error {
	req := new(users.UserRegisterReq)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code, string(signUpAdminErr), err.Error()).Res()
	}

	if !req.CheckEmailPattern() {
		return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code, string(signUpAdminErr), "Invalid email pattern").Res()
	}

	res, err := h.userUsecase.InsertAdmin(req)
	if err != nil {
		switch err.Error() {
		case "username has been used":
			return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code, string(signUpAdminErr), err.Error()).Res()
		case "email has been used":
			return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code, string(signUpAdminErr), err.Error()).Res()
		default:
			return entities.NewResponse(c).Error(fiber.StatusInternalServerError, string(signUpAdminErr), err.Error()).Res()

		}
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, res).Res()
}

func (h *usersHandler) SignIn(c *fiber.Ctx) error {
	req := new(users.UserCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code, string(signInErr), err.Error()).Res()
	}

	passport, err := h.userUsecase.GetPassPort(req)
	if err != nil {
		return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code, string(signInErr), err.Error()).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}

func (h *usersHandler) RefreshPassport(c *fiber.Ctx) error {
	req := new(users.UserRefreshCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code, string(refreshErr), err.Error()).Res()
	}

	passport, err := h.userUsecase.RefreshToken(req)
	if err != nil {
		return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code, string(refreshErr), err.Error()).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}

func (h *usersHandler) SignOut(c *fiber.Ctx) error {
	req := new(users.UserRemoveCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code, string(signOutErr), err.Error()).Res()
	}

	if err := h.userUsecase.DeleteOAuth(req.OAuthId); err != nil {

		return entities.NewResponse(c).Error(fiber.StatusInternalServerError, string(signOutErr), err.Error()).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, nil).Res()
}

func (h *usersHandler) GenerateAdminToken(c *fiber.Ctx) error {
	adminToken, err := auth.NewAuth(auth.Admin, h.cfg.JWT(), nil)
	if err != nil {
		return entities.NewResponse(c).Error(fiber.StatusInternalServerError, string(generateAdminToken), err.Error()).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, &struct {
		Token string `json:"token"`
	}{
		Token: adminToken.SignToken(),
	}).Res()
}

func (h *usersHandler) GetUserProfile(c *fiber.Ctx) error {
	userId := strings.Trim(c.Params("user_id"), " ")

	user, err := h.userUsecase.GetUserProfile(userId)
	if err != nil {
		switch err.Error() {
		case "get user failed: sql: no rows in result set":
			return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code, string(getUserProfileErr), err.Error()).Res()
		default:
			return entities.NewResponse(c).Error(fiber.StatusInternalServerError, string(getUserProfileErr), err.Error()).Res()
		}
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, user).Res()
}
