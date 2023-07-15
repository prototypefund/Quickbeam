package marionette

import (
	"testing"

	"git.sr.ht/~michl/quickbeam/internal/protocol"
)

type testExecuter struct {
	Result string
	Command string
}

func (e *testExecuter) ExecOrEmpty(command string) string {
	e.Command = command
	return e.Result
}

func TestStart(t *testing.T) {
	t.Log("TestStart")
	testCases := []struct {
		firefox *Firefox
		executer *testExecuter
		reportedCommand string
		errCode uint16
	}{
		{
			&Firefox{},
			&testExecuter{"cat", ""},
			"which firefox",
			0,
		},
		{
			&Firefox{
				FirefoxPath: "cat",
			},
			&testExecuter{"firefox", ""},
			"",
			0,
		},
	}

	for _, tc := range testCases {
		err := startFirefox(tc.firefox, tc.executer)
		if err != nil {
			if tc.errCode == 0 {
				t.Errorf("No error expected, this occured: %v", err)
			}
			qe, ok := err.(protocol.QuickbeamError)
			if !ok {
				t.Errorf("Expected error code %d, but returned error is no QuickbeamError", tc.errCode)
			}
			if ok && qe.Code != tc.errCode {
				t.Errorf("Expected error code %d, got %d", tc.errCode, qe.Code)
			}
		}
		if tc.reportedCommand != tc.executer.Command {
			t.Errorf("Expected '%s' to be called, but '%s' was", tc.reportedCommand, tc.executer.Command)
		}
		if tc.errCode > 0 {
			if err == nil {
				t.Errorf("Expected error code %d, but no error occured", tc.errCode)
			}
		}
		tc.firefox.process.Kill()
	}
}
