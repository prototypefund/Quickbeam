package api

import (
	"fmt"
	"reflect"
)

type StateModule struct {
	Identifier string
	Function   interface{}
}

func (a *Api) RegisterState(identifier string, function interface{}) {
	state := StateModule{identifier, function}
	a.states[identifier] = state
}

type StateReturn struct {
	States []string `json:"states"`
}

func (a *Api) getState(app string) (res []string, err error) {
	state, ok := a.states[app]
	stateFunc := state.Function
	if !ok {
		return nil, AppNotAvailable{app}
	}
	argList := []reflect.Value{}
	argList, err = a.appendDependencies(argList, stateFunc, 0)
	results := reflect.ValueOf(stateFunc).Call(argList)
	if !results[1].IsNil() {
		return nil, results[1].Interface().(error)
	}
	res, ok = results[0].Interface().([]string)
	if !ok {
		return nil, InternalDispatchError{
			fmt.Sprintf("State function for app '%v' did not return a list of state descriptors but %v", app, results[0].Type())}
	}
	return
}
