package output

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/coldog/k8pack/addons/logger/pkg/message"
)

const (
	FluentDTimeFormat = "Jan 02 15:04:05"
)

func init() {
	Register("syslog", NewSyslog)
}

// NewSyslog is the tcp constructor.
func NewSyslog() Output { return &Syslog{} }

// Syslog implements the output interface to push to tcp.
type Syslog struct {
	Addr       string
	Format     string
	TimeFormat string
}

// Run kicks off the output.
func (t *Syslog) Run(ctx context.Context, feed chan message.Message) error {
	var writer func(w io.Writer, msg message.Message) (int, error)

	time.RFC3339

	switch t.Format {
	case "rfc5424":
		writer = t.writeSyslog5454
	case "rfc3164":
		writer = t.writeSyslog3164
	default:
		return fmt.Errorf("syslog: format unknown: %s", t.Format)
	}

	conn, err := dial("udp", t.Addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	for msg := range feed {
		_, err := writer(conn, msg)
		if err != nil {
			log.Printf("[WARN] syslog: failed to write: %v", err)
		}
	}
	return nil
}

func (t *Syslog) writeSyslog5454(w io.Writer, msg message.Message) (int, error) {
	return fmt.Fprintf(w, "<%s>1 %s %s %s %s - %s\n",
		msg.Priority,
		msg.Timestamp.Format(t.TimeFormat),
		msg.Hostname,
		msg.Tag,
		msg.Program,
		msg.Message,
	)
}

func (t *Syslog) writeSyslog3164(w io.Writer, msg message.Message) (int, error) {
	return fmt.Fprintf(w, "<%s>%s %s %s[%s]: %s\n",
		msg.Priority,
		msg.Timestamp.Format(t.TimeFormat),
		msg.Hostname,
		msg.Tag,
		msg.Program,
		msg.Message,
	)
}
