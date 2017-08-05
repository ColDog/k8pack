package tail

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/coldog/k8pack/addons/logger/pkg/message"
)

const (
	pollInterval     = 300 * time.Millisecond
	globPollInterval = 1 * time.Second
)

// Tail will tail all files matched by the pattern passed in, pushing events
// into the provided feed.
func Tail(ctx context.Context, pattern string, feed chan message.Message) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	matched := map[string]bool{}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		matches, err := filepath.Glob(pattern)
		if err != nil {
			return err
		}

		for _, match := range matches {
			if !matched[match] {
				log.Printf("[INFO] tail: starting tail: %s", match)
				go tail(ctx, match, feed)
				matched[match] = true
			}
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(globPollInterval):
		}
	}
}

func tail(ctx context.Context, filename string, feed chan message.Message) error {
	f, err := os.Open(filename)
	if err != nil {
		log.Printf("[WARN] reader: error while opening file %s: %v", filename, err)
		return err
	}

	reader := bufio.NewReader(f)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("[WARN] reader: error while reading file %s: %v", filename, err)
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(pollInterval):
			}
		}

		if len(line) > 0 {
			feed <- message.Message{
				Tag: filename,
				Record: map[string]string{
					"message": strings.TrimSpace(line),
				},
			}
		}
	}
}
