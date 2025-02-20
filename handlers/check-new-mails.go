package handlers

import (
	"fmt"
	"log"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

func CheckNewEmails(c *client.Client) ([]uint32, error) {
	_, err := c.Select("INBOX", false)
	if err != nil {
		return nil, fmt.Errorf("failed to select INBOX: %v", err)
	}

	criteria := &imap.SearchCriteria{WithFlags: []string{"\\Seen"}}

	ids, err := c.Search(criteria)
	// ids, err := c.UidSearch(criteria)
	if err != nil {
		log.Println("Retrying IMAP search...")
		time.Sleep(5 * time.Second)
	}
	if len(ids) == 0 {
		log.Println("No emails found. Retrying search....")
		criteria := &imap.SearchCriteria{WithFlags: []string{"\\Seen"}}
		ids, err := c.Search(criteria)
		if err != nil {
			return nil, fmt.Errorf("failed to search emails: %v", err)
		}
		return ids, nil
	}

	return ids, nil
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
