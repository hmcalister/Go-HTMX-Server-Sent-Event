package main

import (
	"context"
	"embed"
	"fmt"
	"hmcalister/HTMXServerSentEvent/api"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"time"
)

var (
	//go:embed static/css/output.css
	embedCSSFile []byte

	//go:embed static/htmx/htmx.js
	embedHTMXFile []byte

	//go:embed static/htmx/sse.js
	embedSSEFile []byte

	//go:embed static/templates/*.html
	templatesFS embed.FS
)

func main() {
	var err error
	applicationState := api.NewApplicationState()

	// Parse templates from embedded file system --------------------------------------------------

	templatesFS, err := fs.Sub(templatesFS, "static/templates")
	if err != nil {
		log.Fatalf("error during embedded file system: %v", err)
	}
	indexTemplate, err := template.ParseFS(templatesFS, "index.html")
	if err != nil {
		log.Fatalf("error parsing template: %v", err)
	}

	// Add handlers for CSS and HTMX files --------------------------------------------------------

	http.HandleFunc("/css/output.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		w.Write(embedCSSFile)
	})

	http.HandleFunc("/htmx/htmx.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/javascript")
		w.Write(embedHTMXFile)
	})

	http.HandleFunc("/htmx/sse.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/javascript")
		w.Write(embedSSEFile)
	})

	// Add handlers for base routes, e.g. initial page --------------------------------------------

	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		err = indexTemplate.Execute(w, nil)
		if err != nil {
			log.Fatalf("error during index template execute: %v", err)
		}
	})

	// Add any API routes -------------------------------------------------------------------------
	http.HandleFunc("/api/addItem", func(w http.ResponseWriter, r *http.Request) {
		cardData := struct {
			ItemID int
		}{
			ItemID: applicationState.AddItem(),
		}
		err = cardTemplate.Execute(w, cardData)
		if err != nil {
			log.Fatalf("error during card template execute: %v", err)
		}
	})

	http.HandleFunc("/api/deleteItem", func(w http.ResponseWriter, r *http.Request) {
		applicationState.DeleteItem()
		w.Write(nil)
	})

	http.HandleFunc("/api/deleteAll", func(w http.ResponseWriter, r *http.Request) {
		applicationState.DeleteAll()
		w.Write(nil)
	})

	// Start server -------------------------------------------------------------------------------

	log.Printf("Serving template at http://localhost:8080/index")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("error during http serving: %v", err)
	}
}
