package handlers

import (
	"context"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// CloudEventReceived handles the cloud event post
func CloudEventReceived(ctx context.Context, event cloudevents.Event) {
	eventChannel <- fmt.Sprintf("Got an Event: %s", event)
}
