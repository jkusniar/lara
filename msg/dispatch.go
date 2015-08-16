package msg

import (
	"encoding/json"
	"fmt"
	"github.com/jkusniar/lara/app"
	"io"
)

type Message struct {
	Name    string          `json:"name"`
	Content json.RawMessage `json:"content"`
}

func (m *Message) String() string {
	return fmt.Sprintf("Message.Name=%s", m.Name)
}

type MessageHandlerFunc func(io.Writer) error

type Dispatcher struct {
	Registry *HandlerRegistry
}

// TODO translate errors to JSON
func (dispatcher *Dispatcher) Dispatch(request io.ReadCloser) (fn MessageHandlerFunc, err error) {
	var message Message
	err = json.NewDecoder(request).Decode(&message)
	if err != nil {
		return
	}

	app.Log.Debug("Dispatching handler for message: ", &message)

	// get callable message handler based on message name
	var handler *Handler
	handler, err = dispatcher.Registry.Handler(message.Name)
	if err != nil {
		return
	}

	// JSON decode handlers' param value from message content
	param := handler.Param()
	if err = json.Unmarshal(message.Content, &param); err != nil {
		return
	}

	app.Log.Debug("Handlers' param: ", param)

	return MessageHandlerFunc(func(w io.Writer) error {
		// call handler function
		resp, e := handler.Call()
		if e != nil {
			return e
		}

		app.Log.Debug("Handlers' result: ", resp)

		// JSON encode response to Writer
		if e = json.NewEncoder(w).Encode(resp); e != nil {
			app.Log.Panic(e)
		}

		return nil
	}), nil
}
