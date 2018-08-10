package main

import (
	"time"

	"github.com/tkovs/courier"
)

func main() {
	mediator := courier.NewMediator()

	mediator.SendMessage(courier.Job{
		Message: courier.Message{
			Content:   "Nunca mais lhe vi. Viajou?",
			Recipient: "558299542550",
		},
		Sender: "558298212550",
	})

	time.Sleep(30 * time.Second)
}
