package main

import (
	"context"
	"embed"
	"flag"
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

	port *int

	initialClickCount *int
)

func init() {
	port = flag.Int("port", 8080, "The port to run the application on.")
	initialClickCount = flag.Int("initialClickCount", 0, "The initial click count when restarting the server.")
	flag.Parse()
}

func main() {
	var err error
	applicationState := api.NewApplicationState(*initialClickCount)

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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err = indexTemplate.Execute(w, nil)
		if err != nil {
			log.Fatalf("error during index template execute: %v", err)
		}
	})

	// Add any API routes -------------------------------------------------------------------------

	http.HandleFunc("/api/click", func(w http.ResponseWriter, r *http.Request) {
		applicationState.AddClick()
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/api/globalClickSSE", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		eventChannel := make(chan int)
		_, cancel := context.WithCancel(r.Context())
		defer cancel()

		go func() {
			for globalClickCount := range eventChannel {
				fmt.Fprintf(w, "data: %v\n\n", globalClickCount)
				w.(http.Flusher).Flush()
			}
		}()

		for {
			select {
			case <-r.Context().Done():
				return
			default:
				eventChannel <- applicationState.GetClicks()
				time.Sleep(1 * time.Second)
			}
		}
	})

	// Start server -------------------------------------------------------------------------------

	log.Printf("Serving template at http://localhost:%v/", *port)
	err = http.ListenAndServe(fmt.Sprintf(":%v", *port), nil)
	if err != nil {
		log.Fatalf("error during http serving: %v", err)
	}
}
