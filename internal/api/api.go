package api

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	Version VersionInfo = VersionInfo{
		Quickbeam: "0.1",
		API: "0.1",
	}
	actions map[string]interface{} =  map[string]interface{}{
		"greet": greet,
	}
)

type VersionInfo struct {
	Quickbeam string `json:"quickbeam"`
	API string `json:"api"`
}

type greetParams struct {
	Name string `json:"name"`
}

func greet(p greetParams) (string, error) {
	name := p.Name
	return fmt.Sprintf("Hello, %s!", name), nil
}

func dispatchAction(action string, params map[string]interface{}) (result interface{}, err error) {
	act, ok := actions[action]
	if !ok {
		return nil, errors.New("Action not available")
	}

	actVal := reflect.ValueOf(act)
	if actVal.Kind() != reflect.Func {
		return nil, errors.New("Internal problem with action (not a function)")
	}
	actType := actVal.Type()
	if actType.NumIn() != 1 {
		return nil, errors.New("Internal problem with action (invalid number of parameters)")
	}
	paramsType := actType.In(0)
	if actType.NumOut() != 2 {
		return nil, errors.New("Internal problem with action (invalid number of return values)")
	}

	paramStruct := reflect.New(paramsType).Elem()

	for k, v := range params {
		field := paramStruct.FieldByName(k)
		if !field.IsValid() {
			return nil, errors.New(fmt.Sprintf("Invalid params: %s does not exist", k))
		}
		if !field.CanSet() {
			return nil, errors.New(fmt.Sprintf("Invalid params: cannot set %s", k))
		}
		val := reflect.ValueOf(v)
		if val.Type() != field.Type() {
			return nil, errors.New(fmt.Sprintf("Invalid params: type for %s does not match: %s != %s", k, val.Type(), field.Type()))
		}
		field.Set(val)
	}

	results := actVal.Call([]reflect.Value{paramStruct})
	resValue := results[0]
	resError := results[1]
	if resError.IsNil() {
		err = nil
	} else {
		err = resError.Interface().(error)
	}

	return resValue.Interface(), err
}

func Dispatch(method string, params interface{}) (result interface{}, err error) {
	switch method {
	case "ping":
		return "pong", nil
	case "version":
		return Version, nil
	case "call":
		p, ok := params.(map[string]interface{})
		if !ok {
			return nil, errors.New("Invalid call: params not a struct")
		}
		a, ok := p["action"]
		if !ok {
			return nil, errors.New("Invalid call: action missing")
		}
		action, ok := a.(string)
		if !ok {
			return nil, errors.New("Invalid call: action not a string")
		}
		actionParams, ok := p["params"]
		if !ok {
			return nil, errors.New("Invalid call: params missing")
		}
		ap, ok := actionParams.(map[string]interface{})
		if !ok {
			return nil, errors.New("Invalid call: action params not a struct")
		}
		return dispatchAction(action, ap)
	}
	return nil, errors.New("Unknown Method")
}
