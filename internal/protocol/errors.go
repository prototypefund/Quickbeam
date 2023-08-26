package protocol

import "fmt"

func InternalError(msg string) error {
	return QuickbeamError{
		Msg:  msg,
		Code: 10000,
	}
}

func ConfigurationError(msg string) error {
	return QuickbeamError{
		Msg:  msg,
		Code: 20000,
	}
}

func EnvironmentError(msg string, a ...interface{}) error {
	message := fmt.Sprintf(msg, a...)
	return QuickbeamError{
		Msg:  message,
		Code: 30000,
	}
}

func WebpageError(msg string) error {
	return QuickbeamError{
		Msg:  msg,
		Code: 40000,
	}
}

func UserError(msg string) error {
	return QuickbeamError{
		Msg:  msg,
		Code: 50000,
	}
}

type QuickbeamError struct {
	Msg  string
	Code uint16
}

func (e QuickbeamError) Error() string {
	return e.Msg
}
