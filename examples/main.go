package main

import (
	"time"

	"../src/courier"
)

func main() {
	var c courier.Courier

	c = courier.Courier{Identity: "558299542550"}
	c.GetReady()

	c.Messages <- courier.Message{"Opa!", "558299406148"}
	c.Messages <- courier.Message{"Verificamos uma pendência de R$500,00 no sistem que está no seu nome.", "558299406148"}
	c.Messages <- courier.Message{"Podemos negociar o pagamento em até 3x sem juros no cartão.", "558299406148"}
	c.Messages <- courier.Message{"Aguardamos contato\n, Plussoft.", "558299406148"}

	time.Sleep(50 * time.Second)
}
