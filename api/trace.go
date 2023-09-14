package api

import (
	"net/http"
)

type LoggingRoundTripper struct {
	Base http.RoundTripper
}

func (r LoggingRoundTripper) RoundTrip(req *http.Request) (res *http.Response, err error) {
	res, err = r.Base.RoundTrip(req)

	return
}
