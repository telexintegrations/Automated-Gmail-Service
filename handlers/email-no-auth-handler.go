package handlers

import (
	"log"
	"time"
)

func EmailNoAuthHandler(email string, password string, username string) {
	c, err := ConnectToImapWithPassword(email, password)
	if err != nil {
		log.Println("IMAP Connection Error:", err)
		return
	}
	defer c.Logout()

	var lastUID uint32

	for {
		err := c.Check()
		if err != nil {
			log.Println("Error checking for emails: ", err)
		}

		ids, err := CheckNewEmails(c, lastUID)
		log.Printf("Found mails with ids: %v\n", ids)

		if err != nil {
			log.Println("Error checking mails:", err)
			time.Sleep(time.Second * 10)
			continue
		}

		if len(ids) == 0 {
			log.Println("No new emails found")
		} else {
			senders, err := FetchEmailSender(c, ids)
			log.Printf("senders: %v\n", senders)
			if err != nil {
				log.Println("Error fetching mail sender:", err)
				time.Sleep(time.Second * 10)
				continue
			}

			ProcessMails(email, password, username, senders)

			markErr := MarkEmailsAsSeen(c, ids)
			if markErr != nil {
				log.Println("Error marking emails as seen: ", markErr)
			}

			if len(ids) > 0 {
				lastUID = ids[len(ids)-1]
			}
		}

		time.Sleep(time.Minute * 1)
		log.Println("Sleeping for one minute")
	}
}
