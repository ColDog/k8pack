package pipeline

import (
	"context"
	"errors"

	"github.com/coldog/k8pack/addons/logger/pkg/filter"
	"github.com/coldog/k8pack/addons/logger/pkg/input"
	"github.com/coldog/k8pack/addons/logger/pkg/message"
	"github.com/coldog/k8pack/addons/logger/pkg/output"
)

var (
	ErrNoInput  = errors.New("pipeline: input not provided")
	ErrNoOutput = errors.New("pipeline: output not provided")
)

// Spec is the serializable pipeline.
type Spec struct {
	BufferSize int
	Input      input.Spec
	Filters    []filter.Spec
	Output     output.Spec
}

// GetPipeline parses the pipeline and returns it.
func (s *Spec) GetPipeline() (pipe *Pipeline, err error) {
	pipe = &Pipeline{}

	pipe.Input, err = s.Input.GetInput()
	if err != nil {
		return nil, err
	}
	pipe.Output, err = s.Output.GetOutput()
	if err != nil {
		return nil, err
	}
	for _, filter := range s.Filters {
		filt, err := filter.GetFilter()
		if err != nil {
			return nil, err
		}
		pipe.Filters = append(pipe.Filters, filt)
	}
	return pipe, nil
}

type Pipeline struct {
	BufferSize int

	Input   input.Input
	Filters []filter.Filter
	Output  output.Output
}

func (p *Pipeline) filter(in chan message.Message, out chan message.Message) {
	for msg := range in {
		for _, filter := range p.Filters {
			msg = filter.Filter(msg)
		}
		out <- msg
	}
}

func (p *Pipeline) validate() error {
	if p.Input == nil {
		return ErrNoInput
	}
	if p.Output == nil {
		return ErrNoOutput
	}
	return nil
}

// Run will run the pipeline.
func (p *Pipeline) Run(ctx context.Context) error {
	if err := p.validate(); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	inputCh := make(chan message.Message, p.BufferSize)
	nextCh := make(chan message.Message, p.BufferSize)

	go func() {
		p.Input.Run(ctx, inputCh)
		close(inputCh)
	}()

	go func() {
		p.filter(inputCh, nextCh)
		close(nextCh)
	}()

	return p.Output.Run(ctx, nextCh)
}
