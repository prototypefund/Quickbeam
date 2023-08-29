package api

import (
	"strconv"
	"time"

	"git.sr.ht/~michl/quickbeam/internal/protocol"
)

type Collection struct {
	Identifier    string
	GetAllMembers interface{}
	Subscribe     interface{}
}

type CollectionGetFunctionResult struct {
	Members []interface{} `json:"members"`
}

type CollectionSubsribeFunctionResult struct {
	Channel <-chan interface{} `json:"channel"`
}

func (a *Api) fetchCollection(collection string) (res CollectionGetFunctionResult, err error) {
	col, err := a.getCollection(collection)
	if err != nil {
		return
	}
	getRes, err := a.dispatchFunc(col.GetAllMembers, map[string]interface{}{})
	if err != nil {
		return
	}
	res, ok := getRes.(CollectionGetFunctionResult)
	if !ok {
		return res, protocol.InternalError("Wrong return type of GetAllMembers function of collection %s", collection)
	}
	return
}

type SubscribeCollectionResult struct {
	Id string `json:"id"`
}

func (a *Api) subscribeCollection(collection string) (res SubscribeCollectionResult, err error) {
	col, err := a.getCollection(collection)
	if err != nil {
		return
	}
	sr, err := a.dispatchFunc(col.Subscribe, map[string]interface{}{})
	if err != nil {
		return
	}
	subscribeRes, ok := sr.(CollectionSubsribeFunctionResult)
	if !ok {
		return res, protocol.InternalError("Wrong return type of Subscribe function of collection %s", collection)
	}
	c := subscribeRes.Channel
	identifier := a.NextIdentifier()
	go func() {
		for {
			<-c
			message := CollectionChangedParams{
				Collection: collection,
				Id:         identifier,
			}
			a.CallBack(callbackTypeCollectionChange, message)
		}
	}()
	return SubscribeCollectionResult{
		Id: identifier,
	}, nil
}

type CollectionChangedParams struct {
	Collection string `json:"collection"`
	Id         string `json:"id"`
}

func (a *Api) getCollection(identifier string) (Collection, error) {
	col, ok := a.collections[identifier]
	if !ok {
		return col, protocol.UserError("Collection %s not found", identifier)
	}
	return col, nil
}

func (a *Api) RegisterCollection(coll Collection) {
	a.collections[coll.Identifier] = coll
}

func (a *Api) NextIdentifier() string {
	a.nextId += 1
	identifier := strconv.Itoa(a.nextId)
	return identifier
}

func getTickingMembers(_ EmptyArgs) (res CollectionGetFunctionResult, err error) {
	for _, tick := range ticks {
		var member interface{} = tick
		res.Members = append(res.Members, member)
	}
	return
}

func subscribeTicking(_ EmptyArgs) (res CollectionSubsribeFunctionResult, err error) {
	if ticker == nil {
		ticker = time.NewTicker(time.Second * 3)
		go func() {
			for {
				t := <-ticker.C
				times := []string{t.Format(time.RFC822)}
				l := len(ticks)
				if l > 9 {
					l = 9
				}
				ticks = append(times, ticks[:l]...)
				for _, c := range tickChannels {
					c <- t
				}
			}
		}()
	}
	c := make(chan interface{})
	tickChannels = append(tickChannels, c)
	res.Channel = c
	return
}

var ticker *time.Ticker
var ticks []string
var tickChannels []chan interface{}
var tickingCollection = Collection{
	Identifier:    "test/tick",
	GetAllMembers: getTickingMembers,
	Subscribe:     subscribeTicking,
}
