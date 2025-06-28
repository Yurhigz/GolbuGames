package main

import (
	"fmt"
	"golbugames/backend/internal/websocket/ws"
	"net/http"
	"time"
)

const (
	smallPayloadHeaderLen  = 6
	mediumPayloadHeaderLen = 8
	largePayloadHeaderLen  = 14
)

func main() {
	HubManager := ws.NewHubManager()
	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		ws.WebsocketHandler(w, r, HubManager)
	})
	fmt.Printf("Listening on port %v ...", 3005)

	http.ListenAndServe(":3005", nil)
	time.Sleep(100 * time.Millisecond)

}
