package tensor

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"wildlife/internal/log"

	"gocv.io/x/gocv"
)

const (
	ClassFile = "/poisonOak.txt"
	ModelFile = "/poisonOak.onnx"
)

// loadClass the class names from a txt file
// These are labels that the model is able to classify
func loadClass(dir string) (err error) {
	// open file assets/poisonOak.txt
	body, err := ioutil.ReadFile(dir)
	if err != nil {
		return
	}
	// split by new line
	ClassNames = strings.Split(string(body), "\n")
	return
}

// BuildModel builds the network model and loads the class names
func BuildModel(assets string, isCuda bool) (err error) {
	// load class names
	err = loadClass(assets + ClassFile)
	if err != nil {
		return
	}
	if os.Getenv("CUDA_ENABLED") != "true" && isCuda {
		err = errors.New("CUDA is not enabled")
	}
	// load the model
	netTemp := gocv.ReadNetFromONNX(assets + ModelFile)
	// allows the network to be used outside this module
	Net = &netTemp
	// set the network to use cuda if applicable
	if isCuda {
		log.Logf("Using CUDA")
		err = Net.SetPreferableBackend(gocv.NetBackendCUDA)
		if err != nil {
			return err
		}
		err = Net.SetPreferableTarget(gocv.NetTargetCUDA)
		if err != nil {
			return err
		}
	} else {
		log.Logf("Using CPU")
		err = Net.SetPreferableBackend(gocv.NetBackendOpenCV)
		if err != nil {
			return err
		}
		err = Net.SetPreferableTarget(gocv.NetTargetCPU)
		if err != nil {
			return err
		}
	}
	// return error if network is not loaded
	if netTemp.Empty() {
		err = errors.New("error loading network model")
		return
	}
	return
}
