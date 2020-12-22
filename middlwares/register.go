package middlwares

import (
	config "fiber-demo/app"
	config2 "fiber-demo/config"
	libraries "fiber-demo/lib"
	"fiber-demo/models"
	"fiber-demo/services/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
)

func ValidateRegisterPost(c *fiber.Ctx) error {
	var register models.RegisterForm

	if err := c.BodyParser(&register); err != nil {
		return config.Flash.WithError(c, fiber.Map{"message": err.Error()}).Redirect("/login")
	}

	v := validate.Struct(register) // 校验器

	if !v.Validate() {
		return config.Flash.WithError(c, fiber.Map{
			"message": v.Errors.One(), // 返回一条报错，验证规则在form中写
		}).Redirect("/register")
	}
	c.Locals("register", register)  // 存储数据
	return c.Next()

}

// 验证token是否正确，邮箱验证token
func ValidateConfirmToken(c *fiber.Ctx) error {
	t := libraries.Decrypt(c.Query("t"), config2.AppConfig.App_Key)
	user, err := models.GetUserByEmail(t)
	if err != nil {
		return config.Flash.WithError(c, fiber.Map{
			"message": err.Error(),
		}).Redirect("/login")
	}

	if user.EmailVerified {
		return config.Flash.WithError(c, fiber.Map{
			"message": "Email was already validated",
		}).Redirect("/login")
	}
	user.EmailVerified = true
	config.DB.Save(&user)

	// 登陆用户
	auth.Login(c, user.ID, config2.AuthConfig.App_Jwt_Secret) //nolint:wsl
	return c.Next()
}
