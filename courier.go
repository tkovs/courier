package courier

import (
	"encoding/gob"
	"os"
	"time"

	whatsapp "github.com/Rhymen/go-whatsapp"
	raven "github.com/getsentry/raven-go"
	"github.com/jmoiron/sqlx"
	"github.com/mitchellh/go-homedir"
)

// Courier is just the messenger
// Identity, actually, is the phone number
// Messages is a channel where the mediator inserts the messages to be send
type Courier struct {
	Identity string
	Messages chan MessageModel
	conn     *whatsapp.Conn
	DB       *sqlx.DB
}

func NewCourier(identity string, bye chan string, db *sqlx.DB) (*Courier, error) {
	var session *whatsapp.Session
	var err error

	// Creates a new Courier
	courier := new(Courier)
	courier.Messages = make(chan MessageModel, 5)
	courier.Identity = identity
	courier.DB = db

	// Regains the session
	if session, err = courier.getSession(); err != nil {
		return nil, err
	}

	err = courier.Connect(session)
	if err != nil {
		return nil, err
	}

	go courier.start(bye)
	return courier, nil
}

func (this *Courier) Connect(session *whatsapp.Session) error {
	var err error

	this.conn, err = whatsapp.NewConn(10 * time.Second)
	this.conn.RestoreSession(*session)

	return err
}

func (this *Courier) start(bye chan<- string) {
	timeout := 10

	for this.conn != nil {
		select {
		case message, ok := <-this.Messages:
			if ok {
				msg := whatsapp.TextMessage{
					Info: whatsapp.MessageInfo{
						RemoteJid: message.ReceiverPhone + "@s.whatsapp.net",
					},

					Text: message.Message,
				}

				err := this.conn.Send(msg)
				if err != nil {
					raven.CaptureErrorAndWait(err, nil)
					message.SetError(this.DB)
					bye <- this.Identity
				}

				message.SetSent(this.DB)
				time.Sleep(5 * time.Second)
			} else {
				bye <- this.Identity
				this.conn = nil
			}
		case <-time.After(time.Duration(timeout) * time.Second):
			bye <- this.Identity
			this.conn = nil
		}
	}
}

func (this *Courier) getSession() (*whatsapp.Session, error) {
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
