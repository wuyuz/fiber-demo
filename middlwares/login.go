package middlwares

import (
	app "fiber-demo/app"
	"fiber-demo/models"
	"fiber-demo/services/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
)

func RedirectToHomePageOnLogin(c *fiber.Ctx) error {
	if auth.IsLoggedIn(c) {
		return c.Redirect("/")
	}
	return c.Next()
}

func ValidateLoginPost(c *fiber.Ctx) error {
	var login models.Login
	if err := c.BodyParser(&login); err != nil {
		return app.Flash.WithError(c, fiber.Map{
			"message": err.Error(),
		}).Redirect("/login")
	}
	v := validate.Struct(login)
	if !v.Validate() {
		return app.Flash.WithError(c, fiber.Map{
			"message": v.Errors.One(),
		}).Redirect("/login")
	}
	user, err := login.CheckLogin() //nolint:wsl

	if err != nil {
		return app.Flash.WithError(c, fiber.Map{
			"message": err.Error(),
		}).Redirect("/login")
	}
	c.Locals("user", user)
	return c.Next()
}
