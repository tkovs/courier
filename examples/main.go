package main

import (
	"fmt"
	"time"

	"github.com/tkovs/courier"
)

func main() {
	var err error
	mediator := courier.NewMediator()

	err = mediator.SendMessage(courier.Job{
		Message: courier.Message{
			Content:   "Nunca mais lhe vi. Viajou?",
			Recipient: "558299542550",
		},
		Sender: "558299542550",
	})

	if err != nil {
		fmt.Println("Erro:", err.Error())
	}

	time.Sleep(30 * time.Second)
}
