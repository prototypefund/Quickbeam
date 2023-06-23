package api

import "testing"

func actionArgs(name string, args DispatchArgs) DispatchArgs {
	return DispatchArgs{
		"action": name,
		"args": args,
	}
}

func TestDispatch(t *testing.T) {
	tests := []struct{
		name string
		method string
		args DispatchArgs
		result interface{}
		err bool
	}{
		{"ping should pong", "ping", nil, "pong", false},
		{"unknown method", "foo", nil, nil, true},
		{"call unknown action", "call",
			actionArgs("foo", make(DispatchArgs)), nil, true},
	}
	for _, tt := range(tests) {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Dispatch(tt.method, tt.args)
			if tt.err && err == nil {
				t.Error("Exptected error, got none")
			}
			if !tt.err && err != nil {
				t.Errorf("Expected success, got error '%v'", err)
			}
		})
	}
}
