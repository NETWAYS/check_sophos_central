package api

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

type LoggingRoundTripper struct {
	Base http.RoundTripper
}

func (r LoggingRoundTripper) RoundTrip(req *http.Request) (res *http.Response, err error) {
	log.WithFields(log.Fields{
		"url":    req.URL,
		"method": req.Method,
	}).Debug("executing http request")

	res, err = r.Base.RoundTrip(req)

	return
}
