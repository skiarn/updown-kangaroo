package subscriber

import "time"

//Message represents a single message.
type Message struct {
	Data    []byte
	Sender  *Subscriber
	Recived time.Time
}
