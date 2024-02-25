package monitorHandlers

import (
	"github.com/TurterDev/basicgolang/config"
	"github.com/TurterDev/basicgolang/modules/entities"
	"github.com/TurterDev/basicgolang/modules/monitor"
	"github.com/gofiber/fiber/v2"
)

type IMonitorHandler interface {
	HealthCheck(c *fiber.Ctx) error
}

type monitorHandler struct {
	cfg config.IConfig
}

func MonitorHandler(cfg config.IConfig) IMonitorHandler {
	return &monitorHandler{
		cfg: cfg,
	}
}

func (h *monitorHandler) HealthCheck(c *fiber.Ctx) error {
	res := &monitor.Monitor{
		Name: h.cfg.App().Name(),
		Version: h.cfg.App().Version(),
	}
	// return c.Status(fiber.StatusOK).JSON(res)
	return entities.NewResponse(c).Success(fiber.StatusOK,res).Res()
}