package tensor

import (
	"os"
	"testing"
	"wildlife/internal/log"

	"gocv.io/x/gocv"
)

func TestAssetsFolder(t *testing.T) {
	log.Testf("Assets folder\n")
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
		t.Errorf("Could not change directory\n")
		t.Fail()
		log.Resf(t)
		return
	}
	files, err := os.ReadDir("assets")
	if err != nil {
		t.Errorf("Assets folder not found\n")
		t.Fail()
		log.Resf(t)
		return
	}
	if len(files) == 0 {
		t.Errorf("Assets folder empty\n")
		t.Fail()
		log.Resf(t)
		return
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
		t.Errorf("Not all files acquired\n")
		t.Fail()
		log.Resf(t)
		return
	}
	log.Resf(t)
}

func TestBuildModelNonExistent(t *testing.T) {
	log.Testf("Non-existent model\n")
	err := BuildModel("WrongFile.onnx", false)
	if err == nil {
		t.Errorf("Building Model shouldn't accept non-existent file\n")
		t.Fail()
	}
	log.Resf(t)
}

func TestBuildModelCuda(t *testing.T) {
	log.Testf("CUDA capabilities\n")
	// environment must be set to use CUDA
	if os.Getenv("CUDA_ENABLED") != "true" {
		// not necessary but if it's enabled and
		// the device can't handle it, it should throw something
		log.Resf(t)
		return
	}
	err := BuildModel("assets", true)
	if err != nil {
		t.Fatalf("Cuda is either not setup or doesn't exist on this device\n")
	}
	log.Resf(t)
}

func TestBuildModelRandomString(t *testing.T) {
	log.Testf("Random string\n")
	err := BuildModel("Some random string", false)
	if err == nil {
		t.Errorf("Model shouldn't accept random string\n")
		t.Fail()
	}
	log.Resf(t)
}

func TestModel(t *testing.T) {
	log.Testf("Model detection\n")
	err := BuildModel("assets", true)
	if err != nil {
		t.Errorf("Model not built\n")
		t.Fail()
		log.Resf(t)
		return
	}
	img := gocv.IMRead("assets/test.jpeg", gocv.IMReadColor)
	if img.Empty() {
		t.Errorf("Image not loaded\n")
		t.Fail()
		log.Resf(t)
		return
	}
	//this image should have results
	ids, _, _ := Detect(img)
	if len(ids) == 0 {
		t.Errorf("No classes detected\n")
		t.Fail()
	}
	log.Resf(t)
}

func TestModelNegative(t *testing.T) {
	log.Testf("Negative model detection\n")
	err := BuildModel("assets", true)
	if err != nil {
		t.Errorf("Model not built\n")
		t.Fail()
		log.Resf(t)
		return
	}
	img := gocv.IMRead("assets/test2.jpeg", gocv.IMReadColor)
	if img.Empty() {
		t.Errorf("Image not loaded\n")
		t.Fail()
		log.Resf(t)
		return
	}
	//this image should have no results
	ids, _, _ := Detect(img)
	if len(ids) != 0 {
		t.Fatalf("Classes detected\n")
		t.Fail()
	}
	log.Resf(t)
}
