package filter

import (
	"encoding/json"
	"fmt"

	"github.com/coldog/k8pack/addons/logger/pkg/message"
)

var filters = map[string]Const{}

// Register a filter.
func Register(name string, filter Const) {
	filters[name] = filter
}

// Const represents a constructor for a filter.
type Const func() Filter

// Spec represnts a generic spec for a filter.
type Spec struct {
	Type   string
	Filter json.RawMessage
}

// GetFilter fetches a filter based on the type name from the global registry.
func (s Spec) GetFilter() (Filter, error) {
	construct, ok := filters[s.Type]
	if !ok {
		return nil, fmt.Errorf("filter: no filter named: %s", s.Type)
	}
	filt := construct()
	if s.Filter != nil {
		err := json.Unmarshal(s.Filter, filt)
		if err != nil {
			return nil, err
		}
	}
	return filt, nil
}

// Filter is the main filter interface.
type Filter interface {
	Filter(msg message.Message) message.Message
}
