package orchestrator

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
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

const thumbnailResolution = 360
const resultResolution = 2560

type LeafProcess struct {
	key    uuid.UUID
	buffer []byte
	result string
	points []interface{}
}

func (l *LeafProcess) Assign(u uuid.UUID) {
	l.key = u
}

// NewLeafProcessJob creates a request for a session token
func NewLeafProcessJob(buffer []byte) uuid.UUID {
	process := LeafProcess{
		buffer: buffer,
		key:    uuid.New(),
	}
	id, err := meta.Enroll(&process)
	if err != nil {
		return id
	}
	return id
}

// Key returns the identifier key of the process
func (l *LeafProcess) Key() uuid.UUID {
	return l.key
}

func NewUpdate(state string, message string, data interface{}) Update {
	return Update{
		Time:    time.Now(),
		State:   state,
		Message: message,
		Data:    data,
	}
}

// Run starts the process
func (l *LeafProcess) Run(c chan Update) error {
	// Send an update to the client confirming upload
	c <- NewUpdate("uploaded", "Processing image for analysis", nil)
	if err := Process(l.buffer, c); err != nil {
		return err
	}
	return nil
}

type Detection struct {
	Bounds      image.Rectangle   `json:"bounds"`
	Confidence  float64           `json:"confidence"`
	Boxes       []image.Rectangle `json:"boxes"`
	Confidences []float64         `json:"confidences"`
	Type        string            `json:"type"`
}

// aggregateBoxes combines duplicated Boxes, and Boxes mostly overlapping
func aggregateBoxes(boxes []image.Rectangle, confidences []float64) ([]Detection, error) {
	var candidates []Detection
	// Iterate through all the detections
	for i, box := range boxes {
		assimilated := false
		boxAvg := float64(box.Dx()+box.Dy()) / 2.0
		// Iterate through existing candidates to see if any of them share a similar disposition
		for _, candidate := range candidates {
			candidateAvg := float64(candidate.Bounds.Dx()+candidate.Bounds.Dy()) / 2.0
			// Calculate the threshold for how far the box should be before it is grouped
			groupingThreshold := ((boxAvg + candidateAvg) / 2.0) / 2.0
			if distance(candidate.Bounds, box) <= groupingThreshold {
				// Create a new box containing both Boxes
				candidate.Bounds = candidate.Bounds.Union(box)
				// Add the current box to the subarray within the detection
				candidate.Boxes = append(candidate.Boxes, box)
				// Similarly, mark the Confidences
				candidate.Confidences = append(candidate.Confidences, confidences[i])
				// Since this box was assimilated into an existing candidate,
				// we don't want to add it to the candidates again.
				assimilated = true
				// Break out of the for loop so this box isn't grouped with other candidates.
				break
			}
		}
		// Check if the box as assimilated
		if !assimilated {
			// If the box is virgin, we will initialize a detection struct for it
			detection := Detection{
				Bounds:      box,
				Boxes:       []image.Rectangle{},
				Confidences: confidences,
				Confidence:  confidences[i],
				Type:        "Poison Oak",
			}
			// Add the box to the candidates array
			candidates = append(candidates, detection)
		}
	}
	// Generate the average Confidence based on the Confidences from the above steps
	// Iterate through all candidates allocated in the above steps
	for _, candidate := range candidates {
		// Initialize a variable with the Confidence of the first detection
		sum := candidate.Confidence
		// Iterate through each child Confidence
		for _, confidence := range candidate.Confidences {
			// Add it to the Confidence sum
			sum += confidence
		}
		// Set the detection level Confidence to the average of the grouped Confidences, plus the original
		candidate.Confidence = sum / float64(len(candidate.Confidences))
	}
	return candidates, nil
}

var textColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}

// drawDetections presents the detections on the provided matrix in a human-readable format
func drawDetections(mat *gocv.Mat, detections []Detection) error {
	// Create a new pseudo-mask
	mask := gocv.NewMatWithSize(mat.Rows(), mat.Cols(), mat.Type())
	// Allocated a new matrix
	dimmed := gocv.NewMat()
	// Copy the input mat to be dimmed
	mat.CopyTo(&dimmed)
	// Adjust the gama of the mat to create a dimmed matrix
	gocv.AddWeighted(*mat, 0.5, mask, 0.5, 0.6, &dimmed)
	// The intensity of the background space blur
	const blurRadius = 20
	// Blur the background to create a sense of depth
	gocv.GaussianBlur(dimmed, &dimmed, image.Point{}, blurRadius, blurRadius, gocv.BorderReflect)
	// Draw the mask of each detection
	for _, detection := range detections {
		gocv.Rectangle(&mask, detection.Bounds, textColor, -1)
	}
	// Create a version of mat where only the detections are visible
	gocv.BitwiseAnd(*mat, mask, mat)
	// Invert the mask to now only highlight empty space
	gocv.BitwiseNot(mask, &mask)
	// Cut the darkened empty space from dimmed
	gocv.BitwiseAnd(dimmed, mask, &dimmed)
	// Combine the two matrices into mat
	gocv.BitwiseOr(*mat, dimmed, mat)

	// Iterate through each detection
	for _, detection := range detections {
		// Define a local variable for readability
		detection.Bounds = detection.Bounds.Add(image.Pt(4, 4))
		err := drawDetectionCorners(mat, detection, color.RGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 0,
		})
		if err != nil {
			return err
		}
		detection.Bounds = detection.Bounds.Sub(image.Pt(4, 4))
		err = drawDetectionCorners(mat, detection, textColor)
		if err != nil {
			return err
		}

	}

	return nil
}

// drawDetectionCorners draws the lines at the corners of a detection and a cross-hair in the middle
func drawDetectionCorners(mat *gocv.Mat, detection Detection, clr color.RGBA) error {
	// Define a local instance of the Bounds rectangle
	rect := detection.Bounds
	// Define an inset distance variable
	const insetDistance = -10
	// Create the inset rect to make the corners appear to have depth
	insetRect := rect.Inset(insetDistance)
	// Define macro variables to hold rect dimensions
	w := insetRect.Dx()
	h := insetRect.Dy()
	// Find the center point
	center := rectCenter(insetRect)
	// Define the four corners of the detection area
	corners := [4]image.Point{
		// Bottom Right
		image.Pt(1, 1),
		// Top Right
		image.Pt(1, -1),
		// Top Left
		image.Pt(-1, 1),
		// Bottom Left
		image.Pt(-1, -1),
	}
	// Define a line length
	const line = 30
	// Define the line thickness
	const lineWeight = 3
	localOffsets := [4][4]int{
		// Bottom Right
		{-line, 0, 0, -line},
		// Top Right
		{0, line, -line, 0},
		// Top Left
		{line, 0, 0, -line},
		// Bottom Left
		{0, line, line, 0},
	}
	// Draw all four corners
	for i, loc := range localOffsets {
		// Compute the point corresponding to the correct corner
		local := center.Add(image.Pt((w/2)*corners[i].X, (h/2)*corners[i].Y))
		// Draw both lines to enclose the pseudo-square
		gocv.Line(mat, local, local.Add(image.Pt(loc[0], loc[1])), clr, lineWeight)
		gocv.Line(mat, local, local.Add(image.Pt(loc[2], loc[3])), clr, lineWeight)
	}

	// Define constants for ui
	const fontScale = 1.3
	const fontThickness = 2
	// Get text metrics for dynamically rendered text
	textDimensions := gocv.GetTextSize(fmt.Sprintf("%.2f%%", detection.Confidence*100), gocv.FontHersheyDuplex, fontScale, fontThickness)
	// Draw the offset text
	gocv.PutText(mat, fmt.Sprintf("%.2f%%", detection.Confidence*100),
		rect.Min.Add(image.Pt(rect.Dx()-textDimensions.X, textDimensions.Y)).Add(image.Pt(-16, 16)),
		gocv.FontHersheyDuplex, fontScale,
		clr, fontThickness)
	// Draw a crosshair in the middle
	gocv.Line(mat, image.Pt(center.X, center.Y-line), image.Pt(center.X, center.Y+line), clr, lineWeight)
	gocv.Line(mat, image.Pt(center.X-line, center.Y), image.Pt(center.X+line, center.Y), clr, lineWeight)
	return nil
}

// Process accepts a buffer and returns a processed buffer
func Process(buffer []byte, client chan Update) error {
	// Read the buffer in
	img, err := gocv.IMDecode(buffer, gocv.IMReadAnyColor)
	// Close the image when the function exits
	if err != nil {
		fmt.Printf("Error opening image buffer\n")
		return nil
	}
	// Defer the closing of the image
	defer func(img *gocv.Mat) {
		err = img.Close()
		if err != nil {
			fmt.Printf("Error closing image buffer\n")
		}
	}(&img)
	// Generate a thumbnail
	thumbnail, err := createThumbnail(img)
	if err != nil {
		return err
	}
	// Send the thumbnail to the user
	client <- NewUpdate("processing", "Processed image for analysis", string(thumbnail))
	// Resize the uploaded buffer to a fixed width
	result, err := resizeMatrixByWidth(img, resultResolution)
	if err != nil {
		return err
	}
	// Run the detection model on the upload
	_, confidences, boxes := tensor.Detect(result)
	// Update the user that the model has run
	client <- NewUpdate("compiling", "", nil)
	// Aggregate all the detections into a few mains ones
	detections, err := aggregateBoxes(boxes, confidences)
	if err != nil {
		return err
	}
	// Draw the detections found previously
	err = drawDetections(&result, detections)
	if err != nil {
		return err
	}
	// Convert the buffer to base64
	bufResults := matToBase64(result)

	marshal, err := json.Marshal(detections)
	if err != nil {
		return err
	}

	// Send the final results to the user
	client <- NewUpdate("results", string(marshal), bufResults)
	return nil
}

// resizeMatrixByWidth resizes a matrix to conform to a provided maxWidth
func resizeMatrixByWidth(src gocv.Mat, maxWidth int) (gocv.Mat, error) {
	// Allocate a new matrix to hold the final thumbnail
	dest := gocv.NewMat()
	// Max Width for thumbnail
	mw := float64(maxWidth)
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

// createThumbnail generates a thumbnail for the client
func createThumbnail(src gocv.Mat) ([]byte, error) {
	// Resize the input matrix to have a max width of 360 pixels
	matrix, err := resizeMatrixByWidth(src, 360)
	if err != nil {
		return nil, err
	}
	// Convert to jpg
	buf := matToBase64(matrix)
	// Return the buffer
	return []byte(buf), nil
}

// matToBase64 converts a gocv matrix into a base64 encoded string
func matToBase64(src gocv.Mat) string {
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

// rectCenter gets the center point of an image.Rectangle
func rectCenter(r1 image.Rectangle) image.Point {
	return image.Pt(r1.Min.X+r1.Dx()/2, r1.Min.Y+r1.Dy()/2)
}

// distance gets the distance between the center of two rectangles
func distance(r1 image.Rectangle, r2 image.Rectangle) float64 {
	return math.Sqrt(math.Pow(float64(rectCenter(r1).X-rectCenter(r2).X),
		2) + math.Pow(float64(rectCenter(r1).Y-rectCenter(r2).Y), 2))
}
