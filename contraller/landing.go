package contraller

import (
	"fiber-demo/services/auth"
	"github.com/gofiber/fiber/v2"
)

func Landing(c *fiber.Ctx) error {
	user, _ := auth.User(c)
	return c.Render("index", fiber.Map{
		"auth": user != nil,
		"user": user,
	})
}

