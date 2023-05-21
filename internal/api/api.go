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

func dispatchFunc(f interface{}, paramsMap map[string]interface{}) (res interface {}, err error) {
	funcVal := reflect.ValueOf(f)
	if funcVal.Kind() != reflect.Func {
		return nil, InternalDispatchError{"dispatchee is not a function"}
	}
	funcType := funcVal.Type()
	if funcType.NumIn() != 1 {
		return nil, InternalDispatchError{"dispatchee does not accept exactly one argument"}
	}
	paramsType := funcType.In(0)
	if paramsKind := paramsType.Kind(); paramsKind != reflect.Struct {
		return nil, InternalDispatchError{"dispatchee parameter is not a struct"}
	}
	if funcType.NumOut() != 2 {
		return nil, InternalDispatchError{"dispatchee does not return exactly two values"}
	}
	if funcType.Out(0).Kind() != reflect.Interface {
		return nil, InternalDispatchError{"dispatchee's first return value is not an interface"}
	}
	errorInterface := reflect.TypeOf((*error)(nil)).Elem()
	if !funcType.Out(1).Implements(errorInterface) {
		return nil, InternalDispatchError{"dispatchee's second return value is not an error"}
	}

	paramsStruct := reflect.New(paramsType).Elem()
	for k, v := range paramsMap {
		field := paramsStruct.FieldByName(k)
		if !field.IsValid() {
			return nil, InvalidArgumentError{paramsType, k}
		}
		if !field.CanSet() {
			return nil, InternalDispatchError{fmt.Sprintf("cannot set field `%s`", k)}
		}
		val := reflect.ValueOf(v)
		if val.Type() != field.Type() {
			return nil, InvalidTypeError{paramsType, k, val.Type()}
		}
		field.Set(val)
	}

	results := funcVal.Call([]reflect.Value{paramsStruct})
	res = results[0].Interface()
	err = results[1].Interface().(error)
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
		if params == nil {
			return nil, ParamMissingError{}
		}
		u, ok := params["url"]
		if !ok {
			return nil, ParamMissingError{"url"}
		}
		url, ok := u.(string)
		if ok {
			return openWebpage(url)
		}
	}
	return nil, errors.New("Unknown Method")
}
