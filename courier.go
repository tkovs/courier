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
	var ticker *time.Ticker
	var tick, ticks int

	this.Messages = make(chan string, 5)
	ticks = 10
	ticksInterval := 500
	tick = ticks
	ticker = time.NewTicker(time.Duration(ticksInterval) * time.Millisecond)

	for range ticker.C {
		select {
		case message, ok := <-this.Messages:
			if ok {
				fmt.Printf("{Courier: %s, Tick: %d, Message: %s}\n", this.CourierID, tick, message)
				tick = ticks
			} else {
				// self destruct
			}
		default:
			fmt.Printf("{Courier: %s, Tick: %d}\n", this.CourierID, tick)
			tick--
		}

		if tick <= 0 {
			ticker.Stop()
			break
		}
	}

	// self destruct
	fmt.Printf("{Courier: %s, Status: %s}", this.CourierID, "My job is done.")
}
