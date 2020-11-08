package main

import (
	"log"
	"net/http"

	"github.com/alecthomas/kong"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var CLI struct {
	Listen string `arg env:"NOTIFY_API_LISTEN" help:"Where to listen" default:":3000"`
}

func main() {
	log.Println("Notify API ⚡")
	kong.Parse(&CLI)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Notify API ⚡"))
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	log.Printf("Listening on %s", CLI.Listen)
	err := http.ListenAndServe(CLI.Listen, r)
	if err != nil {
		log.Fatalf("ListenAndServe failed: %v", err)
	}
}
