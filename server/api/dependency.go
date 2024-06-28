package api

import (
	"realTime/server/db"
	"realTime/server/internal/user"
	"realTime/server/internal/websokcet"
	"realTime/server/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type Api struct {
	Addr string
}

func NewApi(addr string) *Api {
	return &Api{
		Addr: addr,
	}
}

func (a *Api) Run() error {

	fiberApp := fiber.New(fiber.Config{
		CaseSensitive:            true,
		EnableSplittingOnParsers: true,
	})
	fiberApp.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}

		return fiber.ErrUpgradeRequired
	})

	api := fiberApp.Group("/api")

	validate := utils.XValidator{
		Validator: validator.New(),
	}

	db, err := db.NewDatabase()
	if err != nil {
		return err
	}

	userRepo := user.NewRepository(db.DB())
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService, validate)
	userHandler.RegisterRoute(api)

	hub := websokcet.NewHub()
	wsHandler := websokcet.NewHandler(hub)
	wsHandler.RegisterRoute(api)

	go hub.Run()

	return fiberApp.Listen(a.Addr)
}
