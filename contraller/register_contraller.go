package contraller

import (
	config "fiber-demo/app"
	"fiber-demo/models"
	"fiber-demo/services/auth"
	"github.com/gofiber/fiber/v2"
)

func RegisterGet(c *fiber.Ctx) error {
	return c.Render("auth/register", fiber.Map{"title": "注册","message":config.Flash.Get(c)})
}


func RegisterPost(c *fiber.Ctx)  error {
	// 获取注册数据
	register := c.Locals("register").(models.RegisterForm)
	user, err := register.Signup()  // 存储用户数据，但没有激活
	if err != nil {
		//return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": true, "message": "Error on register request", "data": err.Error()}) //nolint:errcheck
		return config.Flash.WithError(c, fiber.Map{
			"message":err.Error(),
		}).Redirect("/register")
	}
	// 发送邮箱认证
	go auth.SendConfirmationEmail(user.Email, c.BaseURL())
	return c.Redirect("/login")
}


func VerifyRegisteredEmail(c *fiber.Ctx) error {
	return c.Redirect("/")
}