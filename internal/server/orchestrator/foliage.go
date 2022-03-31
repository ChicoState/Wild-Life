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
	nw := math.Min(float64(w), mw)
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
		State:   "thumbnail",
		Message: "Thumbnail generated",
		Data:    string(thumbnail),
	}
	// Convert the image to Greyscale
	imgGrey := gocv.NewMat()
	// Close the image when the function exits
	defer imgGrey.Close()
	// Blur the image so we have more unified borders.
	gocv.GaussianBlur(img, &imgGrey, image.Point{}, 40, 40, gocv.BorderDefault)
	// Convert blued image to black and white
	gocv.CvtColor(imgGrey, &imgGrey, gocv.ColorBGRToGray)
	// Take the threshold
	imgThresh := gocv.NewMat()
	// Close the image when the function exits
	defer imgThresh.Close()
	gocv.Threshold(imgGrey, &imgThresh, 60, 255, gocv.ThresholdOtsu|gocv.ThresholdToZero)
	// Find contours
	pv := gocv.FindContours(imgThresh, gocv.RetrievalExternal, gocv.ChainApproxNone)
	// Prepare an image to return
	highlight := gocv.NewMat()
	img.CopyTo(&highlight)
	gocv.DrawContours(&highlight, pv, -1, color.RGBA{R: 60, B: 80, G: 180, A: 128}, 8)
	bufHighlight := MatToBase64(highlight)
	c <- Update{
		Time:    time.Now(),
		State:   "highlight",
		Message: "highlight generated",
		Data:    bufHighlight,
	}
	imgOut := gocv.NewMat()
	// Close the image when the function exits
	// Convert to black and white
	imgGreyNew := gocv.NewMat()
	gocv.CvtColor(img, &imgGreyNew, gocv.ColorBGRToGray)
	gocv.BitwiseAnd(imgGreyNew, imgThresh, &imgOut)
	// Convert back to an RGB color space
	gocv.CvtColor(imgOut, &imgOut, gocv.ColorGrayToBGR)
	// Draw the contours (Black background with green outlines)
	gocv.DrawContours(&imgOut, pv, -1, color.RGBA{R: 60, B: 80, G: 180, A: 128}, 8)
	// Encode the matrix into an image format
	gocv.Normalize(imgOut, &imgOut, 0, 255, gocv.NormMinMax)
	bufThreshold := MatToBase64(imgOut)
	// Send the thumbnail to the user
	c <- Update{
		Time:    time.Now(),
		State:   "threshold",
		Message: "threshold generated",
		Data:    bufThreshold,
	}
	return nil

}

func MatToBase64(src gocv.Mat) string {
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
