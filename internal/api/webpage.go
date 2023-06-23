package api

import (
	"git.sr.ht/~michl/quickbeam/internal/web"
)

type WebpageId uint64

var (
	webpages = make(map[WebpageId] web.Page)
	lastId WebpageId = 0
)

type OpenWebpageReturn struct {
	Id WebpageId `json:"id"`
}

type OpenWebpageArgument struct {
	Url string  `json:"url"`
}

func openWebpage(arg OpenWebpageArgument) (OpenWebpageReturn, error) {
	lastId++
	return OpenWebpageReturn{lastId}, nil
}
