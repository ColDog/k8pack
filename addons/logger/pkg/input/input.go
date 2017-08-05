package input

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/coldog/k8pack/addons/logger/pkg/message"
)

var inputs = map[string]Const{}

// Register an input.
func Register(name string, input Const) {
	inputs[name] = input
}

// Const represents a constructor for an input.
type Const func() Input

// Spec represnts a generic spec for an input.
type Spec struct {
	Type  string
	Input json.RawMessage
}

// GetInput returns the Input interface for the given type and spec.
func (s Spec) GetInput() (Input, error) {
	construct, ok := inputs[s.Type]
	if !ok {
		return nil, fmt.Errorf("input: no input named: %s", s.Type)
	}
	input := construct()
	if s.Input != nil {
		err := json.Unmarshal(s.Input, input)
		if err != nil {
			return nil, err
		}
	}
	return input, nil
}

// Input is a log input that pushes messages onto the provided channel.
type Input interface {
	Run(context.Context, chan message.Message) error
}
