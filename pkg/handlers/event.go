package handlers

import (
	"context"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"log"
	"os"

	"github.com/itsmurugappan/cloud-events-viewer/pkg/codec"
)

var (
	avroDecoder *codec.AvroDecoder
)

func initAvroDecoder() {
	schemaURL, ok := os.LookupEnv("AVRO_SCHEMA_URL")
	if !ok {
		panic("env variable 'AVRO_SCHEMA_URL' not present")
	}
	avroDecoder = codec.NewAvroDecoder(schemaURL)
}

// CloudEventReceived handles the cloud event post
func CloudEventReceived(ctx context.Context, event cloudevents.Event) {
	eventChannel <- fmt.Sprintf("Got an Event: %s", event)
}

// AvroCloudEventReceived handles the cloud event post
func AvroCloudEventReceived(ctx context.Context, event cloudevents.Event) {
	val, err := avroDecoder.Decode(event.Data())
	if err != nil {
		log.Printf("error parsing avro message %v\n", err)
		return
	}
	eventChannel <- string(val)
}
