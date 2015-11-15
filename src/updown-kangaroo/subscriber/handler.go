package subscriber

import (
	"log"
	"time"

	"golang.org/x/net/websocket"
)

// Handler handles communication.
type Handler struct {
	Subscribers map[*Subscriber]bool
	forward     chan *Message
	subscribe   chan *Subscriber
	unsubscribe chan *Subscriber
	send        chan *Message
}

const (
	messageBufferSize = 256
)

// KangarooBroadcastServer broadcasts to subscribers.
func (h *Handler) KangarooBroadcastServer(ws *websocket.Conn) {

	sub := &Subscriber{
		Websocket: ws,
		Send:      make(chan *Message, messageBufferSize),
		handler:   h,
	}
	h.subscribe <- sub
	defer func() { h.unsubscribe <- sub }()
	go sub.write()
	sub.read()
}

// NewHandler creates a new subscriber handler.
func NewHandler() *Handler {
	return &Handler{
		forward:     make(chan *Message),
		subscribe:   make(chan *Subscriber),
		unsubscribe: make(chan *Subscriber),
		Subscribers: make(map[*Subscriber]bool),
		send:        make(chan *Message),
	}
}

// Broadcast message to subscribers.
func (h *Handler) Broadcast(msg *Message) {
	h.send <- msg
}

// Run should be used as a goroutine to handle subscribers.
func (h *Handler) Run() {
	for {
		select {
		case subscriber := <-h.subscribe:
			h.Subscribers[subscriber] = true
			log.Println("New subscriber.")
		case subscriber := <-h.unsubscribe:
			delete(h.Subscribers, subscriber)
			close(subscriber.Send)
			log.Println("Subscriber unsubscribed.")
		case msg := <-h.forward:
			log.Println("Reviced message: ", msg.Sender)
		case msg := <-h.send:
			// broadcast msg.
			for subscriber := range h.Subscribers {
				select {
				case subscriber.Send <- msg:
					log.Println("Time message process took: ", time.Since(msg.Recived))
				default:
					// unable to deliver msg to client.
					delete(h.Subscribers, subscriber)
					close(subscriber.Send)
					log.Println("-- unable to deliver msg to client.")
				}
			}
		}
	}
}
