package orchestrator

import (
	"github.com/google/uuid"
	"testing"
)

func TestNewLeafProcessJob(t *testing.T) {
	orchestrator, err := NewOrchestrator()
	if err != nil {
		t.Error(err)
		return
	}

	job := NewLeafProcessJob([]byte{})
	enroll, err := orchestrator.Enroll(job)
	if err != nil {
		t.Error(err)
		return
	}

	connect, err := orchestrator.Connect(enroll)
	if err != nil {
		t.Error(err)
		return
	}

	if msg := <-connect; msg.State != "queued" {
		t.Fatal("Unexpected state")
	}

	if msg := <-connect; msg.State != "running" {
		t.Fatal("Unexpected state")
	}

}

func TestNewLeafProcessMutability(t *testing.T) {

	job := NewLeafProcessJob([]byte{})
	refUUID := uuid.New()
	job.Assign(refUUID)

	if refUUID != job.Key() {
		t.Fatal("Process key is not mutable")
	}
}
