package msg

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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

	log.Println("Dispatching handler for message: ", &message)

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

	log.Println("Handlers' param: ", param)

	return MessageHandlerFunc(func(w io.Writer) error {
		// call handler function
		resp, e := handler.Call()
		if e != nil {
			return e
		}

		log.Println("Handlers' result: ", resp)

		// JSON encode response to Writer
		// TODO response content type should be application/json
		// /w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if e = json.NewEncoder(w).Encode(resp); e != nil {
			log.Panicln(e)
		}

		return nil
	}), nil
}
