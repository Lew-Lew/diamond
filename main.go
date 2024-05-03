package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

//go:embed template.html
var templateHTML string

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(templateHTML))
	})

	http.HandleFunc("/convert", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		res := convert(r.Body)
		w.Write([]byte(res))
	})

	port := flag.Int("port", 8080, "Port to run the server on")
	flag.Parse()

	slog.Info("Server is running", "addr", fmt.Sprintf("http://127.0.0.1:%v", *port))
	err := http.ListenAndServe(fmt.Sprintf(":%v", *port), nil)
	if err != nil {
		slog.Info("Server stopped", "error", err.Error())
		os.Exit(1)
	}
}
