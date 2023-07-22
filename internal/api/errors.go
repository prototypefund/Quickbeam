package api

import (
	"errors"
	"fmt"
	"reflect"
)

type ParamMissingError struct {
	param string
}

func (e ParamMissingError) Error() string {
	if e.param == "" {
		return "Params required but missing"
	}
	return fmt.Sprintf("Required parameter `%s` missing", e.param)
}

type InvalidArgumentError struct {
	ArgumentType reflect.Type
	Key string
}

func (e InvalidArgumentError) Error() string {
	return fmt.Sprintf("Target function cannot accept argument `%s`", e.Key)
}

type InvalidTypeError struct {
	ArgumentType reflect.Type
	Key string
	ValueType reflect.Type
}

func (e InvalidTypeError) Error() string {
	field, found := e.ArgumentType.FieldByName(e.Key)
	var fieldType reflect.Type
	if found {
		fieldType = field.Type
	}
	return fmt.Sprintf("Wrong type for argument `%s`: expected %v, got %v", e.Key, fieldType, e.ValueType)
}

type AppNotAvailable struct {
	Module string
}

func (e AppNotAvailable) Error() string {
	return fmt.Sprintf("Module '%v' is not available", e.Module)
}

type ActionNotAvailableError struct {
	Action string
}

func (e ActionNotAvailableError) Error() string {
	return fmt.Sprintf("Action '%v' ist not available currently", e.Action)
}

type InternalDispatchError struct{
	message string
}

func (e InternalDispatchError) Error() string {
	return e.message
}

func RuntimeError(wrapped interface{}) error {
	return errors.New(fmt.Sprintf("Runtime error: %v", wrapped))
}
