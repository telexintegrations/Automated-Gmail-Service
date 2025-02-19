package handlers

import (
	"github.com/emersion/go-imap/client"
	"fmt"
	"strings"
)

type OAuthBearer struct {
	User  string
	Token string
}

func (a *OAuthBearer) Name() string {
	return "XOAUTH2"
}

func (a *OAuthBearer) Start() (string, []byte, error) {
	authString := fmt.Sprintf("user=%s\x01auth=Bearer %s\x01\x01", a.User, a.Token)
	return "", []byte(authString), nil
}

func (a *OAuthBearer) Next(_ []byte) ([]byte, error) {
	return nil, nil
}

func getIMAPServer(email string) string {
	if strings.Contains(email, "@gmail.com") {
		return "imap.gmail.com:993"
	} else {
		return ""
	}
	// else if strings.Contains(email, "@outlook.com") {
	// 	return "outlook.office365.com:993"
	// } 
}

func ConnectToImap(email string, token string) (*client.Client, error) {
	server := getIMAPServer(email)
	if server == "" {
		return nil, fmt.Errorf("unsupported email provider")
	}

	// Connect to email server
	c, err := client.DialTLS(server, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v", err)
	}
	
	auth := &OAuthBearer{User: email, Token: token}
	loginError := c.Authenticate(auth)
	if loginError != nil {
		c.Logout()
		return nil, fmt.Errorf("OAuth authentication failed: %v", loginError)
	}

	return c, nil
}