package api

import "git.sr.ht/~michl/quickbeam/internal/webpage"

type WebpageId uint64

var (
	webpages = make(map[WebpageId] webpage.Webpage)
	lastId WebpageId = 0
)

func openWebpage(url string) (WebpageId, error) {
	lastId++
	return lastId, nil
}
