# Fiber

### Fiber 框架架子

[wuyuz/fiber-demo](https://github.com/wuyuz/fiber-demo)

### Django模版语法

最让我惊艳的是在fiber框架中的模版引擎中居然可以使用Django的模版语法

```go
func main() {
	// Create a new engine
	engine := django.New("./views", ".html")

	// Or from an embedded system
	// See github.com/gofiber/embed for examples
	// engine := html.NewFileSystem(http.Dir("./views", ".django"))
	// Pass the engine to the Views
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		// Render with and extends
		return c.Render("index", fiber.Map{
			"Title": "Hello, World!",
		})
	})

	app.Get("/embed", func(c *fiber.Ctx) error {
		// Render index within layouts/main
		return c.Render("embed", fiber.Map{
			"Title": "Hello, World!",
		}, "layouts/main2")
	})

	log.Fatal(app.Listen(":3000"))
}
```

### Air 热加载

[](https://github.com/cosmtrek/air](https://github.com/cosmtrek/air))

Air是 Go 语言的热加载工具，它可以监听文件或目录的变化，自动编译，重启程序。大大提高开发期的工作效率。

- 在项目中执行安装air，下面命令会在$GOPATH/bin目录下生成air命令。我一般会将$GOPATH/bin加入系统PATH中，所以可以方便地在任何地方执行air命令。

```rust
go get -u github.com/cosmtrek/air
```

- 在项目目录下创建.air.toml文件，默认执行air，就会加载此文件，当然可以使用`-c`来指定文件，还可以配置项目根目录，临时文件目录，编译和执行的命令，监听文件目录，监听后缀名，甚至控制台日志颜色都可以配置。具体参看官网的`air_example.toml`

### Casbin 权限管理

规则编辑器

[Casbin · 一个授权库，支持访问控制模型，如 ACL，RBAC，ABAC，支持 Golang，Java，C/C++，Node.js，JavaScript，PHP，Python，.NET(C＃)，Delphi，Rust，Dart/Flutter 和 Elixir。](https://casbin.org/zh-CN/editor)

Casbin是一个强大的、高效的开源访问控制框架，其权限管理机制支持多种访问控制模型。

Casbin只负责访问控制。身份认证 authentication（即验证用户的用户名、密码），需要其他专门的身份认证组件负责。例如（[jwt-go](https://blog.csdn.net/weixin_43746433/article/details/107881613)）

### TailWindCss 样式管理

[Container - Tailwind CSS](https://www.tailwindcss.cn/docs/container)

- 首先在项目中新建一个tailwind目录，进入目录
- 初始化项目，同时安装tailwind 和相关插件、同时初始化tailwind.config.js文件：

```bash
npm init // 生成package.json
cnpm install tailwindcss postcss-cli autoprefixer
npx tailwindcss init
```

- 新建一个style.css 文件，并引入样式

```css
@tailwind base;
@tailwind components;
@tailwind utilities;
```

- 编辑tailwind.config.js,配置css主题，作用的文件等

```bash
module.exports = {
  purge: [
    '../resource/**/*.html',
    '../resource/*.html',
    '../resource/**/*.js',
    '../resource/**/*.vue',
    '../resource/**/*.scss',
    '../resource/**/*.css',
  ],
  theme: {
    extend: {
      colors: {
        black: '#0f1c33',
      },
      margin: {
        '96': '24rem',
        '128': '32rem',
      },
    }
  },
  variants: {
    tableLayout: ['responsive', 'hover', 'focus'],
  },
  plugins: [],
}
```

- 修改package.json 文件的打包命令，同时指定生成的文件，最后执行npm run build 生成css文件，同时在自己的html中引入即可

```bash
"scripts": {
    "build": "postcss style.css -o ../resource/assets/css/tailwind.css"
  },
```

### Validate 验证器

[gookit/validate](https://github.com/gookit/validate/blob/master/README.zh-CN.md)

Go通用的数据验证与过滤库，使用简单，内置大部分常用验证器、过滤器，支持自定义消息、字段翻译。

- 支持验证 `Map` `Struct` `Request`（`Form`，`JSON`，`url.Values`, `UploadedFile`）数据
- 简单方便，支持前置验证检查, 支持添加自定义验证器
- 支持将规则按场景进行分组设置，不同场景验证不同的字段
- 支持在进行验证前对值使用过滤器进行净化过滤，查看 [内置过滤器](https://github.com/gookit/validate/blob/master/README.zh-CN.md#built-in-filters)
- 已经内置了超多（**>70** 个）常用的验证器，查看 [内置验证器](https://github.com/gookit/validate/blob/master/README.zh-CN.md#built-in-validators)
- 方便的获取错误信息，验证后的安全数据获取(*只会收集有规则检查过的数据*)
- 支持自定义每个验证的错误消息，字段翻译，消息翻译(内置`en` `zh-CN` `zh-TW`)
- 完善的单元测试，测试覆盖率 > 90%