package protocol

import (
	"fmt"
	"path"
	"runtime"
)

func InternalError(msg string, args ...interface{}) error {
	return QuickbeamError{
		Msg:  fmt.Sprintf(msg, args...),
		Caller: caller(0),
		Code: 10000,
	}
}

func CallerInternalError(msg string, args ...interface{}) error {
	return QuickbeamError{
		Msg:  fmt.Sprintf(msg, args...),
		Caller: caller(1),
		Code: 10000,
	}
}

func ConfigurationError(msg string, args ...interface{}) error {
	return QuickbeamError{
		Msg:  fmt.Sprintf(msg, args...),
		Caller: caller(0),
		Code: 20000,
	}
}

func EnvironmentError(msg string, args ...interface{}) error {
	return QuickbeamError{
		Msg:  fmt.Sprintf(msg, args...),
		Caller: caller(0),
		Code: 30000,
	}
}

func CallerEnvironmentError(msg string, args ...interface{}) error {
	return QuickbeamError{
		Msg:  fmt.Sprintf(msg, args...),
		Caller: caller(1),
		Code: 30000,
	}
}

func WebpageError(msg string, args ...interface{}) error {
	return QuickbeamError{
		Msg:  fmt.Sprintf(msg, args...),
		Caller: caller(0),
		Code: 40000,
	}
}

func CallerWebpageError(msg string, args ...interface{}) error {
	return QuickbeamError{
		Msg:  fmt.Sprintf(msg, args...),
		Caller: caller(1),
		Code: 40000,
	}
}

func UserError(msg string, args ...interface{}) error {
	return QuickbeamError{
		Msg:  fmt.Sprintf(msg, args...),
		Caller: caller(0),
		Code: 50000,
	}
}

type QuickbeamError struct {
	Msg  string
	Caller string
	Code uint16
}

func (e QuickbeamError) Error() string {
	functionName := path.Base(e.Caller)
	return fmt.Sprintf("%s: %s", functionName, e.Msg)
}

func caller(skip int) string {
	counter, _, _, ok := runtime.Caller(2 + skip)
	var caller string
	if !ok {
		caller = "unknown"
	} else {
		caller = runtime.FuncForPC(counter).Name()
	}
	return caller
}
