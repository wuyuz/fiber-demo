package auth

import (
	config2 "fiber-demo/config"
	libraries "fiber-demo/lib"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func SendConfirmationEmail(email string, baseURL string) {
	confirmLink := GenerateConfirmURL(email, baseURL)
	htmlBody := config2.PrepareHtml("emails/confirm", fiber.Map{
		"confirm_link": confirmLink,
	})
	config2.Send(email, "Is it you? Please confirm!", htmlBody, "", "")
}

func GenerateConfirmURL(email string, baseURL string) string {
	// 加密url：通过包含加密的邮箱进行判断，后激活
	token := libraries.Encrypt(email, config2.AppConfig.App_Key)
	uri := fmt.Sprintf("%s/do/verify-email?t=%s", baseURL, token)
	return uri
}
