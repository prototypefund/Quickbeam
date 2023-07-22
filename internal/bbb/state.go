package bbb

import "git.sr.ht/~michl/quickbeam/internal/web"

func State(w web.Page) (res []string, err error) {
	r, err := w.Root()
	if err != nil {
		return
	}
	_, ok, err := r.MaybeSubNode("footer p", "Powered by Greenlight")
	if err != nil {
		return nil, err
	}
	if ok {
		res = append(res, "greenlight")
	}
	return
}
