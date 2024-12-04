package utils

import (
	"net/http"
	"time"
)

type DelayedTransport struct {
	Transport http.RoundTripper
	Delay     time.Duration
}

func (d *DelayedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	time.Sleep(d.Delay)
	return d.Transport.RoundTrip(req)
}
