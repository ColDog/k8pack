package message

import (
	"time"
)

// Message is a message emitted for each log line.
type Message struct {
	Tag       string
	Data      map[string]string
	Message   string
	Timestamp time.Time
	Priority  string
	Hostname  string
	Program   string
}
