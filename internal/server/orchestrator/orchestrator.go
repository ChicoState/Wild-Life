package orchestrator

import (
	"math/rand"
	"time"
	"wildlife/internal/log"
)

var orch *Orchestrator

func init() {
	receiver := make(chan Request, 8)
	resolver := make(chan Request, 16)
	orch = &Orchestrator{
		receiver: receiver,
		resolver: resolver,
		resolved: map[string]Request{},
		latest:   map[string]chan Update{},
	}
	go func() {
		err := orch.receive()
		if err != nil {
			log.Errf("Orchestrator exited: %s", err.Error())
		}
	}()
	go func() {
		err := orch.resolve()
		if err != nil {
			log.Errf("Orchestrator exited: %s", err.Error())
		}
	}()
}

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

type Update struct {
	Time    time.Time   `json:"time"`
	State   string      `json:"state"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
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

func (o *Orchestrator) update(req Request, state string) {
	o.latest[req.Key()] <- Update{
		Time:  time.Now(),
		State: state,
	}
}

func (o *Orchestrator) closeRequest(req Request) {
	close(o.latest[req.Key()])
}

func Connect(key string) (chan Update, error) {
	return orch.latest[key], nil
}

func (o *Orchestrator) enroll(req Request) error {
	rx := make(chan Update, 8)
	orch.latest[req.Key()] = rx
	o.receiver <- req
	o.update(req, "queued")
	return nil
}

func (o *Orchestrator) complete(request Request) error {
	o.resolver <- request
	return nil
}

// resolve handles all completed requests, should be started on init
func (o *Orchestrator) resolve() error {
	for d := range o.resolver {
		go func(done Request) {
			o.resolved[done.Key()] = done
			o.update(done, "complete")
			o.closeRequest(done)
		}(d)
	}
	return nil
}

// Job requests sent to the receiver channel are handled here, should be started on init
func (o *Orchestrator) receive() error {
	for req := range o.receiver {
		go func(request Request) {
			o.update(request, "running")
			update := o.latest[request.Key()]
			err := request.Run(update)
			if err != nil {
				o.update(request, "failed")
				log.Errf("Orchestrator error: %s", err)
				return
			}
		}(req)
	}
	return nil
}
