package tensor

import (
	"image"
	"math"
	"wildlife/internal/log"

	"gocv.io/x/gocv"
)

// ScaleFactor defines how small and how large a detection can be,
// this value is inversely proportional to the detection size.
const ScaleFactor = 255

// ImgScale is the size of the input files required by the model
const ImgScale = 640

// DefaultOut defines the output vector of the model
const DefaultOut = 5

// MaxConf is the threshold for target detection
const MaxConf = 0.70

var (
	Net        *gocv.Net
	ClassNames []string
)

// formatYoloV5 formats the input image to be a uniform size
func formatYoloV5(img gocv.Mat) gocv.Mat {
	// Define the provided image dimensions
	height := img.Rows()
	width := img.Cols()
	// Determine the scaling ratio for the input image
	max := ImgScale / math.Max(float64(height), float64(width))
	// Initialize a new matrix to contain the new image
	scaled := gocv.NewMat()
	// Resize image by multiplying its dimensions by the scaling factor
	scaleWidth := int(max * float64(width))
	scaleHeight := int(max * float64(height))
	// Use GoCV to resize the provided image matrix to fit the scaled matrix
	gocv.Resize(img, &scaled, image.Pt(scaleWidth, scaleHeight), 0, 0, gocv.InterpolationLinear)
	// Initialize the output matrix with a square scale, abiding by ImgScale
	output := gocv.NewMatWithSize(ImgScale, ImgScale, gocv.MatTypeCV8UC3)
	// Determine the midpoint of the image
	xMidpoint := 0
	yMidpoint := 0
	// Determine the xMidpoint by determining the image orientation
	if width < height {
		xMidpoint = (ImgScale - scaleWidth) / 2
	}
	// Determine the yMidpoint by determining the image orientation
	if height < width {
		yMidpoint = (ImgScale - scaleHeight) / 2
	}
	// Get the target region
	roi := output.Region(image.Rect(xMidpoint, yMidpoint, xMidpoint+scaleWidth, yMidpoint+scaleHeight))
	// Copy image to region
	scaled.CopyTo(&roi)
	// Return value
	return roi
}

// Detect returns class ids, confidences, and bounding boxes of the image after being processed by the network
func Detect(img gocv.Mat) (ids []int, confidences []float64, boxes []image.Rectangle) {
	// FormatYoloV5 converts to provided matrix into a 640x640 matrix for the model to use
	scaled := formatYoloV5(img)
	// Close the input matrix when the function exits
	defer func(mat *gocv.Mat) {
		err := mat.Close()
		if err != nil {
			log.Errf("Failed to close matrix: %s", err.Error())
		}
	}(&scaled)
	// Convert the scaled matrix into an OpenCV 4d Image blob
	blob := gocv.BlobFromImage(scaled, 1.0/ScaleFactor, image.Pt(ImgScale, ImgScale), gocv.NewScalar(0, 0, 0, 0), false, false)
	// Close the blob when the function exits
	defer func(blob *gocv.Mat) {
		err := blob.Close()
		if err != nil {
			log.Errf("Failed to close blob: %s", err.Error())
		}
	}(&blob)
	// Emplace the blob into the network with no name
	Net.SetInput(blob, "")
	// Run the network through an iteration
	result := Net.Forward("")
	// Parse outputs by its known constant size [1][25200][6]
	// Select the size of the 3rd dimension of the result matrix
	localOutputs := result.Size()[2]
	// Get the number of outputs, excluding duplicates
	outputs := result.Total() / localOutputs
	// If there are no detections, exit the function gracefully
	if outputs < 1 {
		return
	}
	// Determine the inverse scale dimensions, the opposite of the process in func formatYoloV5
	scaleWidth := float64(img.Cols()) / float64(ImgScale)
	scaleHeight := float64(img.Rows()) / float64(ImgScale)
	// Iterate through all the detections
	for i := 0; i < outputs; i++ {
		// Read the confidence from the output vector for the section at index 'i'
		confidence := float64(result.GetFloatAt3(0, i, 4))
		// If the confidence for this detection is sufficient, process it, otherwise, save the CPU cycles
		if confidence > MaxConf {
			// Determine the number out output parameters
			localParams := localOutputs - DefaultOut
			// Store the id of the highest value parameter
			var id int64 = 0
			max := 0.0
			for j := 0; j < localParams; j++ {
				val := float64(result.GetFloatAt3(0, i, DefaultOut+j))
				if val > max {
					id = int64(j)
					max = val
				}
			}
			// Determine the dimensions of the rectangle enclosing the detection
			x := float64(result.GetFloatAt3(0, i, 0))
			y := float64(result.GetFloatAt3(0, i, 1))
			w := float64(result.GetFloatAt3(0, i, 2))
			h := float64(result.GetFloatAt3(0, i, 3))
			// Scale the values from the relative scaled values to the original value
			left := int((x - w/2) * scaleWidth)
			top := int((y - h/2) * scaleHeight)
			right := int((x + w/2) * scaleWidth)
			bottom := int((y + h/2) * scaleHeight)
			ids = append(ids, int(id))
			confidences = append(confidences, confidence)
			boxes = append(boxes, image.Rect(left, top, right, bottom))
		}
	}
	return ids, confidences, boxes
}
