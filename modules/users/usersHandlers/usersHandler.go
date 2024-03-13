package usersHandlers

import (
	"strings"

	"github.com/TurterDev/basicgolang/config"
	"github.com/TurterDev/basicgolang/modules/entities"
	"github.com/TurterDev/basicgolang/modules/users"
	"github.com/TurterDev/basicgolang/modules/users/usersUsecases"
	"github.com/TurterDev/basicgolang/pkg/basicgolangauth"
	"github.com/gofiber/fiber/v2"
)

type userHandlersErrCode string

const (
	SignUpCustomerErr     userHandlersErrCode = "users-001"
	SignInErr             userHandlersErrCode = "users-002"
	RefreshPassportErr    userHandlersErrCode = "users-003"
	SignOutErr            userHandlersErrCode = "users-004"
	SignUpAdminErr        userHandlersErrCode = "users-005"
	GenerateAdminTokenErr userHandlersErrCode = "users-006"
	GetUserProfileErr userHandlersErrCode = "users-007"
)

type IUsersHandler interface {
	SignUpCustomer(c *fiber.Ctx) error
	SignIn(c *fiber.Ctx) error
	RefreshPassport(c *fiber.Ctx) error
	SignOut(c *fiber.Ctx) error
	SignUpAdmin(c *fiber.Ctx) error
	GenerateAdminToken(c *fiber.Ctx) error
	GetUserProfile(c *fiber.Ctx) error
}

type usersHandler struct {
	cfg         config.IConfig
	userUsecase usersUsecases.IUsersUsecase
}

func UsersHandler(cfg config.IConfig, userUsecase usersUsecases.IUsersUsecase) IUsersHandler {
	return &usersHandler{
		cfg:         cfg,
		userUsecase: userUsecase,
	}
}

func (h *usersHandler) SignUpCustomer(c *fiber.Ctx) error {
	//Requst Body parser
	req := new(users.UserRegisterReq)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code, //400
			string(SignUpCustomerErr),
			err.Error(),
		).Res()
	}

	//email validation
	if !req.IsEmail() {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(SignUpCustomerErr),
			"email pattern is invalid",
		).Res()
	}

	// Insert
	result, err := h.userUsecase.InsertCustomer(req)
	if err != nil {
		switch err.Error() {
		case "username has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(SignUpCustomerErr),
				err.Error(),
			).Res()
		case "email has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(SignUpCustomerErr),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(SignUpCustomerErr),
				err.Error(),
			).Res()
		}
	}
	return entities.NewResponse(c).Success(fiber.StatusCreated, result).Res()
}

func (h *usersHandler) SignUpAdmin(c *fiber.Ctx) error {
	//Requst Body parser
	req := new(users.UserRegisterReq)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code, //400
			string(SignUpCustomerErr),
			err.Error(),
		).Res()
	}

	//email validation
	if !req.IsEmail() {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(SignUpCustomerErr),
			"email pattern is invalid",
		).Res()
	}
	// Insert
	result, err := h.userUsecase.InsertCustomer(req)
	if err != nil {
		switch err.Error() {
		case "username has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(SignUpCustomerErr),
				err.Error(),
			).Res()
		case "email has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(SignUpCustomerErr),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(SignUpCustomerErr),
				err.Error(),
			).Res()
		}
	}
	return entities.NewResponse(c).Success(fiber.StatusCreated, result).Res()

}

func (h *usersHandler) GenerateAdminToken(c *fiber.Ctx) error {
	adminToken, err := basicgolangauth.NewBasicgolangAuth(basicgolangauth.Admin, h.cfg.Jwt(), nil)

	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(GenerateAdminTokenErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		&struct {
			Token string `json:"token"`
		}{
			Token: adminToken.SignToken(),
		},
	).Res()
}

func (h *usersHandler) SignIn(c *fiber.Ctx) error {
	req := new(users.UserCredentail)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code, //400
			string(SignInErr),
			err.Error(),
		).Res()
	}

	passport, err := h.userUsecase.GetPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code, //400
			string(SignInErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}

func (h *usersHandler) RefreshPassport(c *fiber.Ctx) error {
	req := new(users.UserRefreshCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code, //400
			string(RefreshPassportErr),
			err.Error(),
		).Res()
	}

	passport, err := h.userUsecase.RefreshPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code, //400
			string(RefreshPassportErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}

func (h *usersHandler) SignOut(c *fiber.Ctx) error {

	//รับ bidy
	req := new(users.UserRemoveCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(SignOutErr),
			err.Error(),
		).Res()
	}

	if err := h.userUsecase.DeleteOauth(req.OauthId); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(SignOutErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, nil).Res()
}


func (h *usersHandler) GetUserProfile(c *fiber.Ctx) error {
	//Set params
	userId := strings.Trim(c.Params("user_id"), " ")

	// Get profile
	result, err := h.userUsecase.GetUserProfile(userId)
	if err != nil {
		switch err.Error() {
		case "get user failed: sql: no rows in result set":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(GetUserProfileErr),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(GetUserProfileErr),
				err.Error(),
			).Res()
		}
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, result).Res()
}