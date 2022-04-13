package orchestrator

import (
	"encoding/base64"
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"math"
	"time"
	"wildlife/internal/log"
	"wildlife/internal/server/tensor"
)

func init() {
	log.Logf("Using GoCV v%s", gocv.Version())
	log.Logf("Using OpenCV v%s", gocv.OpenCVVersion())
}

type LeafProcess struct {
	key    string
	buffer []byte
	result string
	points []interface{}
}

// NewLeafProcessJob creates a request for a session token
func NewLeafProcessJob(buffer []byte) string {
	key := randomSequence()
	process := LeafProcess{
		key:    key,
		buffer: buffer,
	}
	err := orch.enroll(&process)
	if err != nil {
		return key
	}
	return key
}

func (l *LeafProcess) Key() string {
	return l.key
}

func (l *LeafProcess) Run(c chan Update) error {
	// Process the image with openCV
	c <- Update{
		Time:    time.Now(),
		State:   "uploaded",
		Message: "Processing image for analysis",
		Data:    "",
	}
	if err := Process(l.buffer, c); err != nil {
		return err
	}

	// DO AI STUFF HERE

	return nil
}

// Generate a thumbnail for later reference
func createThumbnail(src gocv.Mat) ([]byte, error) {
	// Allocate a new matrix to hold the final thumbnail
	dest := gocv.NewMat()
	// Max Width for thumbnail
	mw := 360.0
	// Width
	w := src.Cols()
	// Height
	h := src.Rows()
	// Aspect ratio
	as := float64(h) / float64(w)
	// New Width
	nw := math.Min(mw, float64(w))
	// New Height
	nh := nw * as
	// Resize the source matrix to
	gocv.Resize(src, &dest, image.Pt(int(nw), int(nh)), 0, 0,
		gocv.InterpolationDefault)
	// Convert to jpg
	buf := MatToBase64(dest)
	// Return the buffer
	return []byte(buf), nil
}

// Generate a thumbnail for later reference
func createResult(src gocv.Mat) (gocv.Mat, error) {
	// Allocate a new matrix to hold the final thumbnail
	dest := gocv.NewMat()
	// Max Width for thumbnail
	mw := 2560.0
	// Width
	w := src.Cols()
	// Height
	h := src.Rows()
	// Aspect ratio
	as := float64(h) / float64(w)
	// New Width
	nw := mw
	// New Height
	nh := nw * as
	// Resize the source matrix to
	gocv.Resize(src, &dest, image.Pt(int(nw), int(nh)), 0, 0,
		gocv.InterpolationDefault)
	// Return the buffer
	return dest, nil
}

func rectCenter(r1 image.Rectangle) image.Point {
	return image.Pt(r1.Min.X+r1.Dx()/2, r1.Min.Y+r1.Dy()/2)
}

func distance(r1 image.Rectangle, r2 image.Rectangle) float64 {

	return math.Sqrt(math.Pow(float64(rectCenter(r1).X-rectCenter(r2).X),
		2) + math.Pow(float64(rectCenter(r1).Y-rectCenter(r2).Y), 2))
}

// Process accepts a buffer and returns a processed buffer
func Process(buffer []byte, c chan Update) error {
	// Read the buffer in
	img, err := gocv.IMDecode(buffer, gocv.IMReadAnyColor)
	// Close the image when the function exits
	defer img.Close()
	if err != nil {
		fmt.Printf("Error opening image buffer\n")
		return nil
	}
	// Generate a thumbnail
	thumbnail, err := createThumbnail(img)
	if err != nil {
		return err
	}
	// Send the thumbnail to the user
	c <- Update{
		Time:    time.Now(),
		State:   "processing",
		Message: "Processing image for analysis",
		Data:    string(thumbnail),
	}

	c <- Update{
		Time:    time.Now(),
		State:   "analyzing",
		Message: "",
		Data:    "",
	}

	result, err := createResult(img)
	if err != nil {
		return err
	}

	ids, confidences, boxes := tensor.Detect(result)
	// _max := math.Max(float64(img.Cols()), float64(img.Rows()))
	// scale := _max / 800
	confidence := 0.0
	tmp := gocv.NewMat()
	img.CopyTo(&tmp)

	c <- Update{
		Time:    time.Now(),
		State:   "compiling",
		Message: "",
		Data:    "",
	}

	var primaries []image.Rectangle
	counts := []int{0}

	for i := 0; i < len(ids); i++ {

		if confidences[i] > confidence {
			confidence = confidences[i]
		}
	}

	for i := 0; i < len(ids); i++ {

		center := image.Pt(boxes[i].Min.X+boxes[i].Dx()/2, boxes[i].Min.Y+boxes[i].Dy()/2)

		gocv.Circle(&result, center, 16, color.RGBA{R: 255, G: 255, B: 255, A: 255}, 16)

		accounted := false

		for k := range primaries {
			if distance(primaries[k], boxes[i]) < (float64(boxes[i].Dx()+boxes[i].Dy())/2)/2 {
				primaries[k] = primaries[k].Union(boxes[i])
				counts[k] += 1
				accounted = true
			}
		}

		if !accounted {
			primaries = append(primaries, boxes[i])
			counts = append(counts, 0)
		}

	}

	for i := range primaries {

		gocv.Rectangle(&result, primaries[i].Inset(-10), color.RGBA{R: 108, G: 194, B: 2, A: 255}, 6)
		gocv.Rectangle(&result, primaries[i], color.RGBA{R: 135, G: 242, B: 3, A: 255}, 4)
		gocv.PutText(&result, fmt.Sprintf("%s", "posion oak"), primaries[i].Min.Sub(image.Pt(0, -60)),
			gocv.FontHersheySimplex, 2,
			color.RGBA{R: 135, G: 242,
				B: 3, A: 255}, 12)
	}

	bufResults := MatToBase64(result)
	// Send the thumbnail to the user
	c <- Update{
		Time:    time.Now(),
		State:   "results",
		Message: fmt.Sprintf("%.2f", confidence*100),
		Data:    bufResults,
	}
	// Draw the contours (Black background with green outlines)

	return nil

}

// MatToBase64 converts a gocv matrix into a base64 encoded string
func MatToBase64(src gocv.Mat) string {
	// Encode the matrix into a jpg
	encoded, err := gocv.IMEncode(".jpg", src)
	if err != nil {
		return ""
	}
	// Close the image when the function exits
	defer encoded.Close()
	// Allocate a response buffer
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(encoded.GetBytes())))
	// Encode the result matrix to the user
	base64.StdEncoding.Encode(buf, encoded.GetBytes())
	return string(buf)
}
