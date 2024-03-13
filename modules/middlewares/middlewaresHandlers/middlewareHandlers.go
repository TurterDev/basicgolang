package middlewaresHandlers

import (
	"strings"

	"github.com/TurterDev/basicgolang/config"
	"github.com/TurterDev/basicgolang/modules/entities"
	"github.com/TurterDev/basicgolang/modules/middlewares/middlewaresUsecases"
	"github.com/TurterDev/basicgolang/pkg/basicgolangauth"
	"github.com/TurterDev/basicgolang/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// สร้างเป็น enum middleware ประกาศ type
type middlewareHandlersErrCode string

// ประกาศ enum
const (
	routerCheckErr middlewareHandlersErrCode = "middileware-001"
	JwtAuthErr     middlewareHandlersErrCode = "middileware-002"
	paramsCheckErr middlewareHandlersErrCode = "middileware-003"
	authorizeErr   middlewareHandlersErrCode = "middileware-004"
)

type IMiddlewaresHandler interface {
	//สร้าง func Cors โดย return เป็น fiber.Handler
	Cors() fiber.Handler
	//ทำ Router Check
	RouterCheck() fiber.Handler
	//ทำ Logger
	Logger() fiber.Handler
	//ทำ JwtAuthen
	JwtAuth() fiber.Handler
	//ทำ ParamsCheck
	ParamsCheck() fiber.Handler
	//ทำ Role Based
	Authorize(expectRoleId ...int) fiber.Handler
}

type middlewaresHandler struct {
	//เริ่มทำ CORS
	cfg config.IConfig

	middlewaresUsecase middlewaresUsecases.IMiddlewaresUsecase
}

func MiddlewaresHandler(cfg config.IConfig, middlewaresUsecase middlewaresUsecases.IMiddlewaresUsecase) IMiddlewaresHandler {
	return &middlewaresHandler{
		cfg:                cfg,
		middlewaresUsecase: middlewaresUsecase,
	}
}

// สร้าง func cors
func (h *middlewaresHandler) Cors() fiber.Handler {
	return cors.New(cors.Config{
		Next:             cors.ConfigDefault.Next,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	})
}

// สร้าง func Routerchek
func (h *middlewaresHandler) RouterCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return entities.NewResponse(c).Error(
			fiber.ErrNotFound.Code,
			string(routerCheckErr),
			"router not found",
		).Res()
	}
}

// สร้าง func Logger
func (h *middlewaresHandler) Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} {${ip}} ${status} - ${method} ${path}\n",
		TimeFormat: "02/01/2006",
		TimeZone:   "Bangkok/Asia",
	})
}

// JwtAuthen
func (h *middlewaresHandler) JwtAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		result, err := basicgolangauth.ParseToken(h.cfg.Jwt(), token)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(JwtAuthErr),
				err.Error(),
			).Res()
		}

		claims := result.Claims
		if !h.middlewaresUsecase.FindAccessToken(claims.Id, token) {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(JwtAuthErr),
				"no permission to access",
			).Res()
		}

		//Set UserId
		c.Locals("userId", claims.Id)
		c.Locals("userRoleId", claims.RoleId)

		//c.Next เป็นการส่งการทำงานให้ func ต่อไป
		return c.Next()
	}
}

// ParamsCheck
func (h *middlewaresHandler) ParamsCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("userId")
		if c.Locals("userRoleId").(int) == 2 {
			return c.Next()
		}
		if c.Params("user_id") != userId {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(paramsCheckErr),
				"never gonna give you up",
			).Res()
		}
		return c.Next()
	}
}

// Role Base
func (h *middlewaresHandler) Authorize(expectRoleId ...int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRoleId, ok := c.Locals("userRoleId").(int)
		if !ok {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(authorizeErr),
				"user_id id not int type",
			).Res()
		}
		roles, err := h.middlewaresUsecase.FindRole()
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(authorizeErr),
				err.Error(),
			).Res()
		}
		sum := 0
		for _, v := range expectRoleId {
			sum += v
		}

		// นำ userRole กับ expectRole มา convert เป็น binary
		expectedValueBinary := utils.BinaryConverter(sum, len(roles))
		userValueBibary := utils.BinaryConverter(userRoleId, len(roles))

		//user ->     0 1 0
		//expected -> 1 1 0

		for i := range userValueBibary {
			if userValueBibary[i]&expectedValueBinary[i] == 1 {
				return c.Next()
			}
		}
		return entities.NewResponse(c).Error(
			fiber.ErrUnauthorized.Code,
			string(authorizeErr),
			"no permission to access",
		).Res()
		
	}
}
