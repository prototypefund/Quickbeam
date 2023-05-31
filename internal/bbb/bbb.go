package bbb

import (
	"git.sr.ht/~michl/quickbeam/internal/web"
)

type JoinArgs struct {}
type JoinResult struct {success string}

func Join(p JoinArgs, w web.WebPage) (res JoinResult, err error) {
	// m.b.Page.MustElement("body > div.ReactModalPortal > div > div > div.sc-jObWnj.fWuLOw > div > div > span > button:nth-child(1)").MustClick()
	button := w.Root().SubNode("[aria-label='Microphone']")
	button.Click()
	return JoinResult{"success"}, nil
}
