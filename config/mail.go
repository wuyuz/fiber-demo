package config

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	. "fiber-demo/app"
	"fiber-demo/mail"
	"github.com/valyala/bytebufferpool"
	"log"
	"time"
)

type MailConfiguration struct {
	Mail_Host         string
	Mail_Port         int
	Mail_Username     string
	Mail_Password     string
	Mail_Encryption   string
	Mail_From_Address string
	Mail_From_Name    string
}

var MailConfig *MailConfiguration //nolint:gochecknoglobals

func LoadMailConfig() {
	loadDefaultMailConfig()
	ViperConfig.Unmarshal(&MailConfig)
}

func loadDefaultMailConfig() {
	ViperConfig.SetDefault("MAIL_HOST", "smtp.163.com")
	ViperConfig.SetDefault("MAIL_PORT", "465")
	ViperConfig.SetDefault("MAIL_USERNAME", "python_wlx@163.com")
	ViperConfig.SetDefault("MAIL_PASSWORD", "root123")
	ViperConfig.SetDefault("MAIL_ENCRYPTION", "ssl")
	ViperConfig.SetDefault("MAIL_FROM_ADDRESS", "wanglx<python_wlx@163.com>")
	ViperConfig.SetDefault("MAIL_FROM_NAME", "demo")
}

func Send(to string, subject string, body string, cc string, from string) {
	if MailerServer == nil {
		SetupMailer()
	}
	//New email simple html with inline and CC
	email := mail.NewMSG()
	if from == "" { //nolint:wsl
		from = "wanglx <python_wlx@163.com>"
	}
	email.SetFrom(from). //nolint:wsl
		AddTo(to).
		SetSubject(subject)
	if cc != "" { //nolint:wsl
		email.AddCc(cc)
	}
	email.SetBody(mail.TextHTML, body) //nolint:wsl

	//Call Send and pass the client
	err := email.Send(Mailer)

	if err != nil {
		log.Println(err)
	} else {
		log.Println("Email Sent")
	}
}

func PrepareHtml(view string, body fiber.Map) string {
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)
	// app.Settings.Views.Render
	if err := TemplateEngine.Render(buf, view, body); err != nil {
		// handle err
		fmt.Println("解析错误")
	}
	return buf.String()
}

func SetupMailer() {
	LoadMailConfig()
	var err error
	MailerServer = mail.NewSMTPClient()
	MailerServer.Host = MailConfig.Mail_Host
	MailerServer.Port = MailConfig.Mail_Port
	MailerServer.Username = MailConfig.Mail_Username
	MailerServer.Password = MailConfig.Mail_Password
	if MailConfig.Mail_Encryption == "tls" {
		MailerServer.Encryption = mail.EncryptionTLS
	} else {
		MailerServer.Encryption = mail.EncryptionSSL
	}

	//Variable to keep alive connection
	MailerServer.KeepAlive = false

	//Timeout for connect to SMTP Server
	MailerServer.ConnectTimeout = 10 * time.Second

	//Timeout for send the data and wait respond
	MailerServer.SendTimeout = 10 * time.Second
	Mailer, err = MailerServer.Connect()
	if err != nil {
		log.Print(err)
	}
}
