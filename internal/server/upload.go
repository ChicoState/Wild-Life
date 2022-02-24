package server

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// uploadRouter forwards requests to /upload to the appropriate handlers
func uploadRouter(r chi.Router) {
	// Send post requests to uploadFile function
	r.Post("/", uploadFile)
}

type FileResponse struct {
	Name string `json:"name"`
	Size int    `json:"size"`
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
	// Populate the file response struct
	response := FileResponse{
		Name: m.Filename,
		Size: int(m.Size),
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
