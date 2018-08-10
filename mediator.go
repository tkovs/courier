package courier

type Mediator struct {
	Couriers map[string]*Courier
	Shutdown chan string
}

func NewMediator() *Mediator {
	mediator := new(Mediator)
	mediator.Couriers = make(map[string]*Courier)
	mediator.Shutdown = make(chan string)
	go mediator.reaper()

	return mediator
}

func (this *Mediator) reaper() {
	var identity string

	for {
		identity = <-this.Shutdown
		this.Couriers[identity] = nil
		delete(this.Couriers, identity)
	}
}

func (this *Mediator) SendMessage(job Job) error {
	var c *Courier
	var ok bool
	var err error

	c, ok = this.Couriers[job.Sender]
	if !ok {
		// TODO: Check the error message
		c, err = NewCourier(job.Sender, this.Shutdown)
		if err != nil {
			return err
		}
		this.Couriers[job.Sender] = c
	}

	c.Messages <- job.Message
	return nil
}
