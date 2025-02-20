package handlers

import (
	"fmt"
	"log"
	"time"
)

func EmailNoAuthHandler(email string, password string) {
	c, err := ConnectToImapWithPassword(email, password)
	if err != nil {
		log.Println("IMAP Connection Error:", err)
		return
	}
	defer c.Logout()

	processedEmails := make(map[uint32]bool)

	for {
		err := c.Check()
		if err != nil {
			log.Println("Error checking for emails: ", err)
		}

		ids, err := CheckNewEmails(c)
		fmt.Printf("Found mails with ids: %v\n", ids)

		if err != nil {
			log.Println("Error checking mails:", err)
			time.Sleep(time.Second * 10)
			continue
		}

		var newMails []uint32
		for _, id := range ids {
			if !processedEmails[id] {
				newMails = append(newMails, id)
			}
		}

		if len(newMails) == 0 {
			fmt.Println("No new emails found")
			time.Sleep(time.Second * 60)
			continue
		} else {
			senders, err := FetchEmailSender(c, ids)
			fmt.Printf("senders: %v\n", senders)
			if err != nil {
				log.Println("Error fetching mail sender:", err)
				time.Sleep(time.Second * 10)
				continue
			}
			
			ProcessMails(email, password, senders)

			for _, id := range newMails {
				processedEmails[id] = true
			}

			markErr := MarkEmailsAsSeen(c, newMails)
			if markErr != nil {
				log.Println("Error marking emails as seen: ", markErr)
			}
		}

		time.Sleep(time.Minute * 1)
		fmt.Println("Sleeping for one minute")
	}
}
