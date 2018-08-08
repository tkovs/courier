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

func (this *Courier) readSession() error {
	this.Session = whatsapp.Session{}
	file, err := os.Open("~/git/courier/sessions/" + this.Identity + ".was")
	if err != nil {
		return err
	}

	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&this.Session)

	return err
}

func (this *Courier) GetReady() {
	this.Messages = make(chan Message, 5)

	if err := this.readSession(); err != nil {
		fmt.Println(os.Stderr, "Error on reading "+this.Identity+" session file!\nError description: "+err.Error())
		return
	}

	go this.start()
}

func (this *Courier) start() {
	var timeout int
	var err error
	var wac *whatsapp.Conn

	// Try open the whatsapp websocket
	wac, err = whatsapp.NewConn(10 * time.Second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating connection: %v\n", err)
		return
	}

	wac.RestoreSession(this.Session)

	timeout = 60

	for {
		select {
		case message, ok := <-this.Messages:
			if ok {
				// fmt.Printf("{Courier: %s, Message: %s}\n", this.Identity, message)
				wac.Send(whatsapp.TextMessage{
					Info: whatsapp.MessageInfo{
						RemoteJid: message.Recipient + "@s.whatsapp.net",
					},

					Text: message.Content,
				})

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
