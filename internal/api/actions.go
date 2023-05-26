package api

import "fmt"

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

func dispatchAction(action string, params map[string]interface{}) (result interface{}, err error) {
	act, ok := actions[action]
	if !ok {
		return nil, ActionNotAvailableError{action}
	}
	return dispatchFunc(act, params)
}
