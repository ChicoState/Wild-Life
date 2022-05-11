package server

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"wildlife/internal/server/orchestrator"
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

type FileRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Size int    `json:"size"`
	Data []byte `json:"data"`
}

// func (f *FileRequest) toBase64() (string, error) {
// 	buf := make([]byte, base64.StdEncoding.EncodedLen(len(encoded.GetBytes())))
// 	// Encode the result matrix to the user
// 	base64.StdEncoding.Encode(buf, encoded.GetBytes())
// }

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
	orch := request.Context().Value("orch").(*orchestrator.Orchestrator)

	token := orchestrator.NewLeafProcessJob(orch, buf.Bytes())
	// Populate the file response struct
	fileResponse := FileResponse{
		Name:  m.Filename,
		Size:  int(m.Size),
		Token: token.String(),
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
