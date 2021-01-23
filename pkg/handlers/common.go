package handlers

import (
	"html/template"
	"log"
	"os"
)

const (
	ko_path = "KO_DATA_PATH"
)

var (
	// Templates for handlers
	templates *template.Template
)

// InitHandlers initializes OAuth package
func InitHandlers() {

	// Templates
	tmpls, err := template.ParseGlob(os.Getenv(ko_path) + "/templates/*.html")
	if err != nil {
		log.Fatalf("Error while parsing templates: %v", err)
	}
	templates = tmpls
	initAvroDecoder()
}
