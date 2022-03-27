package server

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"net/http"
)

// uploadRouter forwards requests to /upload to the appropriate handlers
func uploadRouter(r chi.Router) {
	// Send post requests to uploadFile function
	r.Post("/", uploadFile)
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type FileResponse struct {
	Name  string `json:"name"`
	Size  int    `json:"size"`
	Token string `json:"token"`
}

// uploadFile accepts http request containing a multipart file form
func uploadFile(writer http.ResponseWriter, request *http.Request) {
	// Attempt to read the file from the request
	_, m, err := request.FormFile("file")
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, err = writer.Write([]byte("File not uploaded correctly, try again"))
		if err != nil {
			return
		}
		return
	}

	// Do OpenCV, get Token
	open, err := m.Open()
	if err != nil {
		return
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(open)
	if err != nil {
		return
	}

	process, err := Process(buf.Bytes())
	if err != nil {
		return
	}
	// Populate the file response struct
	fileResponse := FileResponse{
		Name:  m.Filename,
		Size:  int(m.Size),
		Token: process,
	}
	response := Response{
		Success: true,
		Message: "success",
		Data:    fileResponse,
	}
	// Convert the struct to json
	marshal, err := json.Marshal(response)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, err = writer.Write([]byte("Failed to form JSON"))
		if err != nil {
			return
		}
		return
	}
	// Send the response to the requester
	_, err = writer.Write(marshal)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func Process(buffer []byte) (string, error) {
	// Read the buffer in
	img, err := gocv.IMDecode(buffer, gocv.IMReadAnyColor)
	// Close the image when the function exits
	defer img.Close()
	if err != nil {
		fmt.Printf("Error opening image buffer\n")
		return "", nil
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
	encoded, err := gocv.IMEncode(".jpg", imgOut)
	if err != nil {
		return "", err
	}
	// Close the image when the function exits
	defer encoded.Close()
	// Allocate a response buffer
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(encoded.GetBytes())))
	// Encode the result matrix to the user
	base64.StdEncoding.Encode(buf, encoded.GetBytes())
	return string(buf), nil

}
