package output

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/coldog/k8pack/addons/logger/pkg/message"
)

var outputs = map[string]Const{}

// Register an output.
func Register(name string, output Const) {
	outputs[name] = output
}

// Const represents a constructor for an output.
type Const func() Output

// Spec represnts a generic spec for an output.
type Spec struct {
	Type   string
	Output json.RawMessage
}

// GetOutput returns the Input interface for the given type and spec.
func (s Spec) GetOutput() (Output, error) {
	construct, ok := outputs[s.Type]
	if !ok {
		return nil, fmt.Errorf("input: no input named: %s", s.Type)
	}
	output := construct()
	if s.Output != nil {
		err := json.Unmarshal(s.Output, output)
		if err != nil {
			return nil, err
		}
	}
	return output, nil
}

// Output is a log output that reads messages from the feed channel. Output
// should exit when the message channel is closed.
type Output interface {
	Run(context.Context, chan message.Message) error
}
