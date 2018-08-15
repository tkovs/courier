package main

import (
	// "github.com/tkovs/courier"

	"fmt"

	courier ".."
)

func main() {
	m := courier.Monitor{}
	err := m.Start()

	if err != nil {
		fmt.Println("Erro:", err.Error())
	}
}
