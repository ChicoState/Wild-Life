package env

import (
	"os"
	"testing"
)

func TestLoadEnv(t *testing.T) {
	err := Load()
	if err == nil {
		t.Fatalf("There should be no .env in this directory\n")
	}
}

func TestLoadEnvTrue(t *testing.T) {
	//change directory to the root of the project
	err := os.Chdir("../../")
	if err != nil {
		t.Fatalf("Could not change directory\n")
	}
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Could not get current directory\n")
	}
	t.Logf("Changed directory to %s\n", pwd)
	err = Load()
	if err != nil {
		t.Fatalf("Could not load .env file\n")
	}
}
