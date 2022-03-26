package server

import (
	"encoding/json"
	"net/http"
	"wildlife/internal/log"

	"github.com/go-chi/chi/v5"
)

// uploadRouter forwards requests to /upload to the appropriate handlers
func uploadRouter(r chi.Router) {
	// Send post requests to uploadFile function
	r.Post("/", uploadFile)
}

type FileResponse struct {
	Name string `json:"name"`
	Size int    `json:"size"`
	Type string `json:"type"`
	Size int    `json:"size"`
	Data []byte `json:"data"`
	Status int64 `json:"status"`
	Plant string `json:"plant"`
	Confidence float64 `json:"confidence"`
}

// uploadFile accepts http request containing a multipart file form
func uploadFile(writer http.ResponseWriter, request *http.Request) {
	// Attempt to read the file from the request
	data, m, err := request.FormFile("file")
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, err = writer.Write([]byte("File not uploaded correctly, try again"))
		log.Errf("File not uploaded correctly, try again -> %s", err)
		if err != nil {
			return
		}
		return
	}

	// Turn data into byte array
	dataBytes := make([]byte, m.Size)
	_, err = data.Read(dataBytes)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, err = writer.Write([]byte("Error reading file"))
		log.Errf("Error reading file -> %s", err)
		if err != nil {
			return
		}
		return
	}

	// Populate the file response struct
	response := FileResponse{
		Name: m.Filename,
		Type: m.Header.Get("Content-Type"),
		Size: int(m.Size),
	}

	// OpenCV process here
	// ...

	// Check if filetype is jpeg or png
	if response.Type != "image/jpeg" && response.Type != "image/png" {
		writer.WriteHeader(http.StatusBadRequest)
		_, err = writer.Write([]byte("File type not supported, try again"))
		log.Errf("File type not supported, try again -> %s", err)
		if err != nil {
			return
		}
		return
	}

	// Convert the struct to json
	marshal, err := json.Marshal(response)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, err = writer.Write([]byte("Failed to form JSON"))
		log.Errf("Failed to form JSON -> %s", err)
		if err != nil {
			return
		}
		return
	}
	// Send the response to the requester
	_, err = writer.Write(marshal)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Errf("Failed to send response -> %s", err)
		return
	}
	log.Logf("File uploaded successfully")
}
