package config

import (
	"crypto/rand"
	. "fiber-demo/app"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/django"
	"github.com/sujit-baniya/flash"
	"net/http"
	"os"
	"path/filepath"
)

// 应用配置结构体
type AppConfiguration struct {
	App_Name        string
	App_Upload_Path string  // app上传文件路径
	App_Upload_Size int   // 文件大小
	App_Env         string
	App_Key         string
	App_Url         string  // 路径
	App_Port        string  // 端口
}

var AppConfig *AppConfiguration //nolint:gochecknoglobals

var BlastlistedDomains = []string{}

func LoadAppConfig() {
	loadDefaultConfig()
	// 反序列化成配置对象，全局变量
	ViperConfig.Unmarshal(&AppConfig)
	if AppConfig.App_Url == "" {
		AppConfig.App_Url = fmt.Sprintf("http://localhost:%s", AppConfig.App_Port)
	}
	AppConfig.App_Upload_Path = filepath.Join(".", AppConfig.App_Upload_Path)
	AppConfig.App_Upload_Size = AppConfig.App_Upload_Size * 1024 * 1024
	if _, err := os.Stat(AppConfig.App_Upload_Path); os.IsNotExist(err) {
		os.MkdirAll(AppConfig.App_Upload_Path, os.ModePerm)  // 如果没有创建则新建
	}
}

func loadDefaultConfig() {
	// 设置变量到viperConfig中
	ViperConfig.SetDefault("APP_NAME", "fiber-demo")
	ViperConfig.SetDefault("APP_ENV", "dev")
	ViperConfig.SetDefault("APP_UPLOAD_PATH", "uploads")
	ViperConfig.SetDefault("APP_UPLOAD_SIZE", 4)
	ViperConfig.SetDefault("APP_KEY", "1894cde6c936a294a478cff0a9227fd276d86df6533b51af5dc59c9064edf426")
	ViperConfig.SetDefault("APP_PORT", "8080")
}

func GenerateAppKey(length int) {
	key := make([]byte, length)

	_, err := rand.Read(key)
	if err != nil {
		// handle error here
	}
}

func BootApp() {
	// 加载环境变量
	LoadAppConfig()
	// 初始化模版文件系统,使用django模版，👍
	TemplateEngine = django.NewFileSystem(http.Dir("resource/views"), ".html")
	App = fiber.New(fiber.Config{
		ErrorHandler:          CustomErrorHandler,
		ServerHeader:          "fiber-demo",
		Prefork:               true,
		DisableStartupMessage: false,
		Views:                 TemplateEngine,
		BodyLimit:             AppConfig.App_Upload_Size,
	})

	App.Use(pprof.New())  // fiber框架性能分析中间件
	App.Use(LoadHeaders)  // 请求头
	App.Use(recover.New())  // fiber报错拦截器
	App.Use(compress.New(compress.Config{  // 中间件压缩文件程度
		Next:  nil,
		Level: compress.LevelBestSpeed,
	}))
	/*App.Use(csrf.New(csrf.Config{
		CookieSecure:   true,
	}))*/

	App.Static("/assets", "resource/assets", fiber.Static{
		Compress: true,
	})

	App.Use(LoadCacheHeaders)
	// 初始化Hash实例，方便后面加解密
	Hash = NewHashDriver()

	// 启动数据库
	_, err := SetupDB()
	if err != nil {
		panic(err)
	}
	// 权限加载
	SetupPermission()

	// session加载
	LoadSession()

	// 闪现，一个消息存储点，闪存，取一次就没有，存在session中
	Flash = &flash.Flash{
		CookiePrefix: "fiber-demo",
	}
}

// 自定义错误返回
func CustomErrorHandler(c *fiber.Ctx, err error) error {
	// StatusCode defaults to 500
	code := fiber.StatusInternalServerError
	//nolint:misspell    // Retrieve the custom statuscode if it's an fiber.*Error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	} //nolint:gofmt,wsl
	if c.Is("json") {
		return c.Status(code).JSON(err)
	} else {
		return c.Status(code).Render(fmt.Sprintf("errors/%d", code), fiber.Map{ //nolint:nolintlint,errcheck
			"error": err,
		})
	}
}
