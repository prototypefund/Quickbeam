package builtin

import "time"

type WaitReturn struct {
	Since string `json:"waited-since"`
	For   int    `json:"waited-for"`
}

func Wait(duration int) WaitReturn {
	now := time.Now()
	time.Sleep(time.Second * time.Duration(duration))
	res := WaitReturn{
		Since: now.Local().String(),
		For:   duration,
	}
	return res
}
