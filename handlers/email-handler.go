package handlers

import (
	"log"
	"time"
	"net/http"
	"fmt"
	"encoding/json"
)

type GoogleUserInfo struct {
	Email string `json:"email"`
}

func EmailHandler(email string, token string) {
	// Connect to email server
	c, err := ConnectToImap(email, token)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Logout()

	for {
		ids, err := CheckNewEmails(c)
		if err != nil {
			log.Println("Error checking mails: ", err)
			continue
		}

		senders, err := FetchEmailSender(c, ids)
		if err != nil {
			log.Println("Error fetching mail sender: ", err)
			continue
		}

		ProcessMails(email, token, senders)

		time.Sleep(time.Minute * 1) // Wait for 1 minute before checking again
	}
}

func GetEmailFromToken (accessToken string) (string, error) {
	url := "https://www.googleapis.com/oauth2/v2/userinfo"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer " + accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error retrieving user info: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch user email from token")
	}

	var userInfo GoogleUserInfo
	parseError := json.NewDecoder(resp.Body).Decode(&userInfo)
	if parseError != nil {
		return "", fmt.Errorf("failed to parse user info")
	}

	if userInfo.Email == "" {
		return "", fmt.Errorf("failed to get user email")
	}

	return userInfo.Email, nil
}