package api

import (
	"errors"
	"fmt"
	"reflect"

	"git.sr.ht/~michl/quickbeam/internal/api/builtin"
)

var (
	actions map[string]interface{} = map[string]interface{}{
		"greet": greet,
	}
)

type greetParams struct {
	Name string `json:"name"`
}

func greet(p greetParams) (string, error) {
	name := p.Name
	return fmt.Sprintf("Hello, %s!", name), nil
}

type ArgumentError struct{
	message string
}

func (e ArgumentError) Error() string {
	return e.message
}

type ActionError struct{
	message string
}

func (e ActionError) Error() string {
	return e.message
}

type InternalDispatchError struct{
	message string
}

func (e InternalDispatchError) Error() string {
	return e.message
}

func dispatchAction(action string, params map[string]interface{}) (result interface{}, err error) {
	act, ok := actions[action]
	if !ok {
		return nil, ActionError{"Action not available"}
	}
	return dispatchFunc(act, params)
}

type Dispatchable func(p interface{}) (interface{}, error)

// getArgumentType returns the type of f's argument. It also checks whether f is dispatchable. This means, it accepts one single argument that is a struct and returns two arguments, a return value and an error. It returns an error if not.
func getArgumentType(f interface{}) reflect.Type {
	return reflect.ValueOf(f).Type().In(0)
}

func assertDispatchable(f interface{}) error {
	funcVal := reflect.ValueOf(f)
	if funcVal.Kind() != reflect.Func {
		return InternalDispatchError{"dispatchee is not a function"}
	}
	funcType := funcVal.Type()
	if funcType.NumIn() != 1 {
		return InternalDispatchError{"dispatchee does not accept exactly one argument"}
	}
	if funcType.In(0).Kind() != reflect.Struct {
		return InternalDispatchError{"dispatchee parameter is not a struct"}
	}
	if funcType.NumOut() != 2 {
		return InternalDispatchError{"dispatchee does not return exactly two values"}
	}
	if funcType.Out(0).Kind() != reflect.Struct {
		return InternalDispatchError{"dispatchee's first return value is not a struct"}
	}
	errorInterface := reflect.TypeOf((*error)(nil)).Elem()
	if !funcType.Out(1).Implements(errorInterface) {
		return InternalDispatchError{"dispatchee's second return value is not an error"}
	}
	return nil
}

func newArgument(t reflect.Type, values map[string]interface{}) (res reflect.Value, err error) {
	res = reflect.New(t)
	dest := res.Elem()
	for k, v := range values {
		field := dest.FieldByName(k)
		if !field.IsValid() {
			return res, InvalidArgumentError{t, k}
		}
		if !field.CanSet() {
			return res, InternalDispatchError{fmt.Sprintf("cannot set field `%s`", k)}
		}
		val := reflect.ValueOf(v)
		if val.Type() != field.Type() {
			return res, InvalidTypeError{t, k, val.Type()}
		}
		field.Set(val)
	}
	return res, nil
}

func dispatchFunc(f interface{}, arguments map[string]interface{}) (res interface {}, err error) {
	if err := assertDispatchable(f); err != nil {
		return nil, err
	}
	argType := getArgumentType(f)
	argPointer, err := newArgument(argType, arguments)
	if err != nil {
		return nil, err
	}

	results := reflect.ValueOf(f).Call([]reflect.Value{argPointer.Elem()})
	res = results[0].Interface()
	if !results[1].IsNil() {
		return nil, results[1].Interface().(error)
	}
	return
}

func Dispatch(method string, params map[string]interface{}) (result interface{}, err error) {
	switch method {
	case "ping":
		return builtin.Ping()
	case "version":
		return builtin.GetVersion(), nil
	case "call":
		a, ok := params["action"]
		if !ok {
			return nil, errors.New("Invalid call: action missing")
		}
		action, ok := a.(string)
		if !ok {
			return nil, errors.New("Invalid call: action not a string")
		}
		actionParams, ok := params["params"]
		if !ok {
			return nil, errors.New("Invalid call: params missing")
		}
		ap, ok := actionParams.(map[string]interface{})
		if !ok {
			return nil, errors.New("Invalid call: action params not a struct")
		}
		return dispatchAction(action, ap)
	case "webpage.open":
		return dispatchFunc(openWebpage, params)
	}
	return nil, errors.New("Unknown Method")
}
