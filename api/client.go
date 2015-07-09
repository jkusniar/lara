package api

import (
	"fmt"
	"log"
	"reflect"
)

type method struct {
	RequestType  reflect.Type
	ResponseType reflect.Type
	HandlerFn    reflect.Value
}

type registry map[string]*method

func (reg registry) RegisterMethod(methodName string, reqType reflect.Type, respType reflect.Type,
	fn reflect.Value) {
	log.Printf("Registering %v with request type %v and response type %v as \"%v\"\n",
		fn.Type(), reqType, respType, methodName)

	// check if method already registered
	if m, ok := reg[methodName]; ok {
		log.Panicf("Method %v already registered as %v", methodName, m.HandlerFn)
	}

	fntype := fn.Type()

	// check if registering a function
	if fntype.Kind() != reflect.Func {
		log.Panicf("`fn` should be %s but is %s", reflect.Func, fntype.Kind())
	}

	// check if function has proper in/out count
	if fntype.NumIn() != 1 {
		log.Panicf("`fn` should have 1 parameter but it has %d parameters", fntype.NumIn())
	}
	if fntype.NumOut() != 2 {
		log.Panicf("`fn` should return 2 values but it returns %d values", fntype.NumOut())
	}

	// check in/out data types
	if fntype.In(0).Kind() != reflect.Ptr {
		log.Panicf("Parameter of `fn` should be %s but is %s", reflect.Ptr, fntype.In(0).Kind())
	}
	if fntype.In(0).Elem().Kind() != reflect.Struct {
		log.Panicf("Parameter of `fn` should point to %s but is pointing to %s",
			reflect.Struct, fntype.In(0).Elem().Kind())
	}

	if fntype.Out(0).Kind() != reflect.Ptr {
		log.Panicf("1st response value of `fn` should be %s but is %s", reflect.Ptr,
			fntype.Out(0).Kind())
	}
	if fntype.Out(0).Elem().Kind() != reflect.Struct {
		log.Panicf("1st response value of `fn` should point to %s but is pointing to %s",
			reflect.Struct, fntype.Out(0).Elem().Kind())
	}

	//check if second response argument is error
	errorType := reflect.TypeOf((*error)(nil)).Elem()
	if fntype.Out(1) != errorType {
		log.Panicf("2nd response value of `fn` should implement error interface but is %s",
			fntype.Out(1))
	}

	//register
	reg[methodName] = &method{reqType, respType, fn}
}

var (
	Registry registry
)

func init() {
	Registry = make(map[string]*method)

	Registry.RegisterMethod("CreateClient",
		reflect.TypeOf(CreateClientRequest{}),
		reflect.TypeOf(CreateClientResponse{}),
		reflect.ValueOf(CreateClient))
}

type CreateClientRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (r CreateClientRequest) String() string {
	return fmt.Sprint("firstName: ", r.FirstName, ", lastName: ", r.LastName)
}

type CreateClientResponse struct {
	Id string `json:"id"`
}

func CreateClient(req *CreateClientRequest) (*CreateClientResponse, error) {
	log.Println("Calling CreateClient ", req.FirstName, ", ", req.LastName)
	return &CreateClientResponse{"f47ac10b-58cc-4372-a567-0e02b2c3d479"}, nil
}
