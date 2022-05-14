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

// making 2 orchestrators to test concurrent access
func TestNewOrchestrator2(t *testing.T) {
	o, err := NewOrchestrator()
	if err != nil {
		t.Error(err)
	}
	err = o.Close()
	if err != nil {
		t.Error(err)
	}
}
