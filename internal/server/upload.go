package server

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	gocv "gocv.io/x/gocv"
	"image"
	"image/color"
	"math/rand"
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

	cascade := gocv.NewCascadeClassifier()
	cascade.Load("haarcascade_frontalface_default.xml")

	eyeRight := gocv.NewCascadeClassifier()
	eyeRight.Load("haarcascade_righteye_2splits.xml")

	eyeLeft := gocv.NewCascadeClassifier()
	eyeLeft.Load("haarcascade_lefteye_2splits.xml")

	img, err := gocv.IMDecode(buffer, gocv.IMReadAnyColor)
	if err != nil {
		fmt.Printf("Error opening image buffer\n")
		return "", nil
	}

	imgGrey := gocv.NewMat()

	gocv.CvtColor(img, &imgGrey, gocv.ColorBGRToGray)

	faces := cascade.DetectMultiScale(imgGrey)
	for _, face := range faces {
		// First Pass Face
		gocv.Ellipse(&img, face.Min.Add(face.Size().Div(2)), face.Size().Div(2), 180,
			0, 360, color.RGBA{
				R: 128,
				G: 128,
				B: 128,
				A: 128,
			}, 2)
		f := face
		// Nested faces
		faceImg := imgGrey.Region(f)
		for i := 0; i < 20; i++ {
			f = face.Inset(-(20 - i))
			faceImg = imgGrey.Region(f.Add(image.Pt(20-rand.Int()%40, 20-rand.Int()%40)))
			nestedFaces := cascade.DetectMultiScale(faceImg)
			if len(nestedFaces) < 1 {
				continue
			}
			// The first nested face
			nestedFace := nestedFaces[0]
			// The main nested face
			gocv.Ellipse(&img, f.Min.Add(nestedFace.Min).Add(nestedFace.Size().Div(2)), nestedFace.Size().Div(2), 180,
				0, 360,
				color.RGBA{
					R: uint8(i % 255),
					G: 64,
					B: 128,
					A: 255,
				}, 1)
			// Confirmed face
		}

	}

	encoded, err := gocv.IMEncode(".jpg", img)
	if err != nil {
		return "", err
	}

	buf := make([]byte, base64.StdEncoding.EncodedLen(len(encoded.GetBytes())))

	base64.StdEncoding.Encode(buf, encoded.GetBytes())

	defer img.Close()

	return string(buf), nil

}
