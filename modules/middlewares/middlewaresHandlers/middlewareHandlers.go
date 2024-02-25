package middlewaresHandlers

import (
	"github.com/TurterDev/basicgolang/config"
	"github.com/TurterDev/basicgolang/modules/entities"
	"github.com/TurterDev/basicgolang/modules/middlewares/middlewaresUsecases"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// สร้างเป็น enum middleware ประกาศ type
type middlewareHandlersErrCode string

// ประกาศ enum
const (
	routerCheckErr middlewareHandlersErrCode = "middileware-001"
)

type IMiddlewaresHandler interface {
	//สร้าง func Cors โดย return เป็น fiber.Handler
	Cors() fiber.Handler
	//ทำ Router Check
	RouterCheck() fiber.Handler
	//ทำ Logger
	Logger() fiber.Handler
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
		Format: "${time} {${ip}} ${status} - ${method} ${path}\n",
		TimeFormat: "02/01/2006",
		TimeZone: "Bangkok/Asia",
	})
}
