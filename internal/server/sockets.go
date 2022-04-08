package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"net/http"
	"wildlife/internal/log"
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
		log.Errf("Socket Error: %s", err)
		return
	}
	defer conn.Close()

	token := chi.URLParam(request, "token")

	rx, err := orchestrator.Connect(token)
	if err != nil {
		fmt.Println(err)
		return
	}

	for update := range rx {
		fmt.Printf("[rx] %s\n", update.State)
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
