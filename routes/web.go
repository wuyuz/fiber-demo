package routes

import (
	. "fiber-demo/app"
	"fiber-demo/config"
	"fiber-demo/contraller"
	"fiber-demo/middlwares"
	"fiber-demo/services/auth"
	"github.com/gofiber/fiber/v2"
)

func WebRoutes() {
	web := App.Group("")
	web.Use(auth.AuthCookie)
	LandingRoutes(web)
}


func LandingRoutes(app fiber.Router) {
	app.Use(middlwares.Authenticate(middlwares.AuthConfig{
		SigningKey:  []byte(config.AuthConfig.App_Jwt_Secret),
		TokenLookup: "cookie:fiber-demo-Token",
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			auth.Logout(ctx)
			return ctx.Next()
		},
	}))

	app.Get("/", contraller.Landing)
	app.Get("/permission",
		contraller.GetPermission,
		)
}
