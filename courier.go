package courier

import (
	"fmt"
	"time"
)

// Courier is just the messenger
type Courier struct {
	CourierID string
	Messages  chan string
}

func (this *Courier) GetReady() {
	// prepare, if need be, and begin
	go this.GetStart()
}

func (this *Courier) GetStart() {
	var timeout int

	this.Messages = make(chan string, 5)
	timeout = 5

	for {
		select {
		case message, ok := <-this.Messages:
			if ok {
				fmt.Printf("{Courier: %s, Message: %s}\n", this.CourierID, message)
			} else {
				// self destruct
			}
		case <-time.After(time.Duration(timeout) * time.Second):
			fmt.Printf("{Courier: %s, Status: %s}\n", this.CourierID, "My job is done.")
			break
		}
	}

	// self destruct
}
