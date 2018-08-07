package main

import (
	"time"

	"github.com/tkovs/courier"
)

func main() {
	var c courier.Courier

	c = courier.Courier{}
	c.CourierID = "5582999542550"
	c.GetReady()

	time.Sleep(2 * time.Second)
	c.Messages <- "Hello, there!"
	time.Sleep(3 * time.Second)
	c.Messages <- "Hello, there!"
	time.Sleep(1 * time.Second)
	c.Messages <- "Byee"
	time.Sleep(10 * time.Second)
}
