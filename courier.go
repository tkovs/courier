package courier

import (
	"encoding/gob"
	"os"
	"time"

	whatsapp "github.com/Rhymen/go-whatsapp"
	"github.com/mitchellh/go-homedir"
)

// Courier is just the messenger
// Identity, actually, is the phone number
// Messages is a channel where the mediator inserts the messages to be send
type Courier struct {
	Identity   string
	Messages   chan Message
	Connection *whatsapp.Conn
}

func NewCourier(identity string, bye chan string) (*Courier, error) {
	var session *whatsapp.Session
	var err error
	// Creates a new Courier
	courier := new(Courier)
	courier.Messages = make(chan Message, 5)
	courier.Identity = identity

	// Regains the session
	if session, err = courier.readSession(); err != nil {
		return nil, err
	}

	courier.Connection, _ = whatsapp.NewConn(10 * time.Second)
	courier.Connection.RestoreSession(*session)

	go courier.start(bye)
	return courier, nil
}

func (this *Courier) start(bye chan<- string) {
	timeout := 10

	for this.Connection != nil {
		select {
		case message, ok := <-this.Messages:
			if ok {
				msg := whatsapp.TextMessage{
					Info: whatsapp.MessageInfo{
						RemoteJid: message.Recipient + "@s.whatsapp.net",
					},

					Text: message.Content,
				}

				this.Connection.Send(msg)

				time.Sleep(5 * time.Second)
			} else {
				this.Connection = nil
			}
		case <-time.After(time.Duration(timeout) * time.Second):
			this.Connection = nil
		}
	}

	bye <- this.Identity
}

func (this *Courier) readSession() (*whatsapp.Session, error) {
	session := new(whatsapp.Session)
	home, _ := homedir.Dir()
	file, err := os.Open(home + "/.courier/sessions/" + this.Identity + ".was")
	if err != nil {
		return nil, err
	}

	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&session)

	return session, err
}
