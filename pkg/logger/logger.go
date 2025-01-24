package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/chakornpat-tn/go-rest-api/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type ILogger interface {
	Print() ILogger
	Save()
	SetQuery(c *fiber.Ctx)
	SetBody(c *fiber.Ctx)
	SetResponse(res any)
}

type logger struct {
	Time       string `json:"time"`
	Ip         string `json:"ip"`
	Method     string `json:"method"`
	StatusCode int    `json:"statusCode"`
	Path       string `json:"path"`
	Query      any    `json:"query"`
	Body       any    `json:"body"`
	Response   any    `json:"response"`
}

func InitLogger(c *fiber.Ctx, res any) ILogger {
	log := &logger{
		Time:       time.Now().Format("2006-01-02 15:04:05"),
		Ip:         c.IP(),
		Method:     c.Method(),
		StatusCode: c.Response().StatusCode(),
		Path:       c.Path(),
	}
	log.SetQuery(c)
	log.SetBody(c)
	log.SetResponse(res)
	return log
}

func (l *logger) Print() ILogger {
	utils.Debug(l)
	return l
}

func (l *logger) Save() {
	data := utils.Output(l)
	fileName := fmt.Sprintf("./assets/logs/%v.txt", strings.ReplaceAll(time.Now().Format("2006-01-02"), "-", ""))
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()
	file.WriteString(string(data) + "\n")
}

func (l *logger) SetQuery(c *fiber.Ctx) {
	var query any
	if err := c.QueryParser(&query); err != nil {
		log.Printf("error parsing query: %v", err)
	}
	l.Query = query
}

func (l *logger) SetBody(c *fiber.Ctx) {
	var body any
	if err := c.BodyParser(&body); err != nil {
		log.Printf("error parsing body: %v", err)
	}

	switch l.Path {
	case "v1/users/signup":
		l.Body = "no pase body for signUp"
	default:
		l.Body = body
	}
}

func (l *logger) SetResponse(res any) {
	l.Response = res
}
