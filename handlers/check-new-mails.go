package handlers

import (
	"fmt"
	"log"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

func CheckNewEmails(c *client.Client, lastUID uint32) ([]uint32, error) {
	_, err := c.Select("INBOX", false)
	if err != nil {
		return nil, fmt.Errorf("failed to select INBOX: %v", err)
	}

	checkErr := c.Check()
	if checkErr != nil {
		log.Println("Error checking for emails: ", checkErr)
	}

	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{"\\Seen"}
	// criteria := &imap.SearchCriteria{WithFlags: []string{"\\Seen"}}
	ids, err := c.Search(criteria)
	if err != nil {
		log.Println("Retrying IMAP search...")
		time.Sleep(5 * time.Second)

		ids, err = c.Search(criteria)
		if err != nil {
			return nil, fmt.Errorf("failed to search emails: %v", err)
		}
	}

	var newEmails []uint32
	for _, id := range ids {
		if id > lastUID {
			newEmails = append(newEmails, id)
		}
	}

	return newEmails, nil
}

func FetchEmailSender(c *client.Client, ids []uint32) ([]string, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	seqset := new(imap.SeqSet)
	seqset.AddNum(ids...)
	messages := make(chan *imap.Message, len(ids))
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	senders := []string{}
	for msg := range messages {
		if msg.Envelope != nil && len(msg.Envelope.From) > 0 {
			senders = append(senders, msg.Envelope.From[0].Address())
		}
	}

	if err := <-done; err != nil {
		return nil, fmt.Errorf("error fetching emails: %v", err)
	}

	return senders, nil
}

// func CheckNewEmails(c *client.Client, lastUID uint32) ([]uint32, error) {
// 	mbox, err := c.Select("INBOX", false)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to select INBOX: %v", err)
// 	}

// 	if mbox.Messages == 0 {
// 		fmt.Printf("No messages in inbox")
// 		return nil, nil
// 	}

// 	seqSet := new(imap.SeqSet)
// 	if mbox.Messages > 50 {
// 		seqSet.AddRange(mbox.Messages-50, mbox.Messages)
// 	} else {
// 		seqSet.AddRange(1, mbox.Messages)
// 	}

// 	items := []imap.FetchItem{imap.FetchEnvelope, imap.FetchFlags, imap.FetchUid}
// 	messages := make(chan *imap.Message, 50)
// 	done := make(chan error, 1)

// 	go func() {
// 		done <- c.Fetch(seqSet, items, messages)
// 	}()

// 	var newEmails []uint32

// 	for msg := range messages {
// 		if msg == nil {
// 			fmt.Println("Received nil message")
// 			continue
// 		}

// 		isUnseen := true
// 		for _, flag := range msg.Flags {
// 			if flag == imap.SeenFlag {
// 				isUnseen = false
// 				break
// 			}
// 		}

// 		if isUnseen && msg.Uid > lastUID {
// 			newEmails = append(newEmails, msg.Uid)
// 		}
// 	}

// 	if err := <-done; err != nil {
// 		return nil, fmt.Errorf("failed to fetch emails: %v", err)
// 	}

// 	if len(newEmails) == 0 {
// 		fmt.Println("No new unseen emails found.")
// 	}

// 	return newEmails, nil
// }

// func FetchEmailSender(c *client.Client, ids []uint32) ([]string, error) {
// 	if len(ids) == 0 {
// 		fmt.Printf("No email IDs for fetching senders")
// 		return nil, nil
// 	}

// 	seqset := new(imap.SeqSet)
// 	seqset.AddNum(ids...)

// 	items := []imap.FetchItem{imap.FetchEnvelope, imap.FetchFlags, imap.FetchUid}

// 	messages := make(chan *imap.Message, len(ids))
// 	done := make(chan error, 1)

// 	go func() {
// 		done <- c.Fetch(seqset, items, messages)
// 	}()

// 	senders := []string{}

// 	for msg := range messages {
// 		if msg == nil {
// 			fmt.Println("Warning: Received nil message")
// 			continue
// 		}

// 		if msg.Envelope == nil {
// 			fmt.Println("Warning: Received nil message", msg.Uid)
// 			continue
// 		}

// 		if len(msg.Envelope.From) == 0 {
// 			fmt.Println("Warning: Received nil message", msg.Uid)
// 			continue
// 		}

// 		senderEmail := msg.Envelope.From[0].Address()
// 		fmt.Printf("Found sender: %s\n", senderEmail)
// 		senders = append(senders, senderEmail)
// 	}

// 	if err := <-done; err != nil {
// 		return nil, fmt.Errorf("error fetching emails: %v", err)
// 	}

// 	fmt.Printf("Found senders: %v\n", senders)
// 	return senders, nil
// }

func ProcessMails(email string, token string, username string, senders []string) {
	for _, sender := range senders {
		fmt.Printf("Sending auto-response to: %s\n", sender)
		SendAutoReply(email, token, username, sender)
	}
}

func MarkEmailsAsSeen(c *client.Client, ids []uint32) error {
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(ids...)

	item := imap.FormatFlagsOp(imap.AddFlags, true)
	flags := []interface{}{imap.SeenFlag}

	return c.Store(seqSet, item, flags, nil)
}
