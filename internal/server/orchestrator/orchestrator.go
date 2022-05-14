package orchestrator

// UUID: Universally Unique Identifier
// Used to identify application w/out central registration authority
// 128 bit

import (
	"fmt"
	"time"
	"wildlife/internal/log"

	"github.com/google/uuid"
)

type Update struct {
	Time    time.Time   `json:"time"`
	State   string      `json:"state"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Request interface {
	Assign(uuid.UUID)
	Run(chan Update) error
	Key() uuid.UUID
}

type Orchestrator struct {
	running  bool
	receiver chan Request
	resolver chan Request
	resolved map[uuid.UUID]Request
	latest   map[uuid.UUID]chan Update
}

// NewOrchestrator initializes a new Orchestrator
func NewOrchestrator() (meta *Orchestrator, err error) {
	if meta != nil {
		return nil, fmt.Errorf("orchestrator has already been initialized")
	}
	// Initialize the Orchestrator struct
	meta = &Orchestrator{
		receiver: make(chan Request),
		resolver: make(chan Request, 16),
		resolved: map[uuid.UUID]Request{},
		latest:   map[uuid.UUID]chan Update{},
	}
	// Begin listening on channels
	err = meta.start()
	if err != nil {
		return nil, err
	}
	return meta, nil
}

// start begins listening on the channels
func (o *Orchestrator) start() error {
	// Start listening on the receiver channel
	go func() {
		err := o.receive()
		if err != nil {
			log.Errf("orchestrator receiver exited: %s", err.Error())
		}
	}()
	// Start listening on the resolver channel
	go func() {
		err := o.resolve()
		if err != nil {
			log.Errf("orchestrator receiver exited: %s", err.Error())
		}
	}()

	return nil
}

// Close will exit the Orchestrator safely
func (o *Orchestrator) Close() error {
	if o.receiver != nil {
		close(o.receiver)
	}

	if o.resolver != nil {
		close(o.resolver)
	}
	return nil
}

// Connect accepts incoming pared connections
func (o *Orchestrator) Connect(token uuid.UUID) (chan Update, error) {
	if o.latest[token] == nil {
		return nil, fmt.Errorf("token not found")
	}
	return o.latest[token], nil
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
func (o *Orchestrator) Enroll(req Request) (uuid.UUID, error) {
	// Generate a new UUID for the incoming request
	uid := uuid.New()
	req.Assign(uid)
	// Open a channel with a ten message buffer
	rx := make(chan Update, 10)
	// Register this channel with the request token identifier
	o.latest[uid] = rx
	// Send the request to the receiver
	o.receiver <- req
	// Notify the client of a connection
	o.update(req, "queued")
	// Return no errors
	return uid, nil
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
	if o.receiver == nil {
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
