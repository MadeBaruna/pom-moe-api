package handlers

import (
	"time"

	"github.com/MadeBaruna/pom-moe/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func ApiStatus(c *fiber.Ctx) error {
	i, _ := utils.GetProxy(c.IP())

	return c.JSON(fiber.Map{
		"name":  "Pom.moe API",
		"time":  time.Now().Format(time.RFC3339),
		"index": i,
	})
}
