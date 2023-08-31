package api

import (
	"errors"
	"fmt"
	"reflect"

	"git.sr.ht/~michl/quickbeam/internal/api/builtin"
	"git.sr.ht/~michl/quickbeam/internal/protocol"
	"git.sr.ht/~michl/quickbeam/internal/web"
)

var (
	callbackTypeCollectionChange = "collection_change"
)

type EmptyArgs struct{}

type Dispatchable func(a Api, p interface{}) (interface{}, error)

func argumentType(f interface{}, num int) reflect.Type {
	return reflect.ValueOf(f).Type().In(num)
}

func assertDispatchable(f interface{}) error {
	funcVal := reflect.ValueOf(f)
	if funcVal.Kind() != reflect.Func {
		return InternalDispatchError{"dispatchee is not a function"}
	}
	funcType := funcVal.Type()
	if funcType.NumIn() < 1 {
		return InternalDispatchError{"dispatchee does not accept at least one argument"}
	}
	if funcType.In(0).Kind() != reflect.Struct {
		return InternalDispatchError{"dispatchee does not accept struct as first parameter"}
	}
	if funcType.NumOut() != 2 {
		return InternalDispatchError{"dispatchee does not return exactly two values"}
	}
	if funcType.Out(0).Kind() != reflect.Struct {
		return InternalDispatchError{"dispatchee's first return value is not a struct"}
	}
	errorInterface := reflect.TypeOf((*error)(nil)).Elem()
	if !funcType.Out(1).Implements(errorInterface) {
		return InternalDispatchError{"dispatchee's second return value is not an error"}
	}
	return nil
}

func newArgument(t reflect.Type, values map[string]interface{}) (res reflect.Value, err error) {
	res = reflect.New(t)
	dest := res.Elem()
	for k, v := range values {
		field := dest.FieldByName(k)
		if !field.IsValid() {
			return res, InvalidArgumentError{t, k}
		}
		if !field.CanSet() {
			return res, InternalDispatchError{fmt.Sprintf("cannot set field `%s`", k)}
		}
		val := reflect.ValueOf(v)
		if val.Type() != field.Type() {
			return res, InvalidTypeError{t, k, val.Type()}
		}
		field.Set(val)
	}
	return res, nil
}

func numberOfArguments(f interface{}) int {
	v := reflect.ValueOf(f)
	t := v.Type()
	return t.NumIn()
}

func (a *Api) appendDependencies(argList []reflect.Value, function interface{}, startingAt int) ([]reflect.Value, error) {
	webpageType := reflect.TypeOf((*web.Page)(nil)).Elem()
	for i := startingAt; i < numberOfArguments(function); i++ {
		t := argumentType(function, i)
		switch t {
		case reflect.TypeOf(a):
			argList = append(argList, reflect.ValueOf(a))
		case webpageType:
			argList = append(argList, reflect.ValueOf(a.WebPage))
		default:
			return nil, InternalDispatchError{fmt.Sprintf("Unknown type for dependency injection: %v", t)}
		}
	}
	return argList, nil
}

func (a *Api) dispatchFunc(f interface{}, arguments map[string]interface{}) (res interface{}, err error) {
	if err := assertDispatchable(f); err != nil {
		return nil, err
	}

	argType := argumentType(f, 0)
	argPointer, err := newArgument(argType, arguments)
	if err != nil {
		return nil, err
	}

	argList := []reflect.Value{argPointer.Elem()}
	argList, err = a.appendDependencies(argList, f, 1)
	if err != nil {
		return nil, err
	}
	results := reflect.ValueOf(f).Call(argList)
	res = results[0].Interface()
	if !results[1].IsNil() {
		return nil, results[1].Interface().(error)
	}
	return
}

func paramMissingError(param string) error {
	return ParamMissingError{param}
}

func paramWrongTypeError(param string, value interface{}, expectedType string) error {
	actualType := reflect.ValueOf(value).Type()
	return errors.New(fmt.Sprintf("Parameter '%s' is of type %v, not %s", param, actualType, expectedType))
}

type Api struct {
	WebPage     web.Page
	collections map[string]Collection
	actions     map[string]Action
	states      map[string]StateModule
	nextId      int
	CallBack    func(string, interface{})
}

func New() Api {
	a := Api{
		CallBack: func(_ string, _ interface{}) {
			return
		},
		collections: make(map[string]Collection),
		actions:     make(map[string]Action),
		states:      make(map[string]StateModule),
	}
	a.RegisterAction("greet", greet)
	a.RegisterCollection(tickingCollection)
	return a
}

func (a *Api) Dispatch(method string, args DispatchArgs) (result interface{}, err error) {
	switch method {
	case "ping":
		return builtin.Ping()
	case "wait":
		duration, err := getArgInt(args, "for")
		if err != nil {
			return nil, err
		}
		res := builtin.Wait(duration)
		return res, nil
	case "version":
		return builtin.GetVersion(), nil
	case "call":
		action, err := getArgString(args, "action")
		if err != nil {
			return nil, err
		}
		actionParams, ok := args["args"]
		if !ok {
			return nil, errors.New("Invalid call: args missing")
		}
		apMap, ok := actionParams.(map[string]interface{})
		ap := DispatchArgs(apMap)
		if !ok && actionParams != nil {
			return nil, errors.New("Invalid call: action params not a struct")
		}
		return a.dispatchAction(action, ap)
	case "open":
		url, err := getArgString(args, "url")
		if err != nil {
			return nil, err
		}
		return nil, a.WebPage.Navigate(url)
	case "state":
		application, err := getArgString(args, "application")
		if err != nil {
			return nil, err
		}
		return a.getState(application)
	case "fetch":
		collection, err := getArgString(args, "collection")
		if err != nil {
			return nil, err
		}
		return a.fetchCollection(collection)
	case "subscribe":
		collection, err := getArgString(args, "collection")
		if err != nil {
			return nil, err
		}
		return a.subscribeCollection(collection)
	}
	return nil, errors.New("Unknown Method")
}

func getArgString(args DispatchArgs, key string) (string, error) {
	v, ok := args[key]
	if !ok {
		return "", protocol.UserError("parameter %s is missing", key)
	}
	value, ok := v.(string)
	if !ok {
		return "", protocol.UserError("parameter %s is not a string", key)
	}
	return value, nil
}

func getArgInt(args DispatchArgs, key string) (int, error) {
	value, ok := args[key]
	if !ok {
		return 0, paramMissingError(key)
	}
	f, ok := value.(float64)
	if !ok {
		return 0, paramWrongTypeError(key, value, "int")
	}
	res := int(f)
	if float64(res) != f {
		return 0, paramWrongTypeError(key, value, "int")
	}
	return res, nil
}

type DispatchArgs map[string]interface{}
