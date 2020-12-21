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

}
