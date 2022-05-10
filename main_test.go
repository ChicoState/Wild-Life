package main
package greetings

import (
	"testing"
	"github.com/joho/godotenv"
	"os"
	"wildlife/internal/env"
)

//TestNAME required format for functions


// main.go 
//Test ONNX
func TestONNX(/*class, foo*/) {
	//msg, err := 
}

//Test Enviornment
// t parameter Fatalf called if func ret no error or non empty str
func TestEnv(t *testing.T) {
	err := env.Load("db_active.env")
	if err != nil {
		t.Fatalf('DB ACTIVE')
	}
}

//Test CNN model
func TestONNXVersion(t *testing.T) { 
	msg := "Program supports 0.0.0 - "
	//Correspond with "i.j.k" in VERSION
	for i := 0; i < 10; ++i {
		for j := 0; j < 10; ++j {
			for k := 0; k < 10; ++k {
				str_i := i.(string)
				str_j := j.(string)
				str_k := k.(string)
				VERSION = str_i + '.' + str_j + '.' + str_k
				err := tensor.BuildModel(OnnxModel, false)
				if err != nil {
					msg += VERSION
					t.Fatalf(msg)
					return
				}

					
			}
		}
	}
}

func TestCNNBuild(t *testing.T) {
	err = tensor.BuildModel("googaga.txt", false)
	if err != nil {
		t.Fatalf("Empty class labels")
	}
}

func TestCNNBuild(t *testing.T) {
	err = tensor.BuildModel("wrong_labels.txt", false)
	if err != nil {
		t.Fatalf("Everything will be False Positives\n")
	}
}

func TestCNNBuild(t *testing.T) {
	err = tensor.BuildModel("file_doesn't exist", false)
	if err != nil {
		t.Fatalf("File does not exist")
	}
}

//Test controller
//See more tests of the database in internal/server/controller/user.go
func TestController {
	os.Setenv("DB_ACTIVE", "false")
	err = controller.InitController()
	if err != nil {
		t.Fatalf("Database is not active. D00d.\n")
	}
}

//Test the web server
//See more tests in internal/server/server.go
func TestServerPort {
	os.Setenv("PORT", "42069")
	err = server.Start()
	if err != nil {
		t.Fatalf("Teehee\n")
	}
}
