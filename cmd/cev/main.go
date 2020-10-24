package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/itsmurugappan/cloud-events-viewer/pkg/handlers"

	cloudevents "github.com/cloudevents/sdk-go/v2"

	"golang.org/x/net/websocket"
)

const (
	static_path = "/static/"
	static_dir  = "/static"
	index       = "/index.html"
	ws          = "/ws"
	ko_path     = "KO_DATA_PATH"
)

func main() {
	// init configs
	handlers.InitHandlers()

	// Static
	http.Handle(static_path, http.StripPrefix(static_path,
		http.FileServer(http.Dir(os.Getenv(ko_path)+static_dir))))

	// UI Handlers
	http.HandleFunc(index, handlers.RootHandler)
	http.Handle(ws, websocket.Handler(handlers.WSHandler))

	ctx := context.Background()
	p, err := cloudevents.NewHTTP()
	if err != nil {
		log.Fatalf("failed to create protocol: %s", err.Error())
	}

	h, err := cloudevents.NewHTTPReceiveHandler(ctx, p, handlers.CloudEventReceived)
	if err != nil {
		log.Fatalf("failed to create handler: %s", err.Error())
	}

	http.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}))

	http.ListenAndServe(":8080", nil)
}
