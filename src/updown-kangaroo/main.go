package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
	"time"
	"updown-kangaroo/subscriber"
	"updown-kangaroo/templates"

	"golang.org/x/net/websocket"
)

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse()

	http.Handle("/", &templates.TemplateHandler{Filename: "index.html"})
	clientHandler := subscriber.NewHandler()
	go clientHandler.Run()
	http.Handle("/out", websocket.Handler(clientHandler.KangarooBroadcastServer))

	http.HandleFunc("/in", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Recived in request:", r)
		type Data struct {
			Host   string `json:"host"`
			Path   string `json:"path"`
			Status string `json:"status"`
		}
		data, err := json.Marshal(Data{Host: r.URL.Host, Path: html.EscapeString(r.URL.Path)})
		if err != nil {
			log.Println("Failed to marshal data of request.", r)
		}
		msg := subscriber.Message{Data: data, Recived: time.Now()}
		clientHandler.Broadcast(&msg)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "")
	})

	// start the web server
	log.Println("Starting the web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
