package output

import (
	"context"
	"fmt"

	"github.com/coldog/k8pack/addons/logger/pkg/message"
)

func init() {
	Register("stdout", NewStdout)
}

// NewStdout is the stdout constructor.
func NewStdout() Output { return &Stdout{} }

// Stdout implements the output interface to push to stdout.
type Stdout struct{}

// Run kicks off the output.
func (t *Stdout) Run(ctx context.Context, feed chan message.Message) error {
	for msg := range feed {
		fmt.Printf("%+v\n", msg)
	}
	return nil
}
