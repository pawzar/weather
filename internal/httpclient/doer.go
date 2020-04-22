package httpclient

import (
	"net/http"
	"time"

	"weather/internal/log"
)

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type DoerFunc func(*http.Request) (*http.Response, error)

func (f DoerFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

func DoerWrappedClient(ttl time.Duration) DoerFunc {
	client := &http.Client{Timeout: ttl}

	return func(req *http.Request) (*http.Response, error) { return client.Do(req) }
}

func LoggingMiddleware(l log.Logger) func(DoerFunc) DoerFunc {
	return func(f DoerFunc) DoerFunc {
		return func(req *http.Request) (*http.Response, error) {
			l.Debugf(">> %s %s %s", req.Method, req.URL, req.Proto)

			res, err := f.Do(req)
			if err != nil {
				return nil, err
			}

			l.Debugf("<< %s %s %d", res.Proto, res.Status, res.ContentLength)

			return res, err
		}
	}
}
