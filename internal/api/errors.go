package api

import (
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
