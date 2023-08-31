package api

import (
	"fmt"

	"git.sr.ht/~michl/quickbeam/internal/web"
)

type Action struct {
	Identifier string
	Function   interface{}
}

func (a *Api) RegisterAction(identifier string, function interface{}) {
	action := Action{identifier, function}
	a.actions[identifier] = action
}

type greetParams struct {
	Name string `json:"name"`
}

type greetReturn struct {
	Greeting string
}

func greet(p greetParams, w web.Page, a *Api) (greetReturn, error) {
	name := p.Name
	greeting := fmt.Sprintf("Hello, %s!", name)
	return greetReturn{greeting}, nil
}

func (a *Api) dispatchAction(action string, params map[string]interface{}) (result interface{}, err error) {
	act, ok := a.actions[action]
	if !ok {
		return nil, ActionNotAvailableError{action}
	}
	return a.dispatchFunc(act.Function, params)
}
