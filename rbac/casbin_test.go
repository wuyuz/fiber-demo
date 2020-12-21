package rbac

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/util"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"testing"
)

var authz *PermissionMiddleware
var PermissionAdapter *gormadapter.Adapter //nolint:gochecknoglobals


func testKeyMatch2( key1 string, key2 string) {
	myRes := util.KeyMatch2(key1, key2)
	fmt.Println(myRes,key1)
}

func TestRequiresPermissions(t *testing.T) {
	app := fiber.New()

	var err error

	// 连接数据库，创建一个casbin的库和一张casbin_rule的表，记录其中的规则和用户关系
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/", "root", "root1234", "127.0.0.1", 3306) //nolint:wsl,lll
	PermissionAdapter, err = gormadapter.NewAdapter("mysql", connectionString)

	if err != nil {
		panic(fmt.Sprintf("failed to initialize casbin adapter: %v", err))
	}
	Enforcer, _ := casbin.NewEnforcer("model.conf", PermissionAdapter) //nolint:wsl

	authz = &PermissionMiddleware{
		Enforcer:      Enforcer, //nolint:gofmt
		PolicyAdapter: PermissionAdapter,
		Lookup: func(ctx *fiber.Ctx) string {
			// 这里需要给出一个计算出当前用户的方法，如果没有需要给一个 最低的默认用户
			user := ctx.Query("user")
			return user
		},
		// 没有认证触发
		Unauthorized: func(c *fiber.Ctx) error {
			var err fiber.Error
			err.Code = fiber.StatusUnauthorized
			err.Message = "Unauthorized"
			return CustomErrorHandler(c, &err)
		},

		// 无权限
		Forbidden: func(c *fiber.Ctx) error {
			var err fiber.Error
			err.Code = fiber.StatusForbidden
			err.Message = "Forbidden!"
			return CustomErrorHandler(c, &err)
		},
	}

	// 验证角色，不判断路径和方法
	app.Get("/blog/:id",
		// 由于该方法中设置的是MatchAll，所以默认必须同时拥有两种角色才可以
		authz.RequiresRoles([]string{"admin", "user"}),
		func(c *fiber.Ctx) error {
			return c.SendString(fmt.Sprintf("Blog updated with Id: %s", c.Params("id")))
		},
	)

	app.Get("/book/:id",
		authz.RequiresRoles([]string{"admin", "user"}, AtLeastOne), // 修改为至少一种角色满足即可
		func(c *fiber.Ctx) error {
			return c.SendString(fmt.Sprintf("Book updated with Id: %s", c.Params("id")))
		},
	)

	// check permission with Method and Path
	// 验证当前用户的角色，和是否具有当前访问路径和方法的权限
	app.Post("/blog",
		authz.RoutePermission(),
		func(c *fiber.Ctx) error {
			// your handler
			return c.SendString(fmt.Sprintf("Blog updated with Id: %s", c.Params("user")))
		},
	)

	app.Post("/blog/:id",
		authz.RequiresPermissions([]string{"blog:create", "blog:update"}, AtLeastOne),
		func(c *fiber.Ctx) error {
			// your handler
			return c.SendString(fmt.Sprintf("Book updated with Id: %s", c.Params("id")))
		},
	)

	// 框架API规则
	// 添加规则
	app.Post("/p/add",
		authz.RequiresRoles([]string{"admin"}),
		func(c *fiber.Ctx) error {

			if ok,_:= authz.Enforcer.AddPolicy("admin","/api/v1/world","GET");!ok{
				fmt.Println("Policy already exist")
			}
			return c.SendString(fmt.Sprint("Policy add OK"))
		},
	)

	// 添加用户角色
	app.Post("/g/add",
		authz.RequiresRoles([]string{"admin"}),
		func(c *fiber.Ctx) error {
			if ok,_:= authz.Enforcer.AddRoleForUser("wu","admin");!ok{
				fmt.Println("g already exist")
			}
			return c.SendString(fmt.Sprint("g add OK"))
		},
	)

	// 获取角色相应的权限
	app.Get("/getPermissionForUser",
		authz.RequiresRoles([]string{"admin"}),
		func(c *fiber.Ctx) error {
			// 代码中调用自带的过滤函数
			testKeyMatch2("/foo/bar", "/foo/*")  // true
			testKeyMatch2( "/resource2", "/:resource")  // true
			user := authz.Enforcer.GetPermissionsForUser(c.Params("user"))
			return c.JSON(fiber.Map{
				"profile":user,
			})
		},
	)

	app.Listen(":8089")
}

func CustomErrorHandler(c *fiber.Ctx, err error) error {
	// StatusCode defaults to 500
	code := fiber.StatusInternalServerError
	//nolint:misspell    // Retrieve the custom statuscode if it's an fiber.*Error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	} //nolint:gofmt,wsl

	return c.Status(code).JSON(err)

}
