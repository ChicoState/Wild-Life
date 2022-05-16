package orchestrator

import (
	"encoding/json"
	"github.com/google/uuid"
	"gocv.io/x/gocv"
	"os"
	"testing"
	"time"
	"wildlife/internal/server/tensor"
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

func (t *TestTask) Assign(uuid uuid.UUID) {
	t.key = uuid
}

func (t *TestTask) Run(updates chan Update) error {
	updates <- Update{
		Time:    time.Time{},
		State:   "working",
		Message: "working order",
		Data:    "abc",
	}
	return nil
}

func (t *TestTask) Key() uuid.UUID {
	return t.key
}

func TestEnroll(t *testing.T) {
	o, err := NewOrchestrator()
	if err != nil {
		t.Fatal(err)
		return
	}

	ta := &TestTask{
		key: uuid.New(),
	}

	id, err := o.Enroll(ta)
	if err != nil {
		t.Fatal("Failed to enroll task")
		return
	}

	// Parse unparse the uuid to ensure it is valid
	err = id.UnmarshalText([]byte(id.String()))
	if err != nil {
		t.Fatal("Invalid UUID")
		return
	}

	connect, err := o.Connect(id)
	if err != nil {
		return
	}

	responses := 0
	for {
		select {
		case <-time.After(time.Millisecond * 100):
			t.Fatal("Orchestrator took more than 100ms to respond to single request")
		case msg := <-connect:
			if msg.State != "queued" && msg.State != "running" && msg.State != "working" {
				t.Fatal("Invalid state returned")
			} else {
				responses++
			}
			break
		}
		if responses == 3 {
			break
		}
	}
	err = o.Close()
	if err != nil {
		t.Fatal("Failed to close orchestrator")
	}
}

func TestFoliageEnroll(t *testing.T) {

	o, err := NewOrchestrator()
	if err != nil {
		t.Fatal(err)
		return
	}
	err = os.Chdir("../../../")
	if err != nil {
		t.Error(err)
	}
	err = tensor.BuildModel("assets", false)
	if err != nil {
		t.Fatal(err)
		return
	}
	img := gocv.IMRead("assets/test.jpeg", gocv.IMReadColor)
	if img.Empty() {
		t.Error(err)
	}
	encoded, err := gocv.IMEncode(".jpeg", img)
	if err != nil {
		t.Error(err)
	}

	job := NewLeafProcessJob(encoded.GetBytes())
	id, err := o.Enroll(job)
	if err != nil {
		t.Fatal("Failed to enroll task")
		return
	}

	// Parse unparse the uuid to ensure it is valid
	err = id.UnmarshalText([]byte(id.String()))
	if err != nil {
		t.Fatal("Invalid UUID")
		return
	}

	connect, err := o.Connect(id)
	if err != nil {
		t.Error(err)
		return
	}

	expected := []string{"queued", "running", "uploaded", "processing", "compiling", "results"}

	responses := 0
	for {
		select {
		case <-time.After(time.Millisecond * 2000):
			t.Fatal("Orchestrator took more than 1s to respond to single request")
		case msg := <-connect:

			responses++
			if msg.State != expected[0] {
				t.Fatal("Unexpected update state")
			}
			expected = expected[1:]
			if msg.State == "results" {
				var ud []Detection
				err = json.Unmarshal([]byte(msg.Message), &ud)

				if err != nil {
					t.Error(err)
					return
				}

				if len(ud) != 6 {
					t.Fatal("returned the incorrect number of detections")
				} else {
					err = o.Close()
					if err != nil {
						t.Fatal("Failed to close orchestrator")
					}
					return
				}
			}
			break
		}
		if responses < 1 {
			break
		}

	}

}

func TestFoliageEnroll2(t *testing.T) {

	o, err := NewOrchestrator()
	if err != nil {
		t.Fatal(err)
		return
	}

	err = tensor.BuildModel("assets", false)
	if err != nil {
		t.Fatal(err)
		return
	}
	img := gocv.IMRead("assets/test2.jpeg", gocv.IMReadColor)
	if img.Empty() {
		t.Error(err)
	}
	encoded, err := gocv.IMEncode(".jpeg", img)
	if err != nil {
		t.Error(err)
	}

	job := NewLeafProcessJob(encoded.GetBytes())
	id, err := o.Enroll(job)
	if err != nil {
		t.Fatal("Failed to enroll task")
		return
	}

	// Parse unparse the uuid to ensure it is valid
	err = id.UnmarshalText([]byte(id.String()))
	if err != nil {
		t.Fatal("Invalid UUID")
		return
	}

	connect, err := o.Connect(id)
	if err != nil {
		t.Error(err)
		return
	}

	expected := []string{"queued", "running", "uploaded", "processing", "compiling", "results"}

	responses := 0
	for {
		select {
		case <-time.After(time.Millisecond * 2000):
			t.Fatal("Orchestrator took more than 1s to respond to single request")
		case msg := <-connect:

			responses++
			if msg.State != expected[0] {
				t.Fatal("Unexpected update state")
			}
			expected = expected[1:]
			if msg.State == "results" {
				var ud []Detection
				err = json.Unmarshal([]byte(msg.Message), &ud)

				if err != nil {
					t.Error(err)
					return
				}

				if len(ud) != 0 {
					t.Fatal("returned the incorrect number of detections")
				} else {
					err = o.Close()
					if err != nil {
						t.Fatal("Failed to close orchestrator")
					}
					return
				}
			}
			break
		}
		if responses < 1 {
			break
		}

	}
	err = o.Close()
	if err != nil {
		t.Fatal("Failed to close orchestrator")
	}
}
