package msg

import (
	"fmt"
	"log"
	"reflect"
)

type UnregisteredMessageError struct {
	MessageName string
}

func (u *UnregisteredMessageError) Error() string {
	return fmt.Sprintf("Message with name %s is not registered", u.MessageName)
}

type HandlerRegistry map[string]reflect.Value

func (registry *HandlerRegistry) RegisterHandler(name string, fn interface{}) {
	fv := reflect.ValueOf(fn)
	fntype := fv.Type()

	log.Printf("Registering %v as \"%v\"\n", fntype, name)

	// check if method already registered
	if f, ok := (*registry)[name]; ok {
		log.Panicf("Handler %v already registered as %v", name, f)
	}

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
	(*registry)[name] = fv
}

type Handler struct {
	ParamPtr reflect.Value
	FuncVal  reflect.Value
}

func (registry *HandlerRegistry) Handler(name string) (handler *Handler, err error) {
	fv, ok := (*registry)[name]
	if !ok {
		err = &UnregisteredMessageError{name}
		return
	}

	handler = new(Handler)
	handler.ParamPtr = reflect.New(fv.Type().In(0).Elem())
	handler.FuncVal = fv

	return
}

func (h *Handler) Param() interface{} {
	return h.ParamPtr.Interface()
}

func (h *Handler) Call() (resp interface{}, err error) {
	in := []reflect.Value{h.ParamPtr}
	out := h.FuncVal.Call(in)

	ierr := out[1].Interface()
	if ierr != nil {
		err = ierr.(error)
		return
	}

	resp = out[0].Interface()
	return
}
