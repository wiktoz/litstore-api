package emails

import (
	"os"

	"gopkg.in/gomail.v2"
)

func SendPasswordResetEmail(email, token string) error {
	/*
		    FOR REACT TEMPLATES

			templateBytes, err := os.ReadFile("emails/reset-password.html")
			if err != nil {
				return err
			}

			body := string(templateBytes)
			body = strings.ReplaceAll(body, "{{.UserName}}", userName)
			body = strings.ReplaceAll(body, "{{.ResetLink}}", "Your token" + token)
	*/

	body := "Your password reset token is: " + token

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL_USER"))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Reset Your Password")
	m.SetBody("text/html", body)

	d := gomail.NewDialer(os.Getenv("EMAIL_HOST"), 587, os.Getenv("EMAIL_USER"), os.Getenv("EMAIL_PASSWORD"))

	return d.DialAndSend(m)
}

func SendVerificationEmail(email, token string) error {
	/*
		    FOR REACT TEMPLATES

			templateBytes, err := os.ReadFile("emails/verify-email.html")
			if err != nil {
				return err
			}

			body := string(templateBytes)
			body = strings.ReplaceAll(body, "{{.UserName}}", userName)
			body = strings.ReplaceAll(body, "{{.VerificationLink}}", "Your token" + token)
	*/

	body := "Your email verification token is: " + token

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL_USER"))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Verify Your Email")
	m.SetBody("text/html", body)

	d := gomail.NewDialer(os.Getenv("EMAIL_HOST"), 587, os.Getenv("EMAIL_USER"), os.Getenv("EMAIL_PASSWORD"))

	return d.DialAndSend(m)
}
