package main

import (
	"fiber-demo/config"
	. "fiber-demo/app"
	libraries "fiber-demo/lib"
	"fiber-demo/middlwares"
	"fiber-demo/migration"
	"fiber-demo/routes"
	"flag"
	"github.com/gofiber/fiber/v2"

	//"github.com/gofiber/fiber/v2"
)

func main() {
	Log = libraries.SetupZeroLog()
	migrate := flag.Bool("migrate", false, "Migrate the pending migration files")
	flag.Parse()
	if *migrate {
		// 	初始化迁移数据库表结构
		migration.InitMigrate()
		return
	}
	Serve()
}

func Serve() {
	Boot()
	App.Use(middlwares.LogMiddleware)
	// 路由加载
	routes.LoadRouter()
	// 视图加载后拦截器
	App.Use(func(c *fiber.Ctx) error {
		var err fiber.Error
		err.Code = fiber.StatusNotFound
		return config.CustomErrorHandler(c, &err)
	})
	// go libraries.Consume("webhook-callback")               //nolint:wsl
	err := App.Listen(":" + config.AppConfig.App_Port) //nolint:wsl
	if err != nil {
		panic("App not starting: " + err.Error() + "on Port: " + config.AppConfig.App_Port)
	}

	defer DB.Close()
}

func Boot() {
	// 加载环境变量
	config.LoadEnv()
	// 初始化app
	config.BootApp()
}

