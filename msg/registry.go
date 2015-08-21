package msg

import (
	"fmt"
	"reflect"
)

// Registry is a map "message name" -> reflect.Value, where Value is a
// message handler function type. There should be an instance of Registry in
// an application created by new function. Message handlers can then be
// registered using Register method.
type Registry map[string]reflect.Value

// Register registeres message handler funcion in reg. Before successful
// registration couple of validation must pass. TODO
func (reg *Registry) Register(name string, fn interface{}) error {
	fv := reflect.ValueOf(fn)
	fntype := fv.Type()

	// check if method already registered
	if f, ok := (*reg)[name]; ok {
		return fmt.Errorf("handler %v already registered as %v",
			name, f)
	}

	// check if registering a function
	if fntype.Kind() != reflect.Func {
		return fmt.Errorf("`fn` should be %s but is %s",
			reflect.Func, fntype.Kind())
	}

	// check if function has proper in/out count
	if fntype.NumIn() != 1 {
		return fmt.Errorf(
			"`fn` should have 1 parameter but it has %d parameters",
			fntype.NumIn())
	}
	if fntype.NumOut() != 2 {
		return fmt.Errorf(
			"`fn` should return 2 values but it returns %d values",
			fntype.NumOut())
	}

	// check in/out data types
	if fntype.In(0).Kind() != reflect.Ptr {
		return fmt.Errorf("parameter of `fn` should be %s but is %s",
			reflect.Ptr, fntype.In(0).Kind())
	}
	if fntype.In(0).Elem().Kind() != reflect.Struct {
		return fmt.Errorf(
			"parameter of `fn` should point to %s but is pointing to %s",
			reflect.Struct, fntype.In(0).Elem().Kind())
	}

	if fntype.Out(0).Kind() != reflect.Ptr {
		return fmt.Errorf(
			"1st response value of `fn` should be %s but is %s",
			reflect.Ptr, fntype.Out(0).Kind())
	}
	if fntype.Out(0).Elem().Kind() != reflect.Struct {
		return fmt.Errorf(
			"1st response value of `fn` should point to %s but is pointing to %s",
			reflect.Struct, fntype.Out(0).Elem().Kind())
	}

	//check if second response argument is error
	errorType := reflect.TypeOf((*error)(nil)).Elem()
	if fntype.Out(1) != errorType {
		return fmt.Errorf(
			"2nd response value of `fn` should implement error interface but is %s",
			fntype.Out(1))
	}

	//register
	(*reg)[name] = fv

	return nil
}

// message handler struct
type handler struct {
	paramPtr reflect.Value // ptr to instance of message handler's param
	funcVal  reflect.Value // message handler (func)
}

// Handler returns message handler struct for a given message name from registry
func (reg *Registry) Handler(name string) (*handler, error) {
	fv, ok := (*reg)[name]
	if !ok {
		return nil,
			fmt.Errorf("message with name %s is not registered",
				name)
	}

	return &handler{
		paramPtr: reflect.New(fv.Type().In(0).Elem()),
		funcVal:  fv}, nil
}

// returns message handler parameter's interface
func (h *handler) Param() interface{} {
	return h.paramPtr.Interface()
}

// calls message handler, returning response or error
func (h *handler) Call() (resp interface{}, err error) {
	in := []reflect.Value{h.paramPtr}
	out := h.funcVal.Call(in)

	ierr := out[1].Interface()
	if ierr != nil {
		err = ierr.(error)
		return
	}

	resp = out[0].Interface()
	return
}
