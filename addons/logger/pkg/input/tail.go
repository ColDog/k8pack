package input

import (
	"context"

	"github.com/coldog/k8pack/addons/logger/pkg/message"
	"github.com/coldog/k8pack/addons/logger/pkg/tail"
)

func init() {
	Register("tail", NewTail)
}

// NewTail is the tail constructor.
func NewTail() Input { return &Tail{} }

// Tail implements the input interface for file tailing.
type Tail struct {
	Match string
}

// Start kicks off the input.
func (t *Tail) Run(ctx context.Context, feed chan message.Message) error {
	return tail.Tail(ctx, t.Match, feed)
}
