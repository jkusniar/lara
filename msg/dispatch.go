package msg

import (
	"encoding/json"
	"fmt"
	"github.com/jkusniar/lara/api"
	"io"
	"log"
	"reflect"
)

type Message struct {
	Name    string          `json:"name"`
	Content json.RawMessage `json:"content"`
}

func (m *Message) String() string {
	return fmt.Sprintf("Message.Name=%s", m.Name)
}

type UnsupportedMessageType struct {
	MessageName string
}

func (u UnsupportedMessageType) Error() string {
	return fmt.Sprintf("Unsupported message type %s", u.MessageName)
}

type MessageHandlerFunc func(io.Writer) error

func Parse(request io.ReadCloser) (fn MessageHandlerFunc, err error) {
	var message Message
	err = json.NewDecoder(request).Decode(&message)
	if err != nil {
		return
	}

	log.Println("Parsed message: ", &message)

	// prepare request -> TODO refactor reflection parts to "api" package
	method := api.Registry[message.Name]
	if method == nil {
		err = UnsupportedMessageType{message.Name}
		return
	}

	reqVal := reflect.New(method.RequestType)
	req := reqVal.Interface()
	if err := json.Unmarshal(message.Content, &req); err != nil {
		return nil, err
	}
	log.Println("req: ", req)

	return MessageHandlerFunc(func(w io.Writer) error {
		// call api handler and encode response to output
		in := []reflect.Value{reqVal}
		out := method.HandlerFn.Call(in)

		ierr := out[1].Interface()
		if ierr != nil {
			return ierr.(error)
		}

		resp := out[0].Interface()
		//w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		//w.WriteHeader(http.StatusOK)
		if e := json.NewEncoder(w).Encode(resp); e != nil {
			log.Panicln(e)
		}

		return nil
	}), nil
}
