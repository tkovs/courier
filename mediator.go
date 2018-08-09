package courier

type Mediator struct {
	Couriers map[string]*Courier
	Shutdown chan string
}

func NewMediator() Mediator {
	mediator := Mediator{}
	mediator.Couriers = make(map[string]*Courier)

	return mediator
}

func (this *Mediator) SendMessage(job Job) error {
	var c *Courier
	var ok bool
	var err error

	c, ok = this.Couriers[job.Sender]
	if !ok {
		// TODO: Check the error message
		c, err = NewCourier(job.Sender)
		if err != nil {
			return err
		}
		this.Couriers[c.Identity] = c
	}

	c.Messages <- job.Message
	return nil
}
