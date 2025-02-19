package handlers

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"fmt"
)

func CheckNewEmails(c *client.Client) ([]uint32, error) {
	_, err := c.Select("INBOX", false)
	if err != nil {
		return nil, fmt.Errorf("failed to select INBOX: %v", err)
	}

	criteria := &imap.SearchCriteria{WithFlags: []string{"\\Unseen"}}
	ids, err := c.Search(criteria)
	if err != nil {
		return nil, fmt.Errorf("failed to search emails: %v", err)
	}

	return ids, nil
}

func FetchEmailSender(c *client.Client, ids []uint32) ([]string, error){
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

func ProcessMails (email string, token string, senders []string) {
	for _, sender := range senders {
		fmt.Printf("Sending auto-response to: %s\n", sender)
		SendAutoReply(email, token, sender)
	}
}