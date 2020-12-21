package config

import (
	"github.com/alexedwards/argon2id"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/session/v2"
	"github.com/gofiber/template/django"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog"
	"github.com/sujit-baniya/flash"
)


var App *fiber.App //nolint:gochecknoglobals

var TemplateEngine *django.Engine

//var Ctx *fiber.Ctx

var DB *gorm.DB //nolint:gochecknoglobals

var Hash *HashDriver //nolint:gochecknoglobals

var Flash *flash.Flash

var Session *session.Session
//
//var MailerServer *mail.SMTPServerzytrew/.,
//
//var Mailer *mail.SMTPClient
//
//var Paypal *paypal.Client

//var RedisClient *redis.Client

var PermissionAdapter *gormadapter.Adapter //nolint:gochecknoglobals
//
var Enforcer *casbin.Enforcer //nolint:gochecknoglobals
//
//var Queue *amqp.Connection

var Log *Logger

type HashConfig struct {
	// Argon2id configuration
	Params *argon2id.Params
}

type HashDriver struct {
	// Configuration for the argon2id driver
	Config *HashConfig
}

type Logger struct {
	*zerolog.Logger
}

func NewHashDriver(config ...HashConfig) *HashDriver {
	var cfg HashConfig
	cfg.Params = argon2id.DefaultParams
	if len(config) > 0 {
		cfg = config[0]
	}
	return &HashDriver{Config: &cfg}
}


func (d *HashDriver) Create(password string) (hash string, err error) {
	return argon2id.CreateHash(password, d.Config.Params)
}

func (d *HashDriver) Match(password string, hash string) (match bool, err error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}