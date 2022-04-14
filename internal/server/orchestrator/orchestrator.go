package orchestrator

import (
	"fmt"
	"math/rand"
	"time"
	"wildlife/internal/log"
)

var orch *Orchestrator

type Request interface {
	Run(chan Update) error
	Key() string
}

type Orchestrator struct {
	receiver chan Request
	resolver chan Request
	resolved map[string]Request
	latest   map[string]chan Update
}

// NewOrchestrator initializes a new Orchestrator
func NewOrchestrator() error {
	// Initialize the Orchestrator struct
	o := Orchestrator{
		receiver: make(chan Request),
		resolver: make(chan Request, 16),
		resolved: map[string]Request{},
		latest:   map[string]chan Update{},
	}
	orch = &o
	// Start listening on the receiver channel
	go func() {
		err := orch.receive()
		if err != nil {
			log.Errf("Orchestrator receiver exited: %s", err.Error())
		}
	}()
	// Start listening on the resolver channel
	go func() {
		err := orch.resolve()
		if err != nil {
			log.Errf("Orchestrator resolver exited: %s", err.Error())
		}
	}()

	// Return no errors
	return nil
}

type Update struct {
	Time    time.Time   `json:"time"`
	State   string      `json:"state"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Connect(token string) (chan Update, error) {
	if orch.latest[token] == nil {
		return nil, fmt.Errorf("invalid token")
	}
	return orch.latest[token], nil
}

func randomSequence() string {
	template := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	var out string
	rand.Seed(time.Now().Unix())
	for i := 0; i < 64; i++ {
		r := rand.Intn(26)
		u := template[r]
		out += string(u)
	}
	return out
}

// update forwards an update message state to a client
func (o *Orchestrator) update(req Request, state string) {
	o.latest[req.Key()] <- Update{
		Time:  time.Now(),
		State: state,
	}
}

// closeRequest closes the connection to an active listener
func (o *Orchestrator) closeRequest(req Request) {
	close(o.latest[req.Key()])
}

// Enroll is used to initialize a connection between a client and the server
func (o *Orchestrator) Enroll(req Request) error {
	// Open a channel with a ten message buffer
	rx := make(chan Update, 10)
	// Register this channel with the request token identifier
	orch.latest[req.Key()] = rx
	// Send the request to the receiver
	o.receiver <- req
	// Notify the client of a connection
	o.update(req, "queued")
	// Return no errors
	return nil
}

// resolve handles messages sent to the resolver channel
func (o *Orchestrator) resolve() error {
	// Make sure channel is open
	if o.resolver == nil {
		return fmt.Errorf("resolver channel is nil")
	}
	// Listen for new messages sent to resolver
	for d := range o.resolver {
		go func(done Request) {
			o.resolved[done.Key()] = done
			o.update(done, "complete")
			o.closeRequest(done)
		}(d)
	}
	// return no errors
	return nil
}

// receive handles messages sent to the receiver channel
func (o *Orchestrator) receive() error {
	// Make sure channel is open
	if o.resolver == nil {
		return fmt.Errorf("receiver channel is nil")
	}
	// Listen for new messages sent to receiver
	for request := range o.receiver {
		o.update(request, "running")
		update := o.latest[request.Key()]
		err := request.Run(update)
		if err != nil {
			o.update(request, "failed")
			log.Errf("Orchestrator error: %s", err)
		}
	}
	return nil
}
