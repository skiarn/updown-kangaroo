package subscriber

import (
	"time"

	"golang.org/x/net/websocket"
)

// Subscriber client.
type Subscriber struct {
	Websocket *websocket.Conn
	Send      chan *Message

	handler *Handler
}

func (s *Subscriber) read() {
	for {
		var msg Message
		buf := make([]byte, 1024)
		if nr, er := s.Websocket.Read(buf); er == nil {
			msg.Data = buf[0:nr]
			msg.Sender = s
			msg.Recived = time.Now()
			s.handler.forward <- &msg
		} else {
			break
		}
	}
	s.Websocket.Close()
}

func (s *Subscriber) write() {
	for msg := range s.Send {
		if _, ew := s.Websocket.Write(msg.Data); ew != nil {
			break
		}
	}
	s.Websocket.Close()
}
