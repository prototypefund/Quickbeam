package api

import "errors"

func Dispatch(method string, params interface{}) (result interface{}, err error) {
	switch method {
	case "ping":
		return "pong", nil
	case "version":
		return "0.1", nil
	}
	return nil, errors.New("Unknown Method")
}
