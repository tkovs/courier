package courier

import (
	"encoding/gob"
	"fmt"
	"os"
	"time"

	whatsapp "github.com/Rhymen/go-whatsapp"
)

// Courier is just the messenger
// Identity, actually, is the phone number
// Messages is a channel where the mediator inserts the messages to be send
type Courier struct {
	Identity string
	Session  whatsapp.Session
	Messages chan Message
}

func NewCourier(identity string) (*Courier, error) {
	// Creates a new Courier
	courier := new(Courier)
	courier.Identity = identity
	courier.Messages = make(chan Message, 5)

	// Regains the session
	if err := courier.readSession(); err != nil {
		fmt.Printf("ERROR NA SESSÃO: %s", err.Error())
		return nil, err
	}

	go courier.start()
	return courier, nil
}

func (this *Courier) start() {
	var timeout int
	var wac *whatsapp.Conn

	fmt.Println(os.Stdout, "{"+this.Identity+"}")

	wac, _ = whatsapp.NewConn(10 * time.Second)
	wac.RestoreSession(this.Session)

	timeout = 60

	for {
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
				// TODO:
				// self destruct
				// return
			}
		case <-time.After(time.Duration(timeout) * time.Second):
			fmt.Printf("{Courier: %s, Status: %s}\n", this.Identity, "My job is done.")
			break
		}
	}

	// TODO:
	// self destruct
	// return
}

func (this *Courier) readSession() error {
	this.Session = whatsapp.Session{}
	file, err := os.Open("/home/tkovs/.courier/sessions/" + this.Identity + ".was")
	if err != nil {
		return err
	}

	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&this.Session)

	return err
}
