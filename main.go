package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/alecthomas/kong"
	"github.com/esiqveland/notify"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/godbus/dbus/v5"
)

var CLI struct {
	Listen string `arg env:"NOTIFY_API_LISTEN" help:"Where to listen" default:":3000"`
}

func main() {
	log.Println("Notify API ⚡")
	kong.Parse(&CLI)

	bus, err := dbus.SessionBusPrivate()
	if err != nil {
		log.Fatalf("Failed to get DB session bus: %v", err)
	}
	defer bus.Close()
	log.Println("Connected to dbus!")

	err = bus.Auth(nil)
	if err != nil {
		log.Fatalf("Failed to authenticate on session bus: %v", err)
	}

	err = bus.Hello()
	if err != nil {
		log.Fatalf("Failed to call session bus: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Notify API ⚡"))
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			App     string `json:"app"`
			Title   string `json:"title"`
			Body    string `json:"body"`
			Expiry  int32  `json:"expiry"`
			Urgency string `json:"urgency"`
		}

		body.App = "notify-api"
		body.Expiry = -1
		body.Urgency = "normal"

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			log.Printf("Failed to parse json body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		n := notify.Notification{
			AppName:       body.App,
			Summary:       body.Title,
			Body:          body.Body,
			ExpireTimeout: body.Expiry,

			Hints: map[string]dbus.Variant{
				"urgency": urgencyVariant(body.Urgency),
			},
		}

		_, err = notify.SendNotification(bus, n)
		if err != nil {
			log.Printf("Failed to call dbus: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	log.Printf("Listening on %s", CLI.Listen)
	err = http.ListenAndServe(CLI.Listen, r)
	if err != nil {
		log.Fatalf("ListenAndServe failed: %v", err)
	}
}

func urgencyVariant(u string) dbus.Variant {
	var n uint8 = 1
	if u == "low" {
		n = 0
	} else if u == "critical" {
		n = 2
	}

	return dbus.MakeVariant(n)
}
