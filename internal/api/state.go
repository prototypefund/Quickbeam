package api

import (
	"fmt"
	"reflect"

	"git.sr.ht/~michl/quickbeam/internal/bbb"
)

var (
	stateModules map[string]interface{} = map[string]interface{}{
		"bbb": bbb.State,
	}
)

type StateReturn struct {
	States []string `json:"states"`
}

func (a *Api) getState(app string) (res []string, err error) {
	stateFunc, ok := stateModules[app]
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
