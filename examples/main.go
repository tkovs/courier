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
			Recipient: "558296495932",
		},
		Sender: "558299542550",
	})

	mediator.SendMessage(courier.Job{
		Message: courier.Message{
			Content:   "Opa! Apenas testando essa nova ferramenta",
			Recipient: "558296495932",
		},
		Sender: "558299011113",
	})

	time.Sleep(10 * time.Second)
}
