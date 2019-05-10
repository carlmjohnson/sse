package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const timeout = 3 * time.Second

func main() {
	addr := os.Getenv("PORT")
	if addr == "" {
		addr = "8080"
	}
	if !strings.Contains(addr, ":") {
		addr = ":" + addr
	}

	http.HandleFunc("/sse.json", func(w http.ResponseWriter, r *http.Request) {
		flush := func() { log.Print("warning: flush not implemented") }
		if f, ok := w.(http.Flusher); ok {
			flush = f.Flush
		}
		w.Header().Set("Content-Type", "text/event-stream")

		enc := json.NewEncoder(w)
		start := time.Now()
		for time.Now().Sub(start) < timeout {
			var data = struct {
				Time time.Time `json:"time"`
			}{
				Time: time.Now(),
			}
			io.WriteString(w, "data: ")
			enc.Encode(&data)
			io.WriteString(w, "\n\n")
			flush()
			time.Sleep(100 * time.Millisecond)
		}
		io.WriteString(w, "event: close\ndata:\n\n")
		flush()
		log.Printf("%q: served events for %v", r.URL.Path, time.Now().Sub(start))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
		log.Printf("%q: served index.html", r.URL.Path)
	})

	log.Printf("starting listener on %s", addr)
	http.ListenAndServe(addr, nil)
}
