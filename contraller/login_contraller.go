package contraller

import (
	. "fiber-demo/app"
	"fiber-demo/config"
	"fiber-demo/models"
	"fiber-demo/services/auth"
	"github.com/gofiber/fiber/v2" //nolint:goimports
)

func LoginGet(c *fiber.Ctx) error {
	// 清空缓存
	Flash.Get(c)
	// 注意使用django的模版后就不用给第三个参数了，原本的模版渲染的第三个参数是给继承模版的
	return c.Render("auth/login",nil)
}

func LoginPost(c *fiber.Ctx) error { //nolint:wsl
	user := c.Locals("user").(*models.User)
	_, _ = auth.Login(c, user.ID, 	config.AuthConfig.App_Jwt_Secret) //nolint:wsl
	return c.Redirect("/")
}

func LogoutPost(c *fiber.Ctx) error { //nolint:nolintlint,wsl
	if auth.IsLoggedIn(c) {
		_ = auth.Logout(c)
	}
	return c.Redirect("/login")
}

