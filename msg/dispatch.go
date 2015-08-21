package msg

import (
	"encoding/json"
	"fmt"
	"io"
)

type message struct {
	Name    string          `json:"name"`
	Content json.RawMessage `json:"content"`
}

func (m *message) String() string {
	return fmt.Sprintf("message{Name=%s}", m.Name)
}

type handlerFn func(io.Writer) error

// Dispatcher is a message dispatcher object. It holds message Registry and
// on calling Dispatch, it returns a handler function for a given message.
type Dispatcher struct {
	reg *Registry
}

// NewDispatcher creates new Dispatcher instance with initialised Registry
func NewDispatcher(reg *Registry) *Dispatcher {
	return &Dispatcher{reg}
}

// TODO translate errors to JSON
func (d *Dispatcher) Dispatch(r io.ReadCloser) (handlerFn, error) {
	var m message
	err := json.NewDecoder(r).Decode(&m)
	if err != nil {
		return nil, err
	}

	// get callable message handler based on message name
	var h *handler
	h, err = d.reg.Handler(m.Name)
	if err != nil {
		return nil, err
	}

	// JSON decode handler's param value from message content
	param := h.Param()
	if err = json.Unmarshal(m.Content, &param); err != nil {
		return nil, err
	}

	return handlerFn(func(w io.Writer) error {
		// call handler function
		resp, e := h.Call()
		if e != nil {
			return e
		}

		// JSON encode response to Writer
		if e = json.NewEncoder(w).Encode(resp); e != nil {
			// TODO what to do? Panic or return error?
			return e
		}

		return nil
	}), nil
}
