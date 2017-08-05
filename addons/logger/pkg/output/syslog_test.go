package output

import (
	"context"
	"testing"
	"time"

	"github.com/coldog/k8pack/addons/logger/pkg/message"
	"github.com/stretchr/testify/assert"
)

func TestSyslog_Run(t *testing.T) {
	feed := make(chan message.Message)

	go func() {
		for i := 0; i < 10; i++ {
			feed <- message.Message{
				Tag:       "test",
				Program:   "test",
				Message:   "testing",
				Timestamp: time.Now(),
				Priority:  "1",
			}
		}
		close(feed)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := (&Syslog{
		Addr:       "127.0.0.1:24224",
		Format:     "rfc5424",
		TimeFormat: FluentDTimeFormat,
	}).Run(ctx, feed)
	assert.Nil(t, err)
}
