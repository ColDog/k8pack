package filter

import (
	"encoding/json"

	"github.com/coldog/k8pack/addons/logger/pkg/message"
)

func init() {
	Register("json", NewJSON)
}

// NewJSON implements Const for the JSONFilter.
func NewJSON() Filter { return &JSONFilter{} }

// JSONFilter parses the message body as json.
type JSONFilter struct {
	RemoveKeys []string
}

// Filter implements JSONFilter by parsing the message as json.
func (f *JSONFilter) Filter(msg message.Message) message.Message {
	data := map[string]string{}
	err := json.Unmarshal([]byte(msg.Record["message"]), &data)
	if err != nil {
		return msg
	}
	for k, v := range data {
		msg.Record[k] = v
	}
	for _, k := range f.RemoveKeys {
		delete(msg.Record, k)
	}
	return msg
}
