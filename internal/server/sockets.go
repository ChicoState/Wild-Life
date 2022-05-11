package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
	"wildlife/internal/server/orchestrator"
)

// socketRouter forwards requests to /upload to the appropriate handlers
func socketRouter(r chi.Router) {
	// Send post requests to uploadFile function
	r.Get("/{token}", connectWebsockets)
}

// connectWebsockets accepts http request containing a multipart file form
func connectWebsockets(writer http.ResponseWriter, request *http.Request) {
	// Initialize an error to manage returns

	// Convert the basic GET request into a WebSocket session
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	var conn *websocket.Conn
	// Upgrade the https session to a web socket session
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		http.Error(writer, "failed to initiate websocket session", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	token := chi.URLParam(request, "token")
	uuidParse, err := uuid.Parse(token)
	if err != nil {
		http.Error(writer, "no token was provided", http.StatusBadRequest)
		return
	}
	orch := request.Context().Value("orch").(*orchestrator.Orchestrator)
	rx, err := orch.Connect(uuidParse)
	if err != nil {
		http.Error(writer, "could not find active job with that token", http.StatusNotFound)
		return
	}

	for update := range rx {
		err = conn.WriteJSON(update)
		if err != nil {
			continue
		}
	}

	err = conn.Close()
	if err != nil {
		return
	}

}
