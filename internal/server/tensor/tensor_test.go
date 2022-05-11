package tensor

import (
	"os"
	"testing"

	"gocv.io/x/gocv"
)

func TestAssetsFolder(t *testing.T) {
	//check files in assets folder from os
	expected := []string{
		"poisonOak.onnx",
		"poisonOak.txt",
		"test.jpeg",
		"test2.jpeg",
		"train.ipynb",
	}
	err := os.Chdir("../../../")
	if err != nil {
		t.Fatalf("Could not change directory\n")
	}
	files, err := os.ReadDir("assets")
	if err != nil {
		t.Fatalf("Assets folder not found\n")
	}
	if len(files) == 0 {
		t.Fatalf("Assets folder empty\n")
	}
	num_acquired := 0
	for _, f := range files {
		for _, e := range expected {
			if f.Name() == e {
				num_acquired++
				continue
			}
		}
	}
	if num_acquired != len(expected) {
		t.Fatalf("Not all files acquired\n")
	}
}

func TestBuildModelNonExistent(t *testing.T) {
	err := BuildModel("WrongFile.onnx", false)
	if err == nil {
		t.Fatalf("Building Model shouldn't accept non-existent file\n")
	}
}

func TestBuildModelCuda(t *testing.T) {
	// environment must be set to use CUDA
	if os.Getenv("CUDA_ENABLED") != "true" {
		// not necessary but if it's enabled and
		// the device can't handle it, it should throw something
		t.Logf("CUDA not enabled in ENV, so this should pass\n")
		return
	}
	err := BuildModel("assets", true)
	if err == nil {
		t.Fatalf("Cuda is either not setup or doesn't exist on this device\n")
	}
}

func TestBuildModelRandomString(t *testing.T) {
	err := BuildModel("Some random string", false)
	if err == nil {
		t.Fatalf("Model shouldn't accept random string\n")
	}
}

func TestModel(t *testing.T) {
	err := BuildModel("assets", true)
	if err != nil {
		t.Fatalf("Model not built\n")
	}
	img := gocv.IMRead("assets/test.jpeg", gocv.IMReadColor)
	if img.Empty() {
		t.Fatalf("Image not loaded\n")
	}
	//this image should have results
	ids, _, _ := Detect(img)
	if len(ids) == 0 {
		t.Fatalf("No classes detected\n")
	}
}

func TestModelNegative(t *testing.T) {
	err := BuildModel("assets", true)
	if err != nil {
		t.Fatalf("Model not built\n")
	}
	img := gocv.IMRead("assets/test2.jpeg", gocv.IMReadColor)
	if img.Empty() {
		t.Fatalf("Image not loaded\n")
	}
	//this image should have no results
	ids, _, _ := Detect(img)
	if len(ids) != 0 {
		t.Fatalf("Classes detected\n")
	}
}
