package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubscribe(t *testing.T) {
	collChanged := make(chan interface{})
	type callbackParams struct {
		Method string
		Params interface{}
	}
	callbacks := []callbackParams{}
	getAll := func(_ EmptyArgs) (res CollectionGetFunctionResult, err error) {
		res.Members = []interface{}{}
		return
	}
	subscribe := func(_ EmptyArgs) (res CollectionSubsribeFunctionResult, err error) {
		res.Channel = collChanged
		return
	}
	callback := func(method string, params interface{}) {
		callbacks = append(callbacks, callbackParams{method, params})
	}
	a := New()
	a.CallBack = callback
	collection := Collection{
		Identifier:    "test",
		GetAllMembers: getAll,
		Subscribe:     subscribe,
	}
	a.RegisterCollection(collection)
	args := DispatchArgs{
		"collection": "test",
	}
	resInterface, err := a.Dispatch("subscribe", args)
	assert.Nilf(t, err, "Got error: %s", err)
	subscribeResult, ok := resInterface.(SubscribeCollectionResult)
	assert.True(t, ok, "Wrong return type for 'subscribe': %T", resInterface)
	collChanged <- "change"
	assert.Len(t, callbacks, 1, "Wrong number of callback invocations")
	expectedCallback := callbackParams{
		Method: callbackTypeCollectionChange,
		Params: CollectionChangedParams{"test", subscribeResult.Id},
	}
	assert.Equal(t, expectedCallback, callbacks[0])
}
