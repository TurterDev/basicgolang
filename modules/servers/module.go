package servers

import (
	"github.com/TurterDev/basicgolang/modules/middlewares/middlewaresHandlers"
	"github.com/TurterDev/basicgolang/modules/middlewares/middlewaresRepositories"
	"github.com/TurterDev/basicgolang/modules/middlewares/middlewaresUsecases"
	"github.com/TurterDev/basicgolang/modules/monitor/monitorHandlers"
	"github.com/TurterDev/basicgolang/modules/users/usersHandlers"
	"github.com/TurterDev/basicgolang/modules/users/usersRepositories"
	"github.com/TurterDev/basicgolang/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
}

type moduleFactory struct {
	r   fiber.Router
	s   *server
	mid middlewaresHandlers.IMiddlewaresHandler
}

func InitModule(r fiber.Router, s *server, mid middlewaresHandlers.IMiddlewaresHandler) IModuleFactory {
	return &moduleFactory{
		r:   r,
		s:   s,
		mid: mid,
	}
}

// func middleware
func InitMiddlewares(s *server) middlewaresHandlers.IMiddlewaresHandler {
	repository := middlewaresRepositories.MiddlewaresRepository(s.db)
	usecase := middlewaresUsecases.MiddlewaresUsecase(repository)
	handler := middlewaresHandlers.MiddlewaresHandler(s.cfg, usecase)
	return handler
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorHandlers.MonitorHandler(m.s.cfg)

	m.r.Get("/", handler.HealthCheck)
}

func (m *moduleFactory) UsersModule() {
	repository := usersRepositories.UsersRepository(m.s.db)
	usecase := usersUsecases.UsersUsecase(m.s.cfg, repository)
	handler := usersHandlers.UsersHandler(m.s.cfg, usecase)

	// /v1/users/sign
	router := m.r.Group("/users")

	router.Post("/signup", handler.SignUpCustomer)
	router.Post("/signin", handler.SignIn)
	//refresh token
	router.Post("/refresh", handler.RefreshPassport)
	//SignOut
	router.Post("/signout", handler.SignOut)

	//GenerateAdminToken
	// router.Get("/secret", handler.GenerateAdminToken)
	// router.Get("/secret", m.mid.JwtAuth(), handler.GenerateAdminToken)
	//Role Based
	router.Get("/admin/secret", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateAdminToken)

	//get user profile
	router.Get("/:user_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), handler.GetUserProfile)
}
