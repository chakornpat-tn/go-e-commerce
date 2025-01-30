package servers

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/chakornpat-tn/go-rest-api/config"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type IServer interface {
	Start()
}

type server struct {
	app *fiber.App
	db  *sqlx.DB
	cfg config.IConfig
}

func NewServer(cfg config.IConfig, db *sqlx.DB) IServer {
	return &server{
		app: fiber.New(
			fiber.Config{
				AppName:      cfg.APP().Name(),
				BodyLimit:    cfg.APP().BodyLimit(),
				ReadTimeout:  cfg.APP().ReadTimeOut(),
				WriteTimeout: cfg.APP().WriteTimeOut(),
				JSONEncoder:  json.Marshal,
				JSONDecoder:  json.Unmarshal,
			},
		),
		db:  db,
		cfg: cfg,
	}
}

func (s *server) Start() {
	//Middleware
	middlewares := InitMiddlewares(s)
	s.app.Use(middlewares.Logger())
	s.app.Use(middlewares.Cors())

	//Modules
	v1 := s.app.Group("v1")
	modules := InitModule(v1, s, middlewares)

	modules.MonitorModule()
	modules.UsersModule()
	s.app.Use(middlewares.RouterCheck())

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		log.Println("server is shutting down...")
		_ = s.app.Shutdown()
	}()

	// Start server
	log.Println("server is running on ", s.cfg.APP().Url())
	s.app.Listen(s.cfg.APP().Url())
}
