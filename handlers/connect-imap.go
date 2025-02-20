package handlers

import (
	"fmt"
	"strings"

	"github.com/emersion/go-imap/client"
)

func getIMAPServerLocal(email string) string {
	if strings.Contains(email, "@gmail.com") {
		return "imap.gmail.com:993"
	} else if strings.Contains(email, "@outlook.com") {
		return "outlook.office365.com:993"
	}
	return ""
}

func ConnectToImapWithPassword(email, password string) (*client.Client, error) {
	server := getIMAPServerLocal(email)
	if server == "" {
		return nil, fmt.Errorf("unsupported email provider")
	}

	c, err := client.DialTLS(server, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v", err)
	}

	if err := c.Login(email, password); err != nil {
		c.Logout()
		return nil, fmt.Errorf("authentication failed: %v", err)
	}

	fmt.Println("Logged in successfully!")
	return c, nil
}
