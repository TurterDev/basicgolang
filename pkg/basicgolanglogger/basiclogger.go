package basicgolanglogger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/TurterDev/basicgolang/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type IBasiclogger interface {
	Print() IBasiclogger
	Save()
	SetQuery(c *fiber.Ctx)
	//http://localhost:3000/v1/products?name=test
	//http://localhost:3000/v1/products/:id/
	SetBody(c *fiber.Ctx)
	SetResponse(res any)
}

type basiclogger struct {
	Time       string `json:"time"`
	Ip         string `json:"ip"`
	Method     string `json:"method"`
	StatusCode int    `json:"status_code"`
	Path       string `json:"path"`
	Query      any    `json:"query`
	Body       any    `json:"body`
	Response   any    `json:response`
}

func InitBasicgolang(c *fiber.Ctx, res any) IBasiclogger {
	log := &basiclogger{
		Time:       time.Now().Local().Format("2006-01-02 15:04:05"),
		Ip:         c.IP(),
		Method:     c.Method(),
		Path:       c.Path(),
		StatusCode: c.Response().StatusCode(),
	}
	log.SetQuery(c)
	log.SetBody(c)
	log.SetResponse(res)
	return log
}

func (l *basiclogger) Print() IBasiclogger {
	utils.Debug(l)
	return l
}

func (l *basiclogger) Save() {
	data := utils.Output(l)
	filename := fmt.Sprintf("./assets/logs/basicgolang_%v.txt", strings.ReplaceAll(time.Now().Format("2006-01-02"), "-", ""))
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()
	file.WriteString(string(data) + "\n")
}

func (l *basiclogger) SetQuery(c *fiber.Ctx) {
	var body any
	if err := c.QueryParser(&body); err != nil {
		log.Printf("qurey parser error: %v", err)
	}
	l.Query = body
}

func (l *basiclogger) SetBody(c *fiber.Ctx) {
	var body any
	if err := c.BodyParser(&body); err != nil {
		log.Printf("body parser error: %v", err)
	}
	switch l.Path {
	case "v1/users/signup":
		l.Body = "never gonna give you up"
	default:
		l.Body = body
	}
}

func (l *basiclogger) SetResponse(res any) {
	l.Response = res
}
