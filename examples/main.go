package main

import (
	// "github.com/tkovs/courier"

	"fmt"

	courier ".."
	raven "github.com/getsentry/raven-go"
)

func init() {
	raven.SetDSN("https://db0b31ef93d043ceaac9aa6c96eb9be7:4f945b825be0415498a85def9e55a125@sentry.io/1263532")
}

func main() {
	m := courier.Monitor{}
	err := m.Start()

	if err != nil {
		fmt.Errorf(err.Error())
	}
}
