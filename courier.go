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
	Identity string
	Session  whatsapp.Session
	Messages chan Message
	Bye      chan string
}

func NewCourier(identity string, bye chan string) (*Courier, error) {
	// Creates a new Courier
	courier := new(Courier)
	courier.Identity = identity
	courier.Messages = make(chan Message, 5)
	courier.Bye = bye

	// Regains the session
	if err := courier.readSession(); err != nil {
		return nil, err
	}

	go courier.start()
	return courier, nil
}

func (this *Courier) start() {
	var timeout int
	var wac *whatsapp.Conn

	wac, _ = whatsapp.NewConn(10 * time.Second)
	wac.RestoreSession(this.Session)

	timeout = 60

	for wac != nil {
		select {
		case message, ok := <-this.Messages:
			if ok {
				msg := whatsapp.TextMessage{
					Info: whatsapp.MessageInfo{
						RemoteJid: message.Recipient + "@s.whatsapp.net",
					},

					Text: message.Content,
				}

				wac.Send(msg)

				time.Sleep(5 * time.Second)
			} else {
				wac = nil
			}
		case <-time.After(time.Duration(timeout) * time.Second):
			wac = nil
		}
	}

	this.Bye <- this.Identity
}

func (this *Courier) readSession() error {
	this.Session = whatsapp.Session{}
	home, _ := homedir.Dir()
	file, err := os.Open(home + "/.courier/sessions/" + this.Identity + ".was")
	if err != nil {
		return err
	}

	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&this.Session)

	return err
}
