package orchestrator

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewOrchestrator(t *testing.T) {
	o, err := NewOrchestrator()
	if err != nil {
		t.Error(err)
	}
	err = o.Close()
	if err != nil {
		t.Error(err)
	}
}

type TestTask struct {
	key uuid.UUID
}

func (t TestTask) Assign(uuid uuid.UUID) {
	t.key = uuid
}

func (t TestTask) Run(updates chan Update) error {
	return nil
}

func (t TestTask) Key() uuid.UUID {
	return t.key
}

// cannot run 2 orchestrators in parallel
// func TestOrchestrator_Enroll(t *testing.T) {
// 	// Create a new Orchestrator
// 	o, err := NewOrchestrator()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	// Create a test task
// 	ta := &TestTask{}
// 	_, err = o.Enroll(ta)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	enrolled := ta.Key()
// 	// Connect to the task to receive updates
// 	connect, err := Connect(enrolled)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	for {
// 		select {
// 		case msg := <-connect:

// 			if msg.State == "complete" {
// 				err = o.Close()
// 				if err != nil {
// 					t.Error(err)
// 				}
// 			}
// 		case <-time.After(time.Second * 3):
// 			t.Errorf("Orchestrator failed to resolve within 3 seconds")
// 			return
// 		default:
// 			return
// 		}
// 	}
// }
