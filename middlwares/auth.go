package middlwares

import (
	"errors"
	"fmt"
	config2 "fiber-demo/config"
	"log"
	"reflect"
	"strings"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

// Config defines the config for BasicAuth middleware
type AuthConfig struct {
	// Filter defines a function to skip middleware.
	// Optional. Default: nil
	Filter func(*fiber.Ctx) bool

	// SuccessHandler defines a function which is executed for a valid token.
	// Optional. Default: nil
	SuccessHandler func(*fiber.Ctx) error

	// ErrorHandler defines a function which is executed for an invalid token.
	// It may be used to define a custom JWT error.
	// Optional. Default: 401 Invalid or expired JWT
	ErrorHandler func(*fiber.Ctx, error) error

	// Signing key to validate token. Used as fallback if SigningKeys has length 0. // 验证令牌，如果长度为0则用作备用
	// Required. This or SigningKeys.
	SigningKey interface{}

	// Map of signing keys to validate token with kid field usage.
	// Required. This or SigningKey.
	SigningKeys map[string]interface{} // key和keys二者选一

	// Signing method, used to check token signing method.
	// Optional. Default: "HS256".
	// Possible values: "HS256", "HS384", "HS512", "ES256", "ES384", "ES512", "RS256", "RS384", "RS512"
	SigningMethod string

	// Context key to store user information from the token into context.
	// Optional. Default: "user".
	ContextKey string  // 存储用户信息默认名字，存在上下文中

	// Claims are extendable claims data defining token content.
	// Optional. Default value jwt.MapClaims
	Claims jwt.Claims  // 定义令牌内容

	// TokenLookup is a string in the form of "<source>:<name>" that is used
	// to extract token from the request.
	// Optional. Default value "header:Authorization".
	// Possible values:
	// - "header:<name>"
	// - "query:<name>"
	// - "param:<name>"
	// - "cookie:<name>"
	TokenLookup string  // 认证的查询方式

	// AuthScheme to be used in the Authorization header.
	// Optional. Default: "Bearer".
	AuthScheme string // token前缀

	keyFunc jwt.Keyfunc
}

// New ...
func Authenticate(config ...AuthConfig) func(*fiber.Ctx) error {
	// Init config
	var cfg AuthConfig
	if len(config) > 0 {
		cfg = config[0]
	}
	if cfg.SuccessHandler == nil {
		// 如果认证成功后的successhandler为nil，则直接通过
		cfg.SuccessHandler = func(c *fiber.Ctx) error {
			return c.Next()
		}
	}
	if cfg.ErrorHandler == nil {
		// 错误处理handler
		cfg.ErrorHandler = func(c *fiber.Ctx, err error) error {
			var er fiber.Error
			if err.Error() == "Missing or malformed JWT" {
				er.Code = fiber.StatusBadRequest
			} else {
				er.Code = fiber.StatusUnauthorized
				return c.SendString("Invalid or expired JWT")
			}
			er.Message = err.Error()
			return config2.CustomErrorHandler(c, &er)
		}
	}
	// 判断唯一key
	if cfg.SigningKey == nil && len(cfg.SigningKeys) == 0 {
		log.Fatal("Fiber: JWT middleware requires signing key")
	}
	if cfg.SigningMethod == "" {
		cfg.SigningMethod = "HS256"
	}
	if cfg.ContextKey == "" {
		cfg.ContextKey = "user"
	}
	if cfg.Claims == nil {
		cfg.Claims = jwt.MapClaims{}
	}
	if cfg.TokenLookup == "" {
		cfg.TokenLookup = "header:" + fiber.HeaderAuthorization
	}
	if cfg.AuthScheme == "" {
		cfg.AuthScheme = "Bearer"
	}
	cfg.keyFunc = func(t *jwt.Token) (interface{}, error) {
		// Check the signing method
		// 检测签名方法
		if t.Method.Alg() != cfg.SigningMethod {
			return nil, fmt.Errorf("Unexpected jwt signing method=%v", t.Header["alg"])
		}
		if len(cfg.SigningKeys) > 0 {
			if kid, ok := t.Header["kid"].(string); ok {
				if key, ok := cfg.SigningKeys[kid]; ok {
					return key, nil
				}
			}
			return nil, fmt.Errorf("Unexpected jwt key id=%v", t.Header["kid"])
		}
		return cfg.SigningKey, nil
	}
	// Initialize
	parts := strings.Split(cfg.TokenLookup, ":")
	extractor := jwtFromHeader(parts[1], cfg.AuthScheme)
	switch parts[0] {
	case "query":
		extractor = jwtFromQuery(parts[1])
	case "param":
		extractor = jwtFromParam(parts[1])
	case "cookie":
		extractor = jwtFromCookie(parts[1])
	}
	// Return middleware handler
	// 验证token
	return func(c *fiber.Ctx) error {
		// Filter request to skip middleware
		// 	一个跳过中间件的方法
		if cfg.Filter != nil && cfg.Filter(c) {
			return c.Next()
		}
		auth, err := extractor(c)
		if err != nil {
			return cfg.ErrorHandler(c, err)
		}
		token := new(jwt.Token)
		if _, ok := cfg.Claims.(jwt.MapClaims); ok {
			token, err = jwt.Parse(auth, cfg.keyFunc)
		} else {
			t := reflect.ValueOf(cfg.Claims).Type().Elem()
			claims := reflect.New(t).Interface().(jwt.Claims)
			token, err = jwt.ParseWithClaims(auth, claims, cfg.keyFunc)
		}
		if err == nil && token.Valid {
			// Store user information from token into context.
			c.Locals(cfg.ContextKey, token)
			return cfg.SuccessHandler(c)
		}
		return cfg.ErrorHandler(c, err)
	}
}

// jwtFromHeader returns a function that extracts token from the request header.
func jwtFromHeader(header string, authScheme string) func(c *fiber.Ctx) (string, error) {
	return func(c *fiber.Ctx) (string, error) {
		auth := c.Get(header)
		l := len(authScheme)
		if len(auth) > l+1 && auth[:l] == authScheme {
			return auth[l+1:], nil
		}
		return "", errors.New("Missing or malformed JWT")
	}
}

// jwtFromQuery returns a function that extracts token from the query string.
func jwtFromQuery(param string) func(c *fiber.Ctx) (string, error) {
	return func(c *fiber.Ctx) (string, error) {
		token := c.Query(param)
		if token == "" {
			return "", errors.New("Missing or malformed JWT")
		}
		return token, nil
	}
}

// jwtFromParam returns a function that extracts token from the url param string.
func jwtFromParam(param string) func(c *fiber.Ctx) (string, error) {
	return func(c *fiber.Ctx) (string, error) {
		token := c.Params(param)
		if token == "" {
			return "", errors.New("Missing or malformed JWT")
		}
		return token, nil
	}
}

// jwtFromCookie returns a function that extracts token from the named cookie.
func jwtFromCookie(name string) func(c *fiber.Ctx) (string, error) {
	return func(c *fiber.Ctx) (string, error) {
		token := c.Cookies(name)
		if token == "" {
			return "", errors.New("Missing or malformed JWT")
		}
		return token, nil
	}
}
