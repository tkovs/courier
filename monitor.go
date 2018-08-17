package courier

import (
	"fmt"
	"time"

	"github.com/getsentry/raven-go"

	_ "github.com/lib/pq"
)

type Job struct {
	Message Message
	Sender  string
}

type Monitor struct {
}

func (m *Monitor) Start() error {
	var err error

	ticker := time.NewTicker(5 * time.Second)
	messages := []MessageModel{}
	mediator := NewMediator()

	db, err := Migrate()
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		return err
	}

	for _ = range ticker.C {
		err = db.Select(&messages, "SELECT * FROM message WHERE (scheduledto < CURRENT_TIMESTAMP) AND (status = 0)")

		for _, message := range messages {
			account, err := message.GetAccount(db)
			if err != nil {
				return err
			}

			err = mediator.SendMessage(Job{
				Message: Message{
					Content:   message.Message,
					Recipient: message.ReceiverPhone,
				},
				Sender: account.Phone,
			})

			if err != nil {
				return err
			}

			fmt.Printf("Sending message %d.\n-> Sender: %s\n-> Recipient: %s\n-> Message: %s\n\n",
				message.ID, account.Phone, message.ReceiverPhone, message.Message)

			message.SetSent(db)
		}

		messages = []MessageModel{}
	}

	return nil
}
