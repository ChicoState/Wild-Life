package tensor

import (
	"errors"
	"io/ioutil"
	"strings"
	"wildlife/internal/log"
	"gocv.io/x/gocv"
)

// loadClass the class names from a txt file
// These are labels that the model is able to classify
func loadClass() (err error) {
	// open file assets/poisonOak.txt
	body, err := ioutil.ReadFile("assets/poisonOak.txt")
	if err != nil {
		return
	}
	// split by new line
	ClassNames = strings.Split(string(body), "\n")
	return
}

// BuildModel builds the network model and loads the class names
func BuildModel(netName string, isCuda bool) (err error) {
	// load class names
	err = loadClass()
	if err != nil {
		return
	}
	// load the model
	netTemp := gocv.ReadNetFromONNX(netName)
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


