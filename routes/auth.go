package routes

import (
	. "fiber-demo/app"
	"fiber-demo/contraller"
	"fiber-demo/middlwares"
)
func AuthRoutes() {
	App.Get("/login",
		middlwares.RedirectToHomePageOnLogin,
		contraller.LoginGet,
	)
	App.Post("/do/login",
		middlwares.ValidateLoginPost,
		contraller.LoginPost,
	)

	App.Get("/register", middlwares.RedirectToHomePageOnLogin, contraller.RegisterGet)
	App.Post("/do/register",
		middlwares.RedirectToHomePageOnLogin,
		middlwares.ValidateRegisterPost, // 验证参数
		contraller.RegisterPost,
	)

	App.Get("/do/verify-email/",
		middlwares.ValidateConfirmToken,
		contraller.VerifyRegisteredEmail,
	)

}
