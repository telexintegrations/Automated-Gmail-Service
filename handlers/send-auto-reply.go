package handlers

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"strings"
)

func Capitalise(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func SendAutoReply(email string, token string, username string, to string) error {
	var SMTPHost string
	var SMTPPort int

	if strings.Contains(email, "@gmail.com") {
		SMTPHost = "smtp.gmail.com"
		SMTPPort = 587
	} else if strings.Contains(email, "@outlook.com") {
		SMTPHost = "smtp.office365.com"
		SMTPPort = 587
	} else {
		return fmt.Errorf("unsupported email provider")
	}

	var capitalisedUsername = Capitalise(username)

	subject := "Thank You for Your Messsage!"
	body := fmt.Sprintf("Hello, \n\nThank you for reaching out to us. We have received your message and will get back to you as soon as possible. \n\nBest regards, \n%s", capitalisedUsername)

	message := gomail.NewMessage()
	message.SetHeader("From", email)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", body)

	dialer := gomail.NewDialer(SMTPHost, SMTPPort, email, token)

	err := dialer.DialAndSend(message)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	fmt.Printf("Autoreply sent to: %v \n", to)
	return nil
}
