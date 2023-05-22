package main

import (
	"github.com/MadeBaruna/pom-moe/internal/handlers"
	"github.com/MadeBaruna/pom-moe/internal/queue"
	"github.com/MadeBaruna/pom-moe/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"
)

func main() {
	utils.LoadEnv()
	utils.InitLogger()

	utils.LoadProxy()
	queue.LoadNats()
	defer queue.Drain()

	app := fiber.New(fiber.Config{
		ProxyHeader:             fiber.HeaderXForwardedFor,
		EnableTrustedProxyCheck: true,
		TrustedProxies:          []string{"127.0.0.1", "172.16.0.0/12"},
	})
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	app.Use(func(c *fiber.Ctx) error {
		log.Info().Str("ip", c.IP()).Str("method", c.Method()).Str("url", c.OriginalURL()).Int("status", c.Response().StatusCode()).Msg("request")
		return c.Next()
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://pom.moe, http://localhost:8151",
	}))

	app.Get("/", handlers.ApiStatus)
	app.Post("/warp", handlers.GetWarp)

	err := app.Listen(":8152")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start api server")
	}
}
