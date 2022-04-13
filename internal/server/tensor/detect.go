package tensor

import (
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"math"
)

const ScaleFactor = 255

// ideal since this is what the model expects most of the time

const ImgScale = 640

// [1 25200 6] index 2 is the number of expected outputs

const NumOut = 3

const MaxConf = 0.82

var HighlightColor = color.RGBA{G: 255, A: 255}

var (
	Net        *gocv.Net
	ClassNames []string
)

// formats the input image to be a uniform size
func formatYoloV5(img gocv.Mat) (out gocv.Mat) {
	// get size of image
	height := img.Rows()
	width := img.Cols()
	// get max dimension
	_max := ImgScale / math.Max(float64(height), float64(width))
	// resize image to fit into IMG_SCALE x IMG_SCALE
	newImg := gocv.NewMat()
	// scale required to fit into IMG_SCALE x IMG_SCALE
	newWidth := int(_max * float64(width))
	newHeight := int(_max * float64(height))
	// resize image
	gocv.Resize(img, &newImg, image.Pt(newWidth, newHeight), 0, 0, gocv.InterpolationLinear)
	// create blank image of size IMG_SCALE
	out = gocv.NewMatWithSize(ImgScale, ImgScale, gocv.MatTypeCV8UC3)
	// set x/y points to center image
	xMidpoint := 0
	yMidpoint := 0
	if width < height {
		xMidpoint = (ImgScale - newWidth) / 2
	}
	if height < width {
		yMidpoint = (ImgScale - newHeight) / 2
	}
	// get region of interest
	roi := out.Region(image.Rect(xMidpoint, yMidpoint, xMidpoint+newWidth, yMidpoint+newHeight))
	// copy image to region of interest
	newImg.CopyTo(&roi)
	return
}

// Detect returns class ids, confidences, and bounding boxes of the image after being processed by the network
func Detect(img gocv.Mat) (ids []int64, confidences []float64, boxes []image.Rectangle) {
	input := formatYoloV5(img)
	defer input.Close()
	// input_image is turned into a blob
	blob := gocv.BlobFromImage(input, 1.0/ScaleFactor, image.Pt(ImgScale, ImgScale), gocv.NewScalar(0, 0, 0, 0), false, false)
	defer blob.Close()
	// Input is the blob
	Net.SetInput(blob, "")
	// Run netwrok
	res := Net.Forward("")
	// get number of outputs according to the shape [1 25200 6]
	// we care about index 1 of the shape
	outs := res.Total() / NumOut
	// get the scale factor required to convert the output to the original image size
	xFactor := float64(img.Cols()) / float64(ImgScale)
	yFactor := float64(img.Rows()) / float64(ImgScale)
	// loop through all outputs
	if outs > 0 {
		for i := 0; i < int(outs); i++ {
			// get the confidence of the prediction
			confidence := float64(res.GetFloatAt3(0, i, 4))
			if confidence > MaxConf {
				// turn into gocv.Mat
				classId := 0
				// pulls coords from the blob
				x := float64(res.GetFloatAt3(0, i, 0))
				y := float64(res.GetFloatAt3(0, i, 1))
				w := float64(res.GetFloatAt3(0, i, 2))
				h := float64(res.GetFloatAt3(0, i, 3))
				// scale to original image
				left := int((x - w/2) * xFactor)
				top := int((y - h/2) * yFactor)
				right := int((x + w/2) * xFactor)
				bottom := int((y + h/2) * yFactor)
				// append items to each slice
				box := image.Rect(left, top, right, bottom)
				boxes = append(boxes, box)
				ids = append(ids, int64(classId))
				confidences = append(confidences, confidence)
			}
		}
	}
	return
}
