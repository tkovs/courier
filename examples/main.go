package main

import (
	// "github.com/tkovs/courier"

	"fmt"

	courier ".."
)

func main() {
	db, err := courier.Migrate()
	if err != nil {
		fmt.Println(err)
	}

	courier.CreateAccount(db, courier.AccountModel{Phone: "558299542550"})
}
