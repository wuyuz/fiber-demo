package contraller

import (
	"fiber-demo/config"
	"github.com/gofiber/fiber/v2"
)

func GetPermission(c *fiber.Ctx) error {
	// 策略
	pAll := config.Auth.Enforcer.GetPolicy()

     	// 角色
	gAll := config.Auth.Enforcer.GetGroupingPolicy()

	return c.Render("auth/permission", fiber.Map{
		"pAll": pAll,
		"gAll": gAll,
	})
}
